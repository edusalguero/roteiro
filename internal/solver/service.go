package solver

import (
	"context"
	"fmt"
	"time"

	"github.com/edusalguero/roteiro.git/internal/algorithms"
	"github.com/edusalguero/roteiro.git/internal/costmatrix"
	"github.com/edusalguero/roteiro.git/internal/distanceestimator"
	"github.com/edusalguero/roteiro.git/internal/logger"
	"github.com/edusalguero/roteiro.git/internal/model"
	"github.com/edusalguero/roteiro.git/internal/problem"
	"github.com/edusalguero/roteiro.git/internal/routeestimator"
)

var ErrBuildingDistanceMatrix = fmt.Errorf("error building distance matrix")
var ErrInAlgo = fmt.Errorf("error processing solve algorithm")

//go:generate mockgen -source=./service.go -destination=./mock/service.go
type Service interface {
	SolveProblem(ctx context.Context, p problem.Problem) (*problem.Solution, error)
}

type Solver struct {
	distanceEstimator distanceestimator.Service
	logger            logger.Logger
}

func NewSolver(d distanceestimator.Service, log logger.Logger) *Solver {
	return &Solver{distanceEstimator: d, logger: log}
}

func (s *Solver) SolveProblem(ctx context.Context, p problem.Problem) (*problem.Solution, error) {
	log := s.logger.WithField("problem_id", p.ID)

	start := time.Now()
	log.Infof("Building Cost Matrix...")
	matrix, err := costmatrix.NewDistanceMatrixBuilder(s.distanceEstimator, s.logger).
		WithAssets(p.Fleet).
		WithRequests(p.Requests).
		Build(ctx)
	if err != nil {
		log.Errorf("Building Cost Matrix done %s", err)
		return nil, ErrBuildingDistanceMatrix
	}
	duration := time.Since(start)
	log.WithField("duration", duration).Infof("Cost Matrix done [%s]", duration)

	routeE := routeestimator.NewEstimator(matrix)
	algo := algorithms.NewSequentialConstruction(s.logger, routeE, matrix)

	algoProblem := NewAlgoProblemFromSolverProblem(p)
	log.Infof("Solving problem...")
	start = time.Now()
	sol, err := algo.Solve(ctx, algoProblem)
	duration = time.Since(start)
	if err != nil {
		log.Errorf("Solving algorithm", err)
		return nil, ErrInAlgo
	}

	log.WithField("duration", duration).Infof("Problem solved [%s]", duration)
	return &problem.Solution{
		ID:       p.ID,
		Solution: *sol,
	}, nil
}

func NewAlgoProblemFromSolverProblem(p problem.Problem) model.Problem {
	var reqs []model.Request
	for _, req := range p.Requests {
		reqs = append(reqs, model.Request{
			RequestID: model.Ref(req.RequestID),
			PickUp:    req.PickUp,
			DropOff:   req.DropOff,
			Load:      model.Load(req.Load),
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
