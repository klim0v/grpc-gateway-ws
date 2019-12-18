package service

import (
	"context"
	"github.com/MinterTeam/minter-go-node/core/types"
	"github.com/klim0v/grpc-gateway-ws/pb"
	"strconv"
)

func (s *Service) Address(_ context.Context, req *pb.AddressRequest) (*pb.AddressResponse, error) {
	cState, err := s.getStateForHeight(req.Height)
	if err != nil {
		return new(pb.AddressResponse), err
	}

	address := types.StringToAddress(req.Address)
	response := &pb.AddressResponse{
		Result: &pb.AddressResponse_Result{
			Balance:          make(map[string]string),
			TransactionCount: strconv.Itoa(int(cState.Accounts.GetNonce(address))),
		},
	}

	balances := cState.Accounts.GetBalances(address)

	for k, v := range balances {
		response.Result.Balance[k.String()] = v.String()
	}

	if _, exists := response.Result.Balance[types.GetBaseCoin().String()]; !exists {
		response.Result.Balance[types.GetBaseCoin().String()] = "0"
	}

	return response, nil
}
