package routeestimator

import (
	"context"
	"testing"

	"github.com/edusalguero/roteiro.git/internal/distanceestimator"
	"github.com/edusalguero/roteiro.git/internal/point"
	"github.com/stretchr/testify/assert"
)

func TestEstimator_GetRouteEstimation(t *testing.T) {
	tests := []struct {
		name    string
		points  []point.Point
		want    Estimation
		wantErr bool
	}{
		{"As Pontes - Miño - Sada",
			[]point.Point{
				point.NewPoint(43.450218, -7.853109),
				point.NewPoint(43.347306, -8.276904),
				point.NewPoint(43.3475, -8.206389),
			},
			Estimation{
				Legs: []Leg{
					{
						point.NewPoint(43.450218, -7.853109),
						point.NewPoint(43.347306, -8.276904),
						36101,
						3249114918406,
					},
					{
						point.NewPoint(43.347306, -8.276904),
						point.NewPoint(43.3475, -8.206389),
						5702,
						513179153771,
					},
				},
				TotalDistance: 41803,
				TotalDuration: 3762294072177,
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := NewEstimator(distanceestimator.NewHaversineDistanceEstimator(40))
			got, err := e.GetRouteEstimation(context.Background(), tt.points)
			assert.NoError(t, err)
			assert.Equal(t, *got, tt.want)
		})
	}
}
