package solver

import (
	"context"
	"fmt"

	"github.com/edusalguero/roteiro.git/internal/algorithms"
	"github.com/edusalguero/roteiro.git/internal/model"
	"github.com/edusalguero/roteiro.git/internal/problem"
)

var ErrInAlgo = fmt.Errorf("error processing solve algorithm")

//go:generate mockgen -source=./service.go -destination=./mock/service.go
type Service interface {
	SolveProblem(ctx context.Context, p problem.Problem) (*problem.Solution, error)
}

type Solver struct {
	algo algorithms.Algorithm
}

func NewSolver(algo algorithms.Algorithm) *Solver {
	return &Solver{algo: algo}
}

func (s *Solver) SolveProblem(ctx context.Context, p problem.Problem) (*problem.Solution, error) {
	algoProblem := NewAlgoProblemFromSolverProblem(p)
	sol, err := s.algo.Solve(ctx, algoProblem)

	if err != nil {
		return nil, ErrInAlgo
	}
	return &problem.Solution{
		ID:       p.ID,
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
