package costmatrix

import (
	"context"
	"fmt"

	"github.com/edusalguero/roteiro.git/internal/cost"
	"github.com/edusalguero/roteiro.git/internal/distanceestimator"
	"github.com/edusalguero/roteiro.git/internal/point"
	"github.com/edusalguero/roteiro.git/internal/problem"
)

var ErrAssetsAreRequired = fmt.Errorf("assets are required")
var ErrRequestsAreRequired = fmt.Errorf("requests are required")

type Builder struct {
	distanceMatrix DistanceMatrix
}

func NewDistanceMatrixBuilder(distanceEstimator distanceestimator.Service) Builder {
	return Builder{
		DistanceMatrix{
			distanceEstimator: distanceEstimator,
			matrix:            make(map[point.Point]map[point.Point]*cost.Cost),
		},
	}
}

func (b Builder) WithAssets(assets []problem.Asset) Builder {
	r := b.distanceMatrix
	r.assets = assets
	return Builder{r}
}

func (b Builder) WithRequests(requests []problem.Request) Builder {
	r := b.distanceMatrix
	r.requests = requests
	return Builder{r}
}

func (b Builder) Build(ctx context.Context) (*DistanceMatrix, error) {
	if b.distanceMatrix.assets == nil {
		return nil, ErrAssetsAreRequired
	}
	if b.distanceMatrix.requests == nil {
		return nil, ErrRequestsAreRequired
	}

	r := b.distanceMatrix
	if err := r.buildMatrix(ctx); err != nil {
		return nil, err
	}
	return &b.distanceMatrix, nil
}
