package distanceestimator

import (
	"context"

	"github.com/edusalguero/roteiro.git/internal/cost"
	"github.com/edusalguero/roteiro.git/internal/point"
)

//go:generate mockgen -source=./service.go -destination=./mock/service.go
type Service interface {
	GetCost(ctx context.Context, from point.Point, to point.Point) (*cost.Cost, error)
}
