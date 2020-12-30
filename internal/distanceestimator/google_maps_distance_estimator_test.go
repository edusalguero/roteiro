package distanceestimator

import (
	"context"
	"math"
	"testing"

	"github.com/edusalguero/roteiro.git/internal/cost"
	"github.com/edusalguero/roteiro.git/internal/logger"
	"github.com/edusalguero/roteiro.git/internal/point"
	"github.com/stretchr/testify/assert"
)

func TestGoogleMapsDistanceEstimator_GetCost(t *testing.T) {
	t.Skip()
	tests := []struct {
		name string
		from point.Point
		to   point.Point
		want *cost.Cost
	}{
		{
			"As Pontes - Sada",
			point.NewPoint(43.450218, -7.853109),
			point.NewPoint(43.347306, -8.276904),
			&cost.Cost{
				Distance: 57406,
				Duration: 2988000000000,
			},
		},
		{
			"Same from and to",
			point.NewPoint(43.450218, -7.853109),
			point.NewPoint(43.450218, -7.853109),
			&cost.Cost{
				Distance: 0,
				Duration: 0,
			},
		},
	}

	g, err := NewGoogleMapsDistanceEstimator(GoogleMapsConf{}, logger.NewNopLogger())
	if err != nil {
		t.Errorf("NewGoogleMapsDistanceEstimator() error = %v", err)
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := g.GetCost(context.TODO(), tt.from, tt.to)
			if err != nil {
				t.Errorf("EstimateDistance() error = %v", err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tt.want.Distance, got.Distance)
			// This is a Padron pepper test.
			// Depending on the time of day the test passes or not because the time estimate changes a few seconds...
			// Rounding off the minutes we ensure that the test passes.
			assert.Equal(t, math.RoundToEven(tt.want.Duration.Minutes()), math.RoundToEven(got.Duration.Minutes()))
		})
	}
}
