package solver

import (
	"context"
	"testing"

	"github.com/edusalguero/roteiro.git/internal/algorithms"
	"github.com/edusalguero/roteiro.git/internal/distanceestimator"
	"github.com/edusalguero/roteiro.git/internal/logger"
	"github.com/edusalguero/roteiro.git/internal/model"
	"github.com/edusalguero/roteiro.git/internal/point"
	"github.com/edusalguero/roteiro.git/internal/problem"
	"github.com/edusalguero/roteiro.git/internal/routeestimator"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func Test_service_SolveProblem(t *testing.T) {
	log := logger.NewNopLogger()

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

	p := problem.NewProblem(
		uuid.New(),
		[]problem.Asset{
			minoAsset,
			aspontesAsset,
		},
		[]problem.Request{
			req1,
			req2,
			req3,
			req4,
		},
		problem.Constraints{
			MaxJourneyTimeFactor: 1.5,
		})

	const SolvedTime = 1182235
	solution := problem.Solution{
		ID: p.ID,
		Solution: model.Solution{
			Metrics: model.SolutionMetrics{
				NumAssets:   2,
				NumRequests: 4,
				Duration:    4632548169904,
				Distance:    102945,
				SolvedTime:  SolvedTime,
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
	e := distanceestimator.NewHaversineDistanceEstimator(80)
	routeE := routeestimator.NewEstimator(e)
	algo := algorithms.NewSequentialConstruction(log, routeE, e)
	s := NewSolver(algo)

	t.Run("Solve problem", func(t *testing.T) {
		got, err := s.SolveProblem(context.Background(), *p)
		assert.NoError(t, err)
		assert.NotNil(t, got)
		got.Metrics.SolvedTime = SolvedTime
		assert.Equal(t, solution, *got)
	})
}
