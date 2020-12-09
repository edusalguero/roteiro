package costmatrix

import (
	"context"
	"fmt"
	"testing"

	"github.com/edusalguero/roteiro.git/internal/cost"
	"github.com/edusalguero/roteiro.git/internal/distanceestimator"
	"github.com/edusalguero/roteiro.git/internal/point"
	"github.com/edusalguero/roteiro.git/internal/problem"
	"github.com/stretchr/testify/assert"
)

var oneOrigin = point.NewPoint(4.68295, -74.04965)

func TestDistanceMatrix_GetDistance(t *testing.T) {
	type query struct {
		from point.Point
		to   point.Point
	}

	tests := []struct {
		name     string
		assets   []problem.Asset
		requests []problem.Request
		query    query
		cost     *cost.Cost
		wantErr  bool
		err      error
	}{
		{
			name: "Invalid query",
			assets: []problem.Asset{
				{

					AssetID:  "Asset 1",
					Location: point.NewPoint(43.3475, -8.206389),
					Capacity: 2,
				},
			},
			requests: []problem.Request{
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
			query: query{
				from: point.Point{},
				to:   point.Point{},
			},
			cost:    nil,
			wantErr: true,
			err:     ErrPointOutOfMatrix,
		},
		{
			name: "A 20x20 matrix",
			assets: []problem.Asset{
				{
					AssetID:  "Asset 1",
					Location: oneOrigin,
					Capacity: 4,
				},
				{
					AssetID:  "Asset 2",
					Location: oneOrigin,
					Capacity: 4,
				},
				{
					AssetID:  "Asset 3",
					Location: oneOrigin,
					Capacity: 4,
				},
				{
					AssetID:  "Asset 4",
					Location: oneOrigin,
					Capacity: 4,
				},
				{
					AssetID:  "Asset 5",
					Location: oneOrigin,
					Capacity: 4,
				},
				{
					AssetID:  "Asset 6",
					Location: oneOrigin,
					Capacity: 4,
				},
				{
					AssetID:  "Asset 7",
					Location: oneOrigin,
					Capacity: 4,
				},
				{
					AssetID:  "Asset 8",
					Location: oneOrigin,
					Capacity: 4,
				},
			},
			requests: manyRequests(t, oneOrigin),
			query: query{
				from: oneOrigin,
				to:   point.Point{4.56418, -74.08705},
			},
			cost: &cost.Cost{
				Distance: 13842,
				Duration: 622883834142,
			},
			wantErr: false,
			err:     nil,
		},
	}

	e := distanceestimator.NewHaversineDistanceEstimator(80)
	b := NewDistanceMatrixBuilder(e)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dm, err := b.WithRequests(tt.requests).WithAssets(tt.assets).Build(context.Background())
			if err != nil {
				t.Errorf("Builder() error = %v", err)
				return
			}
			if dm == nil {
				t.Errorf("Builder() return nil")
				return
			}
			got, err := dm.GetCost(context.Background(), tt.query.from, tt.query.to)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetDistance() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				assert.Error(t, err, tt.err)
			}
			if err == nil {
				assert.Equal(t, got, tt.cost)
			}
		})
	}
}

func manyRequests(t *testing.T, oneOrigin point.Point) []problem.Request {
	t.Helper()
	var many []problem.Request
	for i, p := range []point.Point{
		{4.56418, -74.08705},
		{4.52924, -74.08894},
		{4.62963, -74.17566},
		{4.65039, -74.17011},
		{4.60634, -74.1465},
		{4.69123, -74.14789},
		{4.56294, -74.06794},
		{4.61886, -74.12805},
		{4.58963, -74.1748},
		{4.62566, -74.13724},
		{4.75886, -74.0906},
		{4.75369, -74.10028},
		{4.75828, -74.10493},
		{4.74706, -74.1123},
		{4.75821, -74.10011},
		{4.57929, -74.1545},
		{4.60265, -74.12565},
		{4.68968, -74.11208},
		{4.70884, -74.11448},
		{4.57419, -74.10678},
		{4.63577, -74.15681},
		{4.72291, -74.13105},
		{4.75796, -74.03786},
		{4.72129, -74.0559},
		{4.6329, -74.06282},
		{4.68548, -74.07004},
	} {
		many = append(many, problem.Request{
			RequestID: problem.RequestID(fmt.Sprintf("Rider %d", i)),
			PickUp:    oneOrigin,
			DropOff:   p,
		})
	}
	return many
}
