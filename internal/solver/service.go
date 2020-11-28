package solver

import (
	"context"

	"github.com/edusalguero/roteiro.git/internal/algorithms"
	"github.com/edusalguero/roteiro.git/internal/model"
	"github.com/edusalguero/roteiro.git/internal/problem"
)

type Solver interface {
	SolveProblem(ctx context.Context, problem problem.Problem) (*problem.Solution, error)
}

type Service struct {
	algo algorithms.Algorithm
}

func NewService(algo algorithms.Algorithm) *Service {
	return &Service{algo: algo}
}

func (s *Service) SolveProblem(ctx context.Context, p problem.Problem) (*problem.Solution, error) {
	algoProblem := NewAlgoProblemFromSolverProblem(p)
	sol, err := s.algo.Solve(ctx, algoProblem)

	if err != nil {
		return nil, err
	}
	return &problem.Solution{
		ID:       p.ID.String(),
		Solution: *sol,
	}, nil
}

func NewAlgoProblemFromSolverProblem(p problem.Problem) model.Problem {
	var reqs []model.Request
	for _, req := range p.Requests {
		reqs = append(reqs, model.Request{
			RequestID: model.RequestID(req.RequestID),
			PickUp:    req.PickUp,
			DropOff:   req.DropOff,
		})
	}

	var assets []model.Asset
	for _, asset := range p.Fleet {
		assets = append(assets, model.Asset{
			AssetID:  model.AssetID(asset.AssetID),
			Location: asset.Location,
			Capacity: model.Capacity(asset.Capacity),
		})
	}
	return model.Problem{
		Fleet:    assets,
		Requests: reqs,
		Constraints: model.Constraints{
			MaxJourneyTimeFactor: p.Constraints.MaxJourneyTimeFactor,
		},
	}
}
