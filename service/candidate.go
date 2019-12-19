package service

import (
	"context"
	"fmt"
	"github.com/MinterTeam/minter-go-node/core/state"
	"github.com/MinterTeam/minter-go-node/core/state/candidates"
	"github.com/MinterTeam/minter-go-node/core/types"
	"github.com/klim0v/grpc-gateway-ws/pb"
)

func (s *Service) Candidate(_ context.Context, req *pb.CandidateRequest) (*pb.CandidateResponse, error) {
	cState, err := s.getStateForHeight(req.Height)
	if err != nil {
		return &pb.CandidateResponse{
			Error: &pb.Error{
				Data: err.Error(),
			},
		}, err
	}

	candidate := cState.Candidates.GetCandidate(types.BytesToPubkey([]byte(req.PublicKey)))
	if candidate == nil {
		return &pb.CandidateResponse{
			Error: &pb.Error{
				Code:    "404",
				Message: "Candidate not found",
			},
		}, nil
	}

	result := makeResponseCandidate(cState, *candidate, true)
	return &pb.CandidateResponse{Result: result}, nil
}

func makeResponseCandidate(state *state.State, c candidates.Candidate, includeStakes bool) *pb.CandidateResult {
	candidate := &pb.CandidateResult{
		RewardAddress: c.RewardAddress.String(),
		TotalStake:    state.Candidates.GetTotalStake(c.PubKey).String(),
		PublicKey:     c.PubKey.String(),
		Commission:    fmt.Sprintf("%d", c.Commission),
		Status:        fmt.Sprintf("%d", c.Status),
	}

	if includeStakes {
		stakes := state.Candidates.GetStakes(c.PubKey)
		candidate.Stakes = make([]*pb.CandidateResult_Stake, len(stakes))
		for i, stake := range stakes {
			candidate.Stakes[i] = &pb.CandidateResult_Stake{
				Owner:    stake.Owner.String(),
				Coin:     stake.Coin.String(),
				Value:    stake.Value.String(),
				BipValue: stake.BipValue.String(),
			}
		}
	}

	return candidate
}
