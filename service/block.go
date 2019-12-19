package service

import (
	"bytes"
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/MinterTeam/minter-go-node/core/rewards"
	"github.com/MinterTeam/minter-go-node/core/transaction"
	"github.com/MinterTeam/minter-go-node/core/types"
	_struct "github.com/golang/protobuf/ptypes/struct"
	"github.com/klim0v/grpc-gateway-ws/pb"
	core_types "github.com/tendermint/tendermint/rpc/core/types"
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
			return &pb.BlockResponse{
				Error: &pb.Error{
					Data: err.Error(),
				},
			}, nil
		}

		dataStruct := &_struct.Struct{Fields: make(map[string]*_struct.Value)}
		err = json.Unmarshal(data, dataStruct.Fields)
		if err != nil {
			return &pb.BlockResponse{
				Error: &pb.Error{
					Data: err.Error(),
				},
			}, nil
		}

		txs[i] = &pb.BlockResponse_Result_Transaction{
			Hash:        fmt.Sprintf("Mt%x", rawTx.Hash()),
			RawTx:       fmt.Sprintf("%x", []byte(rawTx)),
			From:        sender.String(),
			Nonce:       fmt.Sprintf("%d", tx.Nonce),
			GasPrice:    fmt.Sprintf("%d", tx.GasPrice),
			Type:        fmt.Sprintf("%d", tx.Type),
			Data:        dataStruct, //todo
			Payload:     tx.Payload,
			ServiceData: tx.ServiceData,
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
			Height:       fmt.Sprintf("%d", block.Block.Height),
			Time:         block.Block.Time.Format(time.RFC3339Nano),
			NumTxs:       fmt.Sprintf("%d", block.Block.NumTxs),
			TotalTxs:     fmt.Sprintf("%d", block.Block.TotalTxs),
			Transactions: txs,
			BlockReward:  rewards.GetRewardForBlock(uint64(req.Height)).String(),
			Size:         fmt.Sprintf("%d", s.cdc.MustMarshalBinaryLengthPrefixed(block)),
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
