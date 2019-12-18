package service

import (
	"context"
	"fmt"
	"github.com/MinterTeam/minter-go-node/core/types"
	"github.com/klim0v/grpc-gateway-ws/pb"
)

func (s *Service) Validators(_ context.Context, req *pb.ValidatorsRequest) (*pb.ValidatorsResponse, error) {
	height := uint64(req.Height)
	if height == 0 {
		height = s.blockchain.Height()
	}

	h := int64(height)
	tmVals, err := s.client.Validators(&h)
	if err != nil {
		return &pb.ValidatorsResponse{Error: &pb.Error{
			Message: err.Error(),
		}}, nil
	}

	responseValidators := make([]*pb.ValidatorsResponse_Result, len(tmVals.Validators))
	for i, val := range tmVals.Validators {
		var pk types.Pubkey
		copy(pk[:], val.PubKey.Bytes()[5:])
		responseValidators[i] = &pb.ValidatorsResponse_Result{
			PublicKey:   pk.String(),
			VotingPower: fmt.Sprintf("%d", val.VotingPower),
		}
	}
	return &pb.ValidatorsResponse{Result: responseValidators}, nil
}
