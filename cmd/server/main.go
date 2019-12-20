package main

import (
	"context"
	"fmt"
	"github.com/MinterTeam/minter-go-node/core/minter"
	"github.com/klim0v/grpc-gateway-ws/service"
	rpc "github.com/tendermint/tendermint/rpc/client"
	"net/http"

	"github.com/golang/glog"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	gw "github.com/klim0v/grpc-gateway-ws/pb"

	"github.com/MinterTeam/minter-go-node/config"
	"github.com/tendermint/go-amino"
)

var (
	cdc        = amino.NewCodec()
	blockchain *minter.Blockchain
	client     *rpc.Local
	minterCfg  *config.Config
	version    string
)

func run() error {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	mux := runtime.NewServeMux()

	srvc := service.NewService(cdc, blockchain, client, minterCfg, version)
	if err := gw.RegisterHttpServiceHandlerServer(ctx, mux, srvc); err != nil {
		return err
	}

	mux.Handle("GET", runtime.MustPattern(runtime.NewPattern(1, []int{2, 0}, []string{"subscribe"}, "", runtime.AssumeColonVerbOpt(true))), srvc.Subscribe)

	fmt.Println("listening")

	if err := http.ListenAndServe(":8841", mux); err != nil {
		return err
	}

	return nil
}

func main() {
	if err := run(); err != nil {
		glog.Fatal(err)
	}
}
