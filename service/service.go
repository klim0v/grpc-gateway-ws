package service

import (
	"context"
	"github.com/MinterTeam/minter-go-node/config"
	"github.com/MinterTeam/minter-go-node/core/minter"
	"github.com/MinterTeam/minter-go-node/core/state"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/klim0v/grpc-gateway-ws/pb"
	"github.com/tendermint/go-amino"
	rpc "github.com/tendermint/tendermint/rpc/client"
)

type Service struct {
	cdc        *amino.Codec
	blockchain *minter.Blockchain
	client     *rpc.Local
	minterCfg  *config.Config
	version    string
}

func (s *Service) Address(context.Context, *pb.AddressRequest) (*pb.AddressResponse, error) {
	panic("implement me")
}

func (s *Service) Addresses(context.Context, *pb.AddressesRequest) (*pb.AddressesResponse, error) {
	panic("implement me")
}

func (s *Service) Block(context.Context, *pb.BlockRequest) (*pb.BlockResponse, error) {
	panic("implement me")
}

func (s *Service) Candidate(context.Context, *pb.CandidateRequest) (*pb.CandidateResponse, error) {
	panic("implement me")
}

func (s *Service) Candidates(context.Context, *pb.CandidatesRequest) (*pb.CandidatesResponse, error) {
	panic("implement me")
}

func (s *Service) CoinInfo(context.Context, *pb.CoinInfoRequest) (*pb.CoinInfoResponse, error) {
	panic("implement me")
}

func (s *Service) EstimateCoinBuy(context.Context, *pb.EstimateCoinBuyRequest) (*pb.EstimateCoinBuyResponse, error) {
	panic("implement me")
}

func (s *Service) EstimateCoinSell(context.Context, *pb.EstimateCoinSellRequest) (*pb.EstimateCoinSellResponse, error) {
	panic("implement me")
}

func (s *Service) EstimateCoinSellAll(context.Context, *pb.EstimateCoinSellAllRequest) (*pb.EstimateCoinSellAllResponse, error) {
	panic("implement me")
}

func (s *Service) EstimateTxCommission(context.Context, *pb.EstimateTxCommissionRequest) (*pb.EstimateTxCommissionResponse, error) {
	panic("implement me")
}

func (s *Service) Events(context.Context, *pb.EventsRequest) (*pb.EventsResponse, error) {
	panic("implement me")
}

func (s *Service) MissedBlocks(context.Context, *pb.MissedBlocksRequest) (*pb.MissedBlocksResponse, error) {
	panic("implement me")
}

func (s *Service) SendTransaction(context.Context, *pb.SendTransactionRequest) (*pb.SendTransactionResponse, error) {
	panic("implement me")
}

func (s *Service) Transaction(context.Context, *pb.TransactionRequest) (*pb.TransactionResponse, error) {
	panic("implement me")
}

func (s *Service) Transactions(context.Context, *pb.TransactionsRequest) (*pb.TransactionsResponse, error) {
	panic("implement me")
}

func (s *Service) UnconfirmedTxs(context.Context, *pb.UnconfirmedTxsRequest) (*pb.UnconfirmedTxsResponse, error) {
	panic("implement me")
}

func (s *Service) Validators(context.Context, *pb.ValidatorsRequest) (*pb.ValidatorsResponse, error) {
	panic("implement me")
}

func (s *Service) Genesis(context.Context, *empty.Empty) (*pb.GenesisResponse, error) {
	panic("implement me")
}

func NewService(blockchain *minter.Blockchain, client *rpc.Local, minterCfg *config.Config, version string) *Service {
	return &Service{blockchain: blockchain, client: client, minterCfg: minterCfg, cdc: amino.NewCodec(), version: version}
}

func (s *Service) getStateForHeight(height int) (*state.State, error) {
	if height > 0 {
		cState, err := s.blockchain.GetStateForHeight(uint64(height))

		return cState, err
	}

	return s.blockchain.CurrentState(), nil
}
