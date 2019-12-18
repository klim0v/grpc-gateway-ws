package service

import (
	"context"
	"encoding/json"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/klim0v/grpc-gateway-ws/pb"
	"strconv"
	"time"
)

func (s *Service) Genesis(context.Context, *empty.Empty) (*pb.GenesisResponse, error) {
	result, err := s.client.Genesis()
	if err != nil {
		return new(pb.GenesisResponse), err //todo
	}

	appState := new(pb.GenesisResponse_Result_Genesis_AppState)
	err = json.Unmarshal(result.Genesis.AppState, appState)
	if err != nil {
		return new(pb.GenesisResponse), err //todo
	}

	return &pb.GenesisResponse{
		Jsonrpc: "2.0",
		Id:      "",
		Result: &pb.GenesisResponse_Result{
			Genesis: &pb.GenesisResponse_Result_Genesis{
				GenesisTime: result.Genesis.GenesisTime.Format(time.RFC3339Nano),
				ChainId:     result.Genesis.ChainID,
				ConsensusParams: &pb.GenesisResponse_Result_Genesis_ConsensusParams{
					Block: &pb.GenesisResponse_Result_Genesis_ConsensusParams_Block{
						MaxBytes:   strconv.Itoa(int(result.Genesis.ConsensusParams.Block.MaxBytes)),
						MaxGas:     strconv.Itoa(int(result.Genesis.ConsensusParams.Block.MaxGas)),
						TimeIotaMs: strconv.Itoa(int(result.Genesis.ConsensusParams.Block.TimeIotaMs)),
					},
					Evidence: &pb.GenesisResponse_Result_Genesis_ConsensusParams_Evidence{
						MaxAge: strconv.Itoa(int(result.Genesis.ConsensusParams.Evidence.MaxAge)),
					},
					Validator: &pb.GenesisResponse_Result_Genesis_ConsensusParams_Validator{
						PublicKeyTypes: result.Genesis.ConsensusParams.Validator.PubKeyTypes,
					},
				},
				AppHash:  result.Genesis.AppHash.String(),
				AppState: appState,
			},
		},
		Error: nil,
	}, nil
}
