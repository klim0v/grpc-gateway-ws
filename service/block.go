package service

import (
	"bytes"
	"context"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/MinterTeam/minter-go-node/core/rewards"
	"github.com/MinterTeam/minter-go-node/core/transaction"
	"github.com/MinterTeam/minter-go-node/core/types"
	_struct "github.com/golang/protobuf/ptypes/struct"
	"github.com/klim0v/grpc-gateway-ws/pb"
	core_types "github.com/tendermint/tendermint/rpc/core/types"
	"strconv"
	"time"
)

func (s *Service) Block(_ context.Context, req *pb.BlockRequest) (*pb.BlockResponse, error) {
	block, err := s.client.Block(&req.Height)
	if err != nil {
		return &pb.BlockResponse{
			Error: &pb.Error{
				Code:    "404",
				Message: "Block not found",
				Data:    err.Error(),
			},
		}, nil
	}

	blockResults, err := s.client.BlockResults(&req.Height)
	if err != nil {
		return &pb.BlockResponse{
			Error: &pb.Error{
				Code:    "404",
				Message: "Block results not found",
				Data:    err.Error(),
			},
		}, nil
	}

	valHeight := req.Height - 1
	if valHeight < 1 {
		valHeight = 1
	}
	tmValidators, err := s.client.Validators(&valHeight)
	if err != nil {
		return &pb.BlockResponse{
			Error: &pb.Error{
				Code:    "404",
				Message: "Validators for block not found",
				Data:    err.Error(),
			},
		}, nil
	}

	txs := make([]*pb.BlockResponse_Result_Transaction, len(block.Block.Data.Txs))
	for i, rawTx := range block.Block.Data.Txs {
		tx, _ := transaction.TxDecoder.DecodeFromBytes(rawTx)
		sender, _ := tx.Sender()

		tags := make(map[string]string)
		for _, tag := range blockResults.Results.DeliverTx[i].Events[0].Attributes {
			tags[string(tag.Key)] = string(tag.Value)
		}

		data, err := s.encodeTxData(tx)
		if err != nil {
			return new(pb.BlockResponse), err
		}

		dataStruct := new(_struct.Struct)
		err = json.Unmarshal(data, dataStruct)
		if err != nil {
			return new(pb.BlockResponse), err
		}

		txs[i] = &pb.BlockResponse_Result_Transaction{
			Hash:        fmt.Sprintf("Mt%x", rawTx.Hash()),
			RawTx:       fmt.Sprintf("%x", []byte(rawTx)),
			From:        sender.String(),
			Nonce:       fmt.Sprintf("%d", tx.Nonce),
			GasPrice:    fmt.Sprintf("%d", tx.GasPrice),
			Type:        fmt.Sprintf("%d", tx.Type),
			Data:        dataStruct, //todo
			Payload:     string(tx.Payload),
			ServiceData: string(tx.ServiceData),
			Gas:         fmt.Sprintf("%d", tx.Gas()),
			GasCoin:     tx.GasCoin.String(),
			Tags:        tags,
			Code:        fmt.Sprintf("%d", blockResults.Results.DeliverTx[i].Code),
			Log:         blockResults.Results.DeliverTx[i].Log,
		}
	}

	var validators []*pb.BlockResponse_Result_Validator
	var proposer string
	if req.Height > 1 {
		p, err := s.getBlockProposer(block)
		if err != nil {
			return &pb.BlockResponse{
				Error: err,
			}, nil
		}

		if p != nil {
			str := p.String()
			proposer = str
		}

		validators = make([]*pb.BlockResponse_Result_Validator, len(tmValidators.Validators))
		for i, tmval := range tmValidators.Validators {
			signed := false
			for _, vote := range block.Block.LastCommit.Precommits {
				if vote == nil {
					continue
				}

				if bytes.Equal(vote.ValidatorAddress.Bytes(), tmval.Address.Bytes()) {
					signed = true
					break
				}
			}

			validators[i] = &pb.BlockResponse_Result_Validator{
				PublicKey: fmt.Sprintf("Mp%x", tmval.PubKey.Bytes()[5:]),
				Signed:    signed,
			}
		}
	}

	return &pb.BlockResponse{
		Result: &pb.BlockResponse_Result{
			Hash:         hex.EncodeToString(block.Block.Hash()),
			Height:       strconv.Itoa(int(block.Block.Height)),
			Time:         block.Block.Time.Format(time.RFC3339Nano),
			NumTxs:       strconv.Itoa(int(block.Block.NumTxs)),
			TotalTxs:     strconv.Itoa(int(block.Block.TotalTxs)),
			Transactions: txs,
			BlockReward:  rewards.GetRewardForBlock(uint64(req.Height)).String(),
			Size:         strconv.Itoa(len(s.cdc.MustMarshalBinaryLengthPrefixed(block))),
			Proposer:     proposer,
			Validators:   validators,
			Evidence: &pb.BlockResponse_Result_Evidence{
				Evidence: make([]*pb.BlockResponse_Result_Evidence_Evidence, len(block.Block.Evidence.Evidence)), // todo
			},
		},
	}, nil
}

func (s *Service) getBlockProposer(block *core_types.ResultBlock) (*types.Pubkey, *pb.Error) {
	vals, err := s.client.Validators(&block.Block.Height)
	if err != nil {
		return nil, &pb.Error{Code: "404", Message: "Validators for block not found", Data: err.Error()}
	}

	for _, tmval := range vals.Validators {
		if bytes.Equal(tmval.Address.Bytes(), block.Block.ProposerAddress.Bytes()) {
			var result types.Pubkey
			copy(result[:], tmval.PubKey.Bytes()[5:])
			return &result, nil
		}
	}

	return nil, &pb.Error{Code: "404", Message: "Block proposer not found"}
}

func (s *Service) encodeTxData(decodedTx *transaction.Transaction) ([]byte, error) {
	switch decodedTx.Type {
	case transaction.TypeSend:
		return s.cdc.MarshalJSON(decodedTx.GetDecodedData().(*transaction.SendData))
	case transaction.TypeRedeemCheck:
		return s.cdc.MarshalJSON(decodedTx.GetDecodedData().(*transaction.RedeemCheckData))
	case transaction.TypeSellCoin:
		return s.cdc.MarshalJSON(decodedTx.GetDecodedData().(*transaction.SellCoinData))
	case transaction.TypeSellAllCoin:
		return s.cdc.MarshalJSON(decodedTx.GetDecodedData().(*transaction.SellAllCoinData))
	case transaction.TypeBuyCoin:
		return s.cdc.MarshalJSON(decodedTx.GetDecodedData().(*transaction.BuyCoinData))
	case transaction.TypeCreateCoin:
		return s.cdc.MarshalJSON(decodedTx.GetDecodedData().(*transaction.CreateCoinData))
	case transaction.TypeDeclareCandidacy:
		return s.cdc.MarshalJSON(decodedTx.GetDecodedData().(*transaction.DeclareCandidacyData))
	case transaction.TypeDelegate:
		return s.cdc.MarshalJSON(decodedTx.GetDecodedData().(*transaction.DelegateData))
	case transaction.TypeSetCandidateOnline:
		return s.cdc.MarshalJSON(decodedTx.GetDecodedData().(*transaction.SetCandidateOnData))
	case transaction.TypeSetCandidateOffline:
		return s.cdc.MarshalJSON(decodedTx.GetDecodedData().(*transaction.SetCandidateOffData))
	case transaction.TypeUnbond:
		return s.cdc.MarshalJSON(decodedTx.GetDecodedData().(*transaction.UnbondData))
	case transaction.TypeMultisend:
		return s.cdc.MarshalJSON(decodedTx.GetDecodedData().(*transaction.MultisendData))
	case transaction.TypeCreateMultisig:
		return s.cdc.MarshalJSON(decodedTx.GetDecodedData().(*transaction.CreateMultisigData))
	case transaction.TypeEditCandidate:
		return s.cdc.MarshalJSON(decodedTx.GetDecodedData().(*transaction.EditCandidateData))
	}

	return nil, errors.New("unknown tx type")
}
