package service

import (
	"context"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/klim0v/grpc-gateway-ws/pb"
)

func (s *Service) Status(context.Context, *empty.Empty) (*pb.StatusResponse, error) {
	return &pb.StatusResponse{
		Jsonrpc: "2.0",
		Id:      "test",
		Result:  nil,
	}, nil
}

func (s *Service) NetInfo(context.Context, *empty.Empty) (*pb.NetInfoResponse, error) {
	panic("implement me")
}

func (s *Service) MinGasPrice(context.Context, *empty.Empty) (*pb.MinGasPriceResponse, error) {
	panic("implement me")
}

func (s *Service) Genesis(context.Context, *empty.Empty) (*pb.GenesisResponse, error) {
	panic("implement me")
}
