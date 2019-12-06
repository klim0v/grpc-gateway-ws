package service

import (
	"context"
	"fmt"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/klim0v/grpc-gateway-ws/pb"
	"time"
)

func (s *Service) Subscribe(request *pb.SubscribeRequest, stream pb.WebsocketService_SubscribeServer) error {
	start := time.Now()
	for i := 0; i < 5; i++ {
		time.Sleep(time.Second)
		if err := stream.Send(&pb.SubscribeResponse{
			Query: "hello there!" + fmt.Sprint(time.Now().Sub(start)),
		}); err != nil {
			return err
		}
	}
	return nil
}

func (s *Service) Unsubscribe(context.Context, *pb.UnsubscribeRequest) (*empty.Empty, error) {
	panic("implement me")
}

func (s *Service) UnsubscribeAll(context.Context, *empty.Empty) (*empty.Empty, error) {
	panic("implement me")
}
