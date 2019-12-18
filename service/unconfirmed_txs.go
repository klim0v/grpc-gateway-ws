package service

import (
	"context"
	"fmt"
	"github.com/klim0v/grpc-gateway-ws/pb"
)

func (s *Service) UnconfirmedTxs(_ context.Context, req *pb.UnconfirmedTxsRequest) (*pb.UnconfirmedTxsResponse, error) {
	txs, err := s.client.UnconfirmedTxs(int(req.Limit))
	if err != nil {
		return &pb.UnconfirmedTxsResponse{Error: &pb.Error{Message: err.Error()}}, nil
	}
	return &pb.UnconfirmedTxsResponse{Result: &pb.UnconfirmedTxsResponse_Result{
		NTxs:       fmt.Sprintf("%d", txs.Count),
		Total:      fmt.Sprintf("%d", txs.Total),
		TotalBytes: string(txs.TotalBytes),
	}}, nil
}
