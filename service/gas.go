package service

import (
	"context"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/klim0v/grpc-gateway-ws/pb"
)

func (s *Service) MinGasPrice(context.Context, *empty.Empty) (*pb.MinGasPriceResponse, error) {
	return &pb.MinGasPriceResponse{Jsonrpc: "2.0", Result: string(s.blockchain.MinGasPrice())}, nil
}

func (s *Service) MaxGas(_ context.Context, req *pb.MaxGasRequest) (*pb.MaxGasResponse, error) {
	cState, err := s.getStateForHeight(int(req.Height))
	if err != nil {
		return new(pb.MaxGasResponse), err //todo
	}

	return &pb.MaxGasResponse{
		Jsonrpc: "2.0",
		Id:      "",
		Result:  string(cState.App.GetMaxGas()),
	}, nil
}
