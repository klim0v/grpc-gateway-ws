package service

import (
	"context"
	"github.com/klim0v/grpc-gateway-ws/pb"
)

func (s *Service) Candidates(_ context.Context, req *pb.CandidatesRequest) (*pb.CandidatesResponse, error) {
	cState, err := s.getStateForHeight(req.Height)
	if err != nil {
		return &pb.CandidatesResponse{
			Error: &pb.Error{
				Data: err.Error(),
			},
		}, nil
	}

	candidates := cState.Candidates.GetCandidates()

	result := &pb.CandidatesResponse{
		Result: make([]*pb.CandidateResult, len(candidates)),
	}
	for i, candidate := range candidates {
		result.Result[i] = makeResponseCandidate(cState, *candidate, req.IncludeStakes)
	}

	return result, nil
}
