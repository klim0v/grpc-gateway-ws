package service

import (
	"context"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/klim0v/grpc-gateway-ws/pb"
	"time"
)

func (s *Service) NetInfo(context.Context, *empty.Empty) (*pb.NetInfoResponse, error) {
	result, err := s.client.NetInfo()
	if err != nil {
		return new(pb.NetInfoResponse), err
	}

	var peers []*pb.NetInfoResponse_Result_Peer
	for _, peer := range result.Peers {
		var channels []*pb.NetInfoResponse_Result_Peer_ConnectionStatus_Channel
		for _, channel := range peer.ConnectionStatus.Channels {
			channels = append(channels, &pb.NetInfoResponse_Result_Peer_ConnectionStatus_Channel{
				ID:                int32(channel.ID),
				SendQueueCapacity: int32(channel.SendQueueCapacity),
				SendQueueSize:     int32(channel.SendQueueSize),
				Priority:          int32(channel.Priority),
				RecentlySent:      channel.RecentlySent,
			})
		}

		peers = append(peers, &pb.NetInfoResponse_Result_Peer{
			NodeInfo: &pb.NodeInfo{
				ProtocolVersion: &pb.NodeInfo_ProtocolVersion{
					P2P:   uint64(peer.NodeInfo.ProtocolVersion.P2P),
					Block: uint64(peer.NodeInfo.ProtocolVersion.Block),
					App:   uint64(peer.NodeInfo.ProtocolVersion.App),
				},
				Id:         string(peer.NodeInfo.ID_),
				ListenAddr: peer.NodeInfo.ListenAddr,
				Network:    peer.NodeInfo.Network,
				Version:    peer.NodeInfo.Version,
				Channels:   peer.NodeInfo.Channels.String(),
				Moniker:    peer.NodeInfo.Moniker,
				Other: &pb.NodeInfo_Other{
					TxIndex:    peer.NodeInfo.Other.TxIndex,
					RpcAddress: peer.NodeInfo.Other.RPCAddress,
				},
			},
			IsOutbound: peer.IsOutbound,
			ConnectionStatus: &pb.NetInfoResponse_Result_Peer_ConnectionStatus{
				Duration: int64(peer.ConnectionStatus.Duration),
				SendMonitor: &pb.NetInfoResponse_Result_Peer_ConnectionStatus_Monitor{
					Active:   false,
					Start:    peer.ConnectionStatus.SendMonitor.Start.Format(time.RFC3339Nano),
					Duration: peer.ConnectionStatus.SendMonitor.Duration.Nanoseconds(),
					Idle:     peer.ConnectionStatus.SendMonitor.Idle.Nanoseconds(),
					Bytes:    peer.ConnectionStatus.SendMonitor.Bytes,
					Samples:  peer.ConnectionStatus.SendMonitor.Samples,
					InstRate: peer.ConnectionStatus.SendMonitor.InstRate,
					CurRate:  peer.ConnectionStatus.SendMonitor.CurRate,
					AvgRate:  peer.ConnectionStatus.SendMonitor.AvgRate,
					PeakRate: peer.ConnectionStatus.SendMonitor.PeakRate,
					BytesRem: peer.ConnectionStatus.SendMonitor.BytesRem,
					TimeRem:  peer.ConnectionStatus.SendMonitor.TimeRem.Nanoseconds(),
					Progress: uint32(peer.ConnectionStatus.SendMonitor.Progress),
				},
				RecvMonitor: nil,
				Channels:    channels,
			},
			RemoteIp: peer.RemoteIP,
		})
	}

	return &pb.NetInfoResponse{
		Jsonrpc: "2.0",
		Id:      "",
		Result: &pb.NetInfoResponse_Result{
			Listening: result.Listening,
			Listeners: result.Listeners,
			NPeers:    int32(result.NPeers),
			Peers:     peers,
		},
	}, nil
}
