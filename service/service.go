package service

import (
	"fmt"
	"github.com/klim0v/grpc-gateway-ws/pb"
	"time"
)

type WebsocketService struct {
}

func (w *WebsocketService) Echo(_ *pb.Empty, stream pb.WebsocketService_EchoServer) error {
	start := time.Now()
	for i := 0; i < 5; i++ {
		time.Sleep(time.Second)
		if err := stream.Send(&pb.Response{
			Value: "hello there!" + fmt.Sprint(time.Now().Sub(start)),
		}); err != nil {
			return err
		}
	}
	return nil
}
