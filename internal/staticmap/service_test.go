package staticmap

import (
	"testing"

	"github.com/edusalguero/roteiro.git/internal/model"
	"github.com/edusalguero/roteiro.git/internal/point"
	"github.com/edusalguero/roteiro.git/internal/problem"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestService_Render(t *testing.T) {
	var minoLoc = point.NewPoint(43.3475, -8.206389)
	var aspontesLoc = point.NewPoint(43.450218, -7.853109)
	var sadaLoc = point.NewPoint(43.347306, -8.276904)

	minoAsset := problem.Asset{
		AssetID:  "Miño Asset",
		Location: minoLoc,
		Capacity: 2,
	}
	aspontesAsset := problem.Asset{
		AssetID:  "As Pontes Asset",
		Location: aspontesLoc,
		Capacity: 2,
	}
	req1 := problem.Request{
		RequestID: "As Pontes 1",
		PickUp:    aspontesLoc,
		DropOff:   sadaLoc,
	}
	req2 := problem.Request{
		RequestID: "As Pontes 2",
		PickUp:    aspontesLoc,
		DropOff:   sadaLoc,
	}
	req3 := problem.Request{
		RequestID: "As Pontes 3",
		PickUp:    aspontesLoc,
		DropOff:   sadaLoc,
	}
	req4 := problem.Request{
		RequestID: "As Pontes 4",
		PickUp:    aspontesLoc,
		DropOff:   sadaLoc,
	}
	solution := problem.Solution{
		ID: uuid.New(),
		Solution: model.Solution{
			Metrics: model.SolutionMetrics{
				NumAssets:   2,
				NumRequests: 4,
				Duration:    4632548169904,
				Distance:    102945,
				SolvedTime:  1,
			},
			Routes: []model.SolutionRoute{
				{
					Asset: model.Asset{
						AssetID:  model.AssetID(minoAsset.AssetID),
						Location: minoAsset.Location,
						Capacity: model.Capacity(minoAsset.Capacity),
					},
					Requests: []model.Request{
						{
							RequestID: model.RequestID(req1.RequestID),
							PickUp:    req1.PickUp,
							DropOff:   req1.DropOff,
						},
						{
							RequestID: model.RequestID(req2.RequestID),
							PickUp:    req2.PickUp,
							DropOff:   req2.DropOff,
						},
					},
					Waypoints: []model.Waypoint{
						{
							Location: minoLoc,
							Load:     0,
							Activities: []model.Activity{
								{
									ActivityType: model.ActivityTypeStart,
									Ref:          model.Ref(minoAsset.AssetID),
								},
							},
						},
						{
							Location: aspontesLoc,
							Load:     2,
							Activities: []model.Activity{
								{
									ActivityType: model.ActivityTypePickUp,
									Ref:          model.Ref(req1.RequestID),
								},
								{
									ActivityType: model.ActivityTypePickUp,
									Ref:          model.Ref(req2.RequestID),
								},
							},
						},
						{
							Location: sadaLoc,
							Load:     0,
							Activities: []model.Activity{
								{
									ActivityType: model.ActivityTypeDropOff,
									Ref:          model.Ref(req1.RequestID),
								},
								{
									ActivityType: model.ActivityTypeDropOff,
									Ref:          model.Ref(req2.RequestID),
								},
							},
						},
					},
					Metrics: model.RouteMetrics{
						Duration: 3007990710701,
						Distance: 66844,
					},
				},
				{
					Asset: model.Asset{
						AssetID:  model.AssetID(aspontesAsset.AssetID),
						Location: aspontesAsset.Location,
						Capacity: model.Capacity(aspontesAsset.Capacity),
					},
					Requests: []model.Request{
						{
							RequestID: model.RequestID(req3.RequestID),
							PickUp:    req1.PickUp,
							DropOff:   req1.DropOff,
						},
						{
							RequestID: model.RequestID(req4.RequestID),
							PickUp:    req2.PickUp,
							DropOff:   req2.DropOff,
						},
					},
					Waypoints: []model.Waypoint{
						{
							Location: aspontesLoc,
							Load:     2,
							Activities: []model.Activity{
								{
									ActivityType: model.ActivityTypeStart,
									Ref:          model.Ref(aspontesAsset.AssetID),
								},
								{
									ActivityType: model.ActivityTypePickUp,
									Ref:          model.Ref(req3.RequestID),
								},
								{
									ActivityType: model.ActivityTypePickUp,
									Ref:          model.Ref(req4.RequestID),
								},
							},
						},
						{
							Location: sadaLoc,
							Load:     0,
							Activities: []model.Activity{
								{
									ActivityType: model.ActivityTypeDropOff,
									Ref:          model.Ref(req3.RequestID),
								},
								{
									ActivityType: model.ActivityTypeDropOff,
									Ref:          model.Ref(req4.RequestID),
								},
							},
						},
					},
					Metrics: model.RouteMetrics{
						Duration: 1624557459203,
						Distance: 36101,
					},
				},
			},
			Unassigned: []model.Request{},
		},
	}

	tests := []struct {
		name     string
		solution problem.Solution
		wantErr  bool
	}{
		{
			"One example",
			solution,
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := StaticMap{}
			got, err := s.Render(&tt.solution)
			if (err != nil) != tt.wantErr {
				t.Errorf("Render() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.NotNil(t, got)
		})
	}
}
