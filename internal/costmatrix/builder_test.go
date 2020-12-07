package costmatrix

import (
	"context"
	"testing"

	"github.com/edusalguero/roteiro.git/internal/distanceestimator"
	"github.com/edusalguero/roteiro.git/internal/point"
	"github.com/edusalguero/roteiro.git/internal/problem"
	"github.com/stretchr/testify/assert"
)

func TestNewDistanceMatrixBuilder(t *testing.T) {
	e := distanceestimator.NewHaversineDistanceEstimator(80)

	tests := []struct {
		name     string
		assets   []problem.Asset
		requests []problem.Request
		wantErr  bool
		err      error
	}{
		{
			"Without assets",
			nil,
			nil,
			true,
			ErrAssetsAreRequired,
		},
		{
			"Without requests",
			[]problem.Asset{
				{

					AssetID:  "Asset 1",
					Location: point.NewPoint(43.3475, -8.206389),
					Capacity: 2,
				},
			},
			nil,
			true,
			ErrRequestsAreRequired,
		},
		{
			"Without errors",
			[]problem.Asset{
				{

					AssetID:  "Asset 1",
					Location: point.NewPoint(43.3475, -8.206389),
					Capacity: 2,
				},
			},
			[]problem.Request{
				{
					RequestID: "Request 1",
					PickUp:    point.NewPoint(43.450218, -7.853109),
					DropOff:   point.NewPoint(43.3475, -8.206389),
				},
				{
					RequestID: "Request 1",
					PickUp:    point.NewPoint(43.450218, -7.853109),
					DropOff:   point.NewPoint(43.347306, -8.276904),
				},
			},
			false,
			nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			builder := NewDistanceMatrixBuilder(e)
			m, err := builder.WithAssets(tt.assets).WithRequests(tt.requests).Build(context.Background())

			if (err != nil) != tt.wantErr {
				t.Errorf("Build() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr == true && err != tt.err {
				t.Errorf("Build() should have error")
			}

			if m != nil {
				assert.NotEmpty(t, m.matrix)
			}
		})
	}
}
