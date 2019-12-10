package main

import (
	"context"
	"fmt"
	"github.com/klim0v/grpc-gateway-ws/service"
	"golang.org/x/sync/errgroup"
	"net/http"

	"github.com/golang/glog"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	gw "github.com/klim0v/grpc-gateway-ws/pb"
)

func run() error {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	mux := runtime.NewServeMux()

	if err := gw.RegisterHttpServiceHandlerServer(ctx, mux, &service.Service{}); err != nil {
		return err
	}

	mux.Handle("GET", runtime.MustPattern(runtime.NewPattern(1, []int{2, 0}, []string{"subscribe"}, "", runtime.AssumeColonVerbOpt(true))), func(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {

	})

	fmt.Println("listening")

	var group errgroup.Group
	group.Go(func() error {
		return http.ListenAndServe(":8000", mux)
	})

	if err := group.Wait(); err != nil {
		return err
	}

	return nil
}

func main() {
	if err := run(); err != nil {
		glog.Fatal(err)
	}
}
