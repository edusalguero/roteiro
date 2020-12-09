package costmatrix

import (
	"context"
	"fmt"
	"sync"

	"github.com/edusalguero/roteiro.git/internal/cost"
	"github.com/edusalguero/roteiro.git/internal/distanceestimator"
	"github.com/edusalguero/roteiro.git/internal/logger"
	"github.com/edusalguero/roteiro.git/internal/point"
	"github.com/edusalguero/roteiro.git/internal/problem"
)

var ErrPointOutOfMatrix = fmt.Errorf("point out of matrix")
var ErrBuildingMatrix = fmt.Errorf("error building matrix")

type Service interface {
	GetCost(ctx context.Context, from point.Point, to point.Point) (*cost.Cost, error)
}

type DistanceMatrix struct {
	distanceEstimator distanceestimator.Service
	assets            []problem.Asset
	requests          []problem.Request
	logger            logger.Logger
	matrix            costMap
	errors            []error
}

func (d *DistanceMatrix) GetCost(_ context.Context, from, to point.Point) (*cost.Cost, error) {
	e, ok := d.matrix[costPath{from, to}]
	if !ok {
		return nil, ErrPointOutOfMatrix
	}

	return e, nil
}

func (d *DistanceMatrix) buildMatrix(ctx context.Context) error {
	var lock = sync.Mutex{}
	var points []point.Point
	for _, r := range d.requests {
		points = append(points, r.DropOff)
		points = append(points, r.PickUp)
	}
	points = point.UniquePoints(points)

	limiter := make(chan int, 10000) // This is just a buffer to limit the number of concurrent goroutines
	costs := make(chan costResponse)
	wg := sync.WaitGroup{}

	go d.handleCostResults(costs, &lock, &wg)

	for _, a := range d.assets {
		d.getCostFrom2Points(ctx, costs, limiter, &wg, a.Location, points...)
	}

	for _, p := range points {
		d.getCostFrom2Points(ctx, costs, limiter, &wg, p, points...)
	}

	wg.Wait()
	if len(d.errors) > 0 {
		d.logger.Errorf("Error building Cost MatrixL: %s", d.errors)
		return ErrBuildingMatrix
	}
	return nil
}

func (d *DistanceMatrix) getCostFrom2Points(ctx context.Context, costResults chan costResponse, sem chan int, wg *sync.WaitGroup, from point.Point, to ...point.Point) {
	for _, p := range to {
		sem <- 1
		wg.Add(1)
		go func(p point.Point) {
			c, err := d.distanceEstimator.GetCost(ctx, from, p)
			costMap := costMap{}
			costMap[costPath{from, p}] = c
			costResults <- costResponse{costMap: costMap, err: err}
			<-sem
		}(p)
	}
}

func (d *DistanceMatrix) handleCostResults(costs chan costResponse, lock sync.Locker, wg *sync.WaitGroup) {
	for result := range costs {
		lock.Lock()
		if result.err != nil {
			d.errors = append(d.errors, result.err)
		}
		for path, c := range result.costMap {
			d.matrix[path] = c
		}
		lock.Unlock()
		wg.Done()
	}
}

type costPath [2]point.Point
type costMap map[costPath]*cost.Cost

type costResponse struct {
	costMap costMap
	err     error
}
