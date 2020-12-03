package config

import (
	httpwrapper "github.com/edusalguero/roteiro.git/internal/utils/httpserver"
	"github.com/edusalguero/roteiro.git/internal/utils/shutdown"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Log               Logger
	Shutdown          shutdown.Config
	DistanceEstimator DistanceEstimator
	Server            httpwrapper.Config
}

type Logger struct {
	Level string `default:"debug"`
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
