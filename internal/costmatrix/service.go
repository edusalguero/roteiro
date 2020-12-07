package costmatrix

import (
	"context"
	"fmt"

	"github.com/edusalguero/roteiro.git/internal/cost"
	"github.com/edusalguero/roteiro.git/internal/distanceestimator"
	"github.com/edusalguero/roteiro.git/internal/point"
	"github.com/edusalguero/roteiro.git/internal/problem"
)

var ErrPointOutOfMatrix = fmt.Errorf("point out of matrix")

type Service interface {
	GetCost(ctx context.Context, from point.Point, to point.Point) (*cost.Cost, error)
}

type DistanceMatrix struct {
	distanceEstimator distanceestimator.Service
	assets            []problem.Asset
	requests          []problem.Request
	matrix            map[point.Point]map[point.Point]*cost.Cost
}

func (d *DistanceMatrix) GetCost(_ context.Context, from, to point.Point) (*cost.Cost, error) {
	f, ok := d.matrix[from]
	if !ok {
		return nil, ErrPointOutOfMatrix
	}

	estimation, ok := f[to]
	if !ok {
		return nil, ErrPointOutOfMatrix
	}

	return estimation, nil
}

func (d *DistanceMatrix) buildMatrix(ctx context.Context) error {
	err := d.assets2Requests(ctx)
	if err != nil {
		return err
	}

	err = d.request2Requests(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (d *DistanceMatrix) request2Requests(ctx context.Context) error {
	for _, r1 := range d.requests {
		if _, ok := d.matrix[r1.PickUp]; !ok {
			d.matrix[r1.PickUp] = make(map[point.Point]*cost.Cost)
		}
		if _, ok := d.matrix[r1.DropOff]; !ok {
			d.matrix[r1.DropOff] = make(map[point.Point]*cost.Cost)
		}

		c, err := d.distanceEstimator.GetCost(ctx, r1.PickUp, r1.DropOff)
		if err != nil {
			return err
		}
		d.matrix[r1.PickUp][r1.DropOff] = c

		c, err = d.distanceEstimator.GetCost(ctx, r1.DropOff, r1.PickUp)
		if err != nil {
			return err
		}
		d.matrix[r1.DropOff][r1.PickUp] = c

		for _, r2 := range d.requests {
			c, err := d.distanceEstimator.GetCost(ctx, r1.PickUp, r2.PickUp)
			if err != nil {
				return err
			}
			d.matrix[r1.PickUp][r2.PickUp] = c

			c, err = d.distanceEstimator.GetCost(ctx, r1.PickUp, r2.DropOff)
			if err != nil {
				return err
			}
			d.matrix[r1.PickUp][r2.DropOff] = c

			c, err = d.distanceEstimator.GetCost(ctx, r1.DropOff, r2.PickUp)
			if err != nil {
				return err
			}
			d.matrix[r1.DropOff][r2.PickUp] = c

			c, err = d.distanceEstimator.GetCost(ctx, r1.DropOff, r2.DropOff)
			if err != nil {
				return err
			}
			d.matrix[r1.DropOff][r2.DropOff] = c
		}
	}
	return nil
}

func (d *DistanceMatrix) assets2Requests(ctx context.Context) error {
	for _, a := range d.assets {
		if _, ok := d.matrix[a.Location]; !ok {
			d.matrix[a.Location] = make(map[point.Point]*cost.Cost)
		}

		for _, r := range d.requests {
			c, err := d.distanceEstimator.GetCost(ctx, a.Location, r.PickUp)
			if err != nil {
				return err
			}
			d.matrix[a.Location][r.PickUp] = c

			c, err = d.distanceEstimator.GetCost(ctx, a.Location, r.DropOff)
			if err != nil {
				return err
			}
			d.matrix[a.Location][r.DropOff] = c
		}
	}
	return nil
}
