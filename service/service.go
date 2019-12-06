package service

import (
	"github.com/MinterTeam/minter-go-node/config"
	"github.com/MinterTeam/minter-go-node/core/minter"
	"github.com/tendermint/go-amino"
	rpc "github.com/tendermint/tendermint/rpc/client"
)

type Service struct {
	cdc        *amino.Codec
	blockchain *minter.Blockchain
	client     *rpc.Local
	minterCfg  *config.Config
}

func NewService(blockchain *minter.Blockchain, client *rpc.Local, minterCfg *config.Config) *Service {
	return &Service{blockchain: blockchain, client: client, minterCfg: minterCfg, cdc: amino.NewCodec()}
}
