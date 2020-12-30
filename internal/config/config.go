package config

import (
	"github.com/edusalguero/roteiro.git/internal/distanceestimator"
	"github.com/edusalguero/roteiro.git/internal/logger"
	"github.com/edusalguero/roteiro.git/internal/solver"
	httpwrapper "github.com/edusalguero/roteiro.git/internal/utils/httpserver"
	"github.com/edusalguero/roteiro.git/internal/utils/shutdown"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Log               logger.Config
	Shutdown          shutdown.Config
	DistanceEstimator distanceestimator.Config
	Server            httpwrapper.Config
	Solver            solver.Config
}

func Get() (Config, error) {
	var cfg Config
	if err := envconfig.Process("roteiro", &cfg); err != nil {
		return Config{}, err
	}
	return cfg, nil
}
