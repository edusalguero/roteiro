package distanceestimator

import (
	"context"
	"time"

	"github.com/edusalguero/roteiro.git/internal/point"
)

//go:generate mockgen -source=./service.go -destination=./mock/service.go
type Service interface {
	EstimateDistance(ctx context.Context, from point.Point, to point.Point) (*RouteEstimation, error)
}

type RouteEstimation struct {
	From     point.Point
	To       point.Point
	Distance float64
	Duration time.Duration
}
