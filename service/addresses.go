package service

import (
	"context"
	"fmt"
	"github.com/MinterTeam/minter-go-node/core/types"
	"github.com/klim0v/grpc-gateway-ws/pb"
)

func (s *Service) Addresses(_ context.Context, req *pb.AddressesRequest) (*pb.AddressesResponse, error) {
	cState, err := s.getStateForHeight(req.Height)
	if err != nil {
		return &pb.AddressesResponse{
			Error: &pb.Error{
				Data: err.Error(),
			},
		}, nil
	}

	response := &pb.AddressesResponse{
		Result: make([]*pb.AddressesResponse_Result, len(req.Addresses)),
	}

	for i, address := range req.Addresses {
		addr := types.StringToAddress(address)
		data := &pb.AddressesResponse_Result{
			Address:          address,
			Balance:          make(map[string]string),
			TransactionCount: fmt.Sprintf("%d", cState.Accounts.GetNonce(addr)),
		}

		balances := cState.Accounts.GetBalances(addr)
		for k, v := range balances {
			data.Balance[k.String()] = v.String()
		}

		if _, exists := data.Balance[types.GetBaseCoin().String()]; !exists {
			data.Balance[types.GetBaseCoin().String()] = "0"
		}

		response.Result[i] = data
	}

	return response, nil
}
