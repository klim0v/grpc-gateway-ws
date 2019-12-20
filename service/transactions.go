package service

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/MinterTeam/minter-go-node/core/transaction"
	_struct "github.com/golang/protobuf/ptypes/struct"
	"github.com/klim0v/grpc-gateway-ws/pb"
	"github.com/tendermint/tendermint/libs/common"
)

func (s *Service) Transactions(_ context.Context, req *pb.TransactionsRequest) (*pb.TransactionsResponse, error) {
	page := int(req.Page)
	if page == 0 {
		page = 1
	}
	perPage := int(req.PerPage)
	if perPage == 0 {
		perPage = 100
	}

	rpcResult, err := s.client.TxSearch(req.Query, false, page, perPage)
	if err != nil {
		return &pb.TransactionsResponse{
			Error: &pb.Error{
				Data: err.Error(),
			},
		}, nil
	}

	result := make([]*pb.TransactionResult, len(rpcResult.Txs))
	for i, tx := range rpcResult.Txs {
		decodedTx, _ := transaction.TxDecoder.DecodeFromBytes(tx.Tx)
		sender, _ := decodedTx.Sender()

		tags := make(map[string]string)
		for _, tag := range tx.TxResult.Events[0].Attributes {
			tags[string(tag.Key)] = string(tag.Value)
		}

		data, err := s.encodeTxData(decodedTx)
		if err != nil {
			return &pb.TransactionsResponse{
				Error: &pb.Error{
					Data: err.Error(),
				},
			}, nil
		}

		dataStruct := &_struct.Struct{Fields: make(map[string]*_struct.Value)}
		err = json.Unmarshal(data, dataStruct.Fields)
		if err != nil {
			return &pb.TransactionsResponse{
				Error: &pb.Error{
					Data: err.Error(),
				},
			}, nil
		}

		result[i] = &pb.TransactionResult{
			Hash:     common.HexBytes(tx.Tx.Hash()).String(),
			RawTx:    fmt.Sprintf("%x", []byte(tx.Tx)),
			Height:   fmt.Sprintf("%d", tx.Height),
			Index:    fmt.Sprintf("%d", tx.Index),
			From:     sender.String(),
			Nonce:    fmt.Sprintf("%d", decodedTx.Nonce),
			GasPrice: fmt.Sprintf("%d", decodedTx.GasPrice),
			GasCoin:  decodedTx.GasCoin.String(),
			Gas:      fmt.Sprintf("%d", decodedTx.Gas()),
			Type:     fmt.Sprintf("%d", uint8(decodedTx.Type)),
			Data:     dataStruct,
			Payload:  decodedTx.Payload,
			Tags:     tags,
			Code:     fmt.Sprintf("%d", tx.TxResult.Code),
			Log:      tx.TxResult.Log,
		}
	}

	return &pb.TransactionsResponse{
		Result: result,
	}, nil
}
