package main

import (
	"context"
	"fmt"
	"github.com/klim0v/grpc-gateway-ws/service"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"log"
	"net"
	"net/http"

	"github.com/golang/glog"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	gw "github.com/klim0v/grpc-gateway-ws/pb"
	"github.com/tmc/grpc-websocket-proxy/wsproxy"
)

func startGRPC() error {
	lis, err := net.Listen("tcp", ":8081")
	if err != nil {
		return err
	}
	grpcServer := grpc.NewServer()
	gw.RegisterWebsocketServiceServer(grpcServer, &service.Service{})
	gw.RegisterHttpServiceServer(grpcServer, &service.Service{})
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

	if err := gw.RegisterWebsocketServiceHandlerFromEndpoint(ctx, mux, ":8081", opts); err != nil {
		return err
	}

	if err := gw.RegisterHttpServiceHandlerFromEndpoint(ctx, mux, ":8081", opts); err != nil {
		return err
	}

	fmt.Println("listening")

	var group errgroup.Group
	group.Go(func() error {
		return http.ListenAndServe(":8000", mux)
	})
	group.Go(func() error {
		return http.ListenAndServe(":8000", wsproxy.WebsocketProxy(mux))
	})

	if err := group.Wait(); err != nil {
		return err
	}

	return nil
}

func main() {

	if err := startGRPC(); err != nil {
		glog.Fatal(err)
	}

	if err := run(); err != nil {
		glog.Fatal(err)
	}
}
