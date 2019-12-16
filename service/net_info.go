package service

import (
	"context"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/klim0v/grpc-gateway-ws/pb"
	"strconv"
	"time"
)

func (s *Service) NetInfo(context.Context, *empty.Empty) (*pb.NetInfoResponse, error) {
	result, err := s.client.NetInfo()
	if err != nil {
		return new(pb.NetInfoResponse), err //todo
	}

	var peers []*pb.NetInfoResponse_Result_Peer
	for _, peer := range result.Peers {
		var channels []*pb.NetInfoResponse_Result_Peer_ConnectionStatus_Channel
		for _, channel := range peer.ConnectionStatus.Channels {
			channels = append(channels, &pb.NetInfoResponse_Result_Peer_ConnectionStatus_Channel{
				ID:                strconv.Itoa(int(channel.ID)),
				SendQueueCapacity: strconv.Itoa(channel.SendQueueCapacity),
				SendQueueSize:     strconv.Itoa(channel.SendQueueSize),
				Priority:          strconv.Itoa(channel.Priority),
				RecentlySent:      strconv.Itoa(int(channel.RecentlySent)),
			})
		}

		peers = append(peers, &pb.NetInfoResponse_Result_Peer{
			NodeInfo: &pb.NodeInfo{
				ProtocolVersion: &pb.NodeInfo_ProtocolVersion{
					P2P:   strconv.Itoa(int(peer.NodeInfo.ProtocolVersion.P2P)),
					Block: strconv.Itoa(int(peer.NodeInfo.ProtocolVersion.Block)),
					App:   strconv.Itoa(int(peer.NodeInfo.ProtocolVersion.App)),
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
				Duration: strconv.Itoa(int(peer.ConnectionStatus.Duration)),
				SendMonitor: &pb.NetInfoResponse_Result_Peer_ConnectionStatus_Monitor{
					Active:   false,
					Start:    peer.ConnectionStatus.SendMonitor.Start.Format(time.RFC3339Nano),
					Duration: strconv.Itoa(int(peer.ConnectionStatus.SendMonitor.Duration.Nanoseconds())),
					Idle:     strconv.Itoa(int(peer.ConnectionStatus.SendMonitor.Idle.Nanoseconds())),
					Bytes:    strconv.Itoa(int(peer.ConnectionStatus.SendMonitor.Bytes)),
					Samples:  strconv.Itoa(int(peer.ConnectionStatus.SendMonitor.Samples)),
					InstRate: strconv.Itoa(int(peer.ConnectionStatus.SendMonitor.InstRate)),
					CurRate:  strconv.Itoa(int(peer.ConnectionStatus.SendMonitor.CurRate)),
					AvgRate:  strconv.Itoa(int(peer.ConnectionStatus.SendMonitor.AvgRate)),
					PeakRate: strconv.Itoa(int(peer.ConnectionStatus.SendMonitor.PeakRate)),
					BytesRem: strconv.Itoa(int(peer.ConnectionStatus.SendMonitor.BytesRem)),
					TimeRem:  strconv.Itoa(int(peer.ConnectionStatus.SendMonitor.TimeRem.Nanoseconds())),
					Progress: strconv.Itoa(int(peer.ConnectionStatus.SendMonitor.Progress)),
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
			NPeers:    strconv.Itoa(result.NPeers),
			Peers:     peers,
		},
	}, nil
}
