package cost

import (
	"context"

	"github.com/edusalguero/roteiro.git/internal/point"
)

type Service interface {
	GetCost(ctx context.Context, from point.Point, to point.Point) (*Cost, error)
}
