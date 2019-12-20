package service

import (
	"context"
	"fmt"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/klim0v/grpc-gateway-ws/pb"
)

func (s *Service) MinGasPrice(context.Context, *empty.Empty) (*pb.MinGasPriceResponse, error) {
	return &pb.MinGasPriceResponse{
		Jsonrpc: "2.0",
		Id:      "",
		Result:  fmt.Sprintf("%d", s.blockchain.MinGasPrice()),
	}, nil
}

func (s *Service) MaxGas(_ context.Context, req *pb.MaxGasRequest) (*pb.MaxGasResponse, error) {
	cState, err := s.getStateForHeight(req.Height)
	if err != nil {
		return &pb.MaxGasResponse{
			Error: &pb.Error{
				Data: err.Error(),
			},
		}, nil
	}

	return &pb.MaxGasResponse{
		Jsonrpc: "2.0",
		Id:      "",
		Result:  fmt.Sprintf("%d", cState.App.GetMaxGas()),
	}, nil
}
