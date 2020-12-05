package distanceestimator

import (
	"context"
	"math"
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
			assert.Equal(t, tt.want.From, got.From)
			assert.Equal(t, tt.want.To, got.To)
			assert.Equal(t, tt.want.Distance, got.Distance)
			// This is a Padron pepper test.
			// Depending on the time of day the test passes or not because the time estimate changes a few seconds...
			// Rounding off the minutes we ensure that the test passes.
			assert.Equal(t, math.RoundToEven(tt.want.Duration.Minutes()), math.RoundToEven(got.Duration.Minutes()))
		})
	}
}
