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
	"github.com/edusalguero/roteiro.git/internal/store"
)

var ErrSavingProblem = fmt.Errorf("error saving problem")
var ErrSavingSolution = fmt.Errorf("error saving solution")
var ErrBuildingDistanceMatrix = fmt.Errorf("error building distance matrix")
var ErrInAlgo = fmt.Errorf("error processing solve algorithm")

//go:generate mockgen -source=./service.go -destination=./mock/service.go
type Service interface {
	SolveProblem(ctx context.Context, p problem.Problem) (*problem.Solution, error)
}

type Solver struct {
	logger            logger.Logger
	cnf               Config
	distanceEstimator distanceestimator.Service
	repository        store.Repository
}

func NewSolver(log logger.Logger, conf Config, r store.Repository, d distanceestimator.Service) *Solver {
	return &Solver{distanceEstimator: d, logger: log, repository: r, cnf: conf}
}

func (s *Solver) SolveProblem(ctx context.Context, p problem.Problem) (*problem.Solution, error) {
	log := s.logger.WithField("problem_id", p.ID)

	err := s.repository.AddProblem(ctx, &p)
	if err != nil {
		log.Errorf("Adding problem to the repository %s", err)
		return nil, ErrSavingProblem
	}

	start := time.Now()
	log.Infof("Building Cost Matrix...")
	matrix, err := costmatrix.NewDistanceMatrixBuilder(s.distanceEstimator, s.logger).
		WithAssets(p.Fleet).
		WithRequests(p.Requests).
		Build(ctx)
	if err != nil {
		log.Errorf("Building Cost Matrix done %s", err)
		if err := s.repository.SetError(ctx, p.ID, err); err != nil {
			return nil, ErrBuildingDistanceMatrix
		}
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
		if err := s.repository.SetError(ctx, p.ID, err); err != nil {
			log.Errorf("Setting error: %s", err)
			return nil, ErrInAlgo
		}
		log.Errorf("Solving algorithm", err)
		return nil, ErrInAlgo
	}

	log.WithField("duration", duration).Infof("Problem solved [%s]", duration)
	solution := &problem.Solution{
		ID:       p.ID,
		Solution: *sol,
	}
	if err := s.repository.SetSolution(ctx, p.ID, solution); err != nil {
		log.Errorf("setting error: %s", err)
		return nil, ErrSavingSolution
	}

	return solution, nil
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
