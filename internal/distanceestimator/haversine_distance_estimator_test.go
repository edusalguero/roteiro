package distanceestimator

import (
	"context"
	"reflect"
	"testing"

	"github.com/edusalguero/roteiro.git/internal/cost"
	"github.com/edusalguero/roteiro.git/internal/point"
)

func Test_HaversineDistanceEstimator_GetCost(t *testing.T) {
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
				Distance: 36101,
				Duration: 2599291934724,
			},
		},
	}
	for _, tt := range tests {
		g := NewHaversineDistanceEstimator(50)
		got, err := g.GetCost(context.TODO(), tt.from, tt.to)
		if err != nil {
			t.Errorf("EstimateDistance() error = %v", err)
			return
		}
		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("EstimateDistance() got = %v, want %v", got, tt.want)
		}
	}
}
