package distanceestimator

import (
	"context"
	"testing"

	"github.com/edusalguero/roteiro.git/internal/config"
	"github.com/edusalguero/roteiro.git/internal/logger"
	"github.com/edusalguero/roteiro.git/internal/point"
	"github.com/stretchr/testify/assert"
)

func TestGoogleMapsDistanceEstimator_EstimateRouteDistances(t *testing.T) {
	cnf, err := config.Get()
	if err != nil {
		t.Errorf("Invalid config = %v", err)
	}
	tests := []struct {
		name string
		from point.Point
		to   point.Point
		want *RouteEstimation
	}{
		{
			"As Pontes - Sada",
			point.NewPoint(43.450218, -7.853109),
			point.NewPoint(43.347306, -8.276904),
			&RouteEstimation{
				From:     point.NewPoint(43.450218, -7.853109),
				To:       point.NewPoint(43.347306, -8.276904),
				Distance: 57406,
				Duration: 2988000000000,
			},
		},
	}

	g, err := NewGoogleMapsDistanceEstimator(cnf.DistanceEstimator.GoogleMaps.APIKey, logger.NewNopLogger())
	if err != nil {
		t.Errorf("NewGoogleMapsDistanceEstimator() error = %v", err)
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := g.EstimateDistance(context.TODO(), tt.from, tt.to)
			if err != nil {
				t.Errorf("EstimateDistance() error = %v", err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
