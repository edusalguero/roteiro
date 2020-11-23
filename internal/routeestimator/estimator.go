package routeestimator

import (
	"context"
	"time"

	"github.com/edusalguero/roteiro.git/internal/distanceestimator"
	"github.com/edusalguero/roteiro.git/internal/point"
)

type Estimator struct {
	de distanceestimator.Service
}

func NewEstimator(de distanceestimator.Service) Estimator {
	return Estimator{de: de}
}
func (e Estimator) GetRouteEstimation(ctx context.Context, points []point.Point) (*Estimation, error) {

	l := len(points)
	tDistance := 0.0
	tDuration := time.Duration(0)
	var legs []Leg
	for i := range points {
		if i+1 >= l {
			break
		}
		re, err := e.de.EstimateDistance(ctx, points[i], points[i+1])
		if err != nil {
			return nil, err
		}
		legs = append(legs, Leg{
			From:     points[i],
			To:       points[i+1],
			Distance: re.Distance,
			Duration: re.Duration,
		})
		tDistance += re.Distance
		tDuration += re.Duration

	}

	return &Estimation{
		Legs:          legs,
		TotalDistance: tDistance,
		TotalDuration: tDuration,
	}, nil
}

type Leg struct {
	From     point.Point
	To       point.Point
	Distance float64
	Duration time.Duration
}

type Estimation struct {
	Legs          []Leg
	TotalDistance float64
	TotalDuration time.Duration
}
