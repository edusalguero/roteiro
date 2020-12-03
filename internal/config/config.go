package config

import (
	"github.com/edusalguero/roteiro.git/internal/logger"
	httpwrapper "github.com/edusalguero/roteiro.git/internal/utils/httpserver"
	"github.com/edusalguero/roteiro.git/internal/utils/shutdown"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Log               logger.Config
	Shutdown          shutdown.Config
	DistanceEstimator DistanceEstimator
	Server            httpwrapper.Config
}

type DistanceEstimator struct {
	GoogleMaps GoogleMaps
}

type GoogleMaps struct {
	Enabled bool   `default:"false"`
	APIKey  string `required:"true"`
}

func Get() (Config, error) {
	var cfg Config
	if err := envconfig.Process("roteiro", &cfg); err != nil {
		return Config{}, err
	}
	return cfg, nil
}
