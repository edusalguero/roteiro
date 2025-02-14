package main

import (
	"github.com/edusalguero/roteiro.git/internal/config"
	"github.com/edusalguero/roteiro.git/internal/distanceestimator"
	"github.com/edusalguero/roteiro.git/internal/logger"
	"github.com/edusalguero/roteiro.git/internal/roteiro"
	"github.com/edusalguero/roteiro.git/internal/solver"
	"github.com/edusalguero/roteiro.git/internal/store"
	httpwrapper "github.com/edusalguero/roteiro.git/internal/utils/httpserver"
	"github.com/edusalguero/roteiro.git/internal/utils/shutdown"
)

func main() {
	cnf, err := config.Get()
	if err != nil {
		panic("could not load config: " + err.Error())
	}
	defer shutdown.Gracefully(cnf.Shutdown)

	log, err := logger.New(cnf.Log)
	if err != nil {
		panic("could not initialize logger: " + err.Error())
	}

	e := distanceestimator.NewHaversineDistanceEstimator(30)
	if cnf.DistanceEstimator.GoogleMaps.Enabled {
		e, err = distanceestimator.NewGoogleMapsDistanceEstimator(
			cnf.DistanceEstimator.GoogleMaps,
			log,
		)
		if err != nil {
			log.Panicf("NewGoogleMapsDistanceEstimator() error = %v", err)
		}
	}

	problemRepo := store.NewInMemoryRepository()
	solverService := solver.NewSolver(log, cnf.Solver, problemRepo, e)
	httpServerWrapper := httpwrapper.NewHTTPServerWrapper(cnf.Server)
	httpServerWrapper.AddController(roteiro.NewStatusController())
	httpServerWrapper.AddController(roteiro.NewSolverController(log, solverService, roteiro.IDGenerator))
	httpServerWrapper.AddController(roteiro.NewProblemController(log, problemRepo))

	log.Info("Starting Roteiro API Server")
	shutdown.First().AfterStarting(httpServerWrapper)
	shutdown.WaitForStopSignal()
}
