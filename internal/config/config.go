package config

import (
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Log               Logger
	DistanceEstimator DistanceEstimator
}
type DistanceEstimator struct {
	GoogleMaps GoogleMaps
}

type GoogleMaps struct {
	APIKey string `required:"true"`
}

type Logger struct {
	Level string `default:"debug"`
}

func Get() (Config, error) {
	var cfg Config
	if err := envconfig.Process("roteiro", &cfg); err != nil {
		return Config{}, err
	}
	return cfg, nil
}
