package main

import (
	"context" // Use "golang.org/x/net/context" for Golang version <= 1.6
	"fmt"
	"github.com/klim0v/grpc-gateway-ws/service"
	"google.golang.org/grpc"
	"log"
	"net"
	"net/http"

	"github.com/golang/glog"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	gw "github.com/klim0v/grpc-gateway-ws/pb" // Update
	"github.com/tmc/grpc-websocket-proxy/wsproxy"
)

func startGRPC() error {
	lis, err := net.Listen("tcp", ":8081")
	if err != nil {
		return err
	}
	grpcServer := grpc.NewServer()
	gw.RegisterWebsocketServiceServer(grpcServer, &service.WebsocketService{})
	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			log.Println("serveGRPC err:", err)
		}
	}()
	return nil
}

func run() error {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithInsecure()}
	err := gw.RegisterWebsocketServiceHandlerFromEndpoint(ctx, mux, ":8081", opts)
	if err != nil {
		return err
	}

	fmt.Println("listening")

	return http.ListenAndServe(":8000", wsproxy.WebsocketProxy(mux))
}

func main() {

	if err := startGRPC(); err != nil {
		glog.Fatal(err)
	}

	if err := run(); err != nil {
		glog.Fatal(err)
	}
}
