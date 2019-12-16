package service

import (
	"context"
	"fmt"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/klim0v/grpc-gateway-ws/pb"
	"time"
)

func (s *Service) Status(context.Context, *empty.Empty) (*pb.StatusResponse, error) {
	result, err := s.client.Status()
	if err != nil {
		return new(pb.StatusResponse), err
	}
	return &pb.StatusResponse{
		Jsonrpc: "2.0",
		Id:      "",
		Result: &pb.StatusResponse_Result{
			Version:           s.version,
			LatestBlockHash:   fmt.Sprintf("%X", result.SyncInfo.LatestBlockHash),
			LatestAppHash:     fmt.Sprintf("%X", result.SyncInfo.LatestAppHash),
			LatestBlockHeight: result.SyncInfo.LatestBlockHeight,
			LatestBlockTime:   result.SyncInfo.LatestBlockTime.Format(time.RFC3339Nano),
			KeepLastStates:    s.minterCfg.BaseConfig.KeepLastStates,
			TmStatus: &pb.StatusResponse_Result_TmStatus{
				NodeInfo: &pb.NodeInfo{
					ProtocolVersion: &pb.NodeInfo_ProtocolVersion{
						P2P:   uint64(result.NodeInfo.ProtocolVersion.P2P),
						Block: uint64(result.NodeInfo.ProtocolVersion.Block),
						App:   uint64(result.NodeInfo.ProtocolVersion.App),
					},
					Id:         string(result.NodeInfo.ID_),
					ListenAddr: result.NodeInfo.ListenAddr,
					Network:    result.NodeInfo.Network,
					Version:    result.NodeInfo.Version,
					Channels:   result.NodeInfo.Channels.String(),
					Moniker:    result.NodeInfo.Moniker,
					Other: &pb.NodeInfo_Other{
						TxIndex:    result.NodeInfo.Other.TxIndex,
						RpcAddress: result.NodeInfo.Other.RPCAddress,
					},
				},
				SyncInfo: &pb.StatusResponse_Result_TmStatus_SyncInfo{
					LatestBlockHash:   result.SyncInfo.LatestBlockHash.String(),
					LatestAppHash:     result.SyncInfo.LatestAppHash.String(),
					LatestBlockHeight: result.SyncInfo.LatestBlockHeight,
					LatestBlockTime:   result.SyncInfo.LatestBlockTime.Format(time.RFC3339Nano),
					CatchingUp:        result.SyncInfo.CatchingUp,
				},
				ValidatorInfo: &pb.StatusResponse_Result_TmStatus_ValidatorInfo{
					Address: result.ValidatorInfo.Address.String(),
					PubKey:  &pb.StatusResponse_Result_TmStatus_ValidatorInfo_PubKey{
						//Type:  result.ValidatorInfo.PubKey., todo
						//Value: result.ValidatorInfo.PubKey., todo
					},
					VotingPower: result.ValidatorInfo.VotingPower,
				},
			},
		},
	}, nil
}
