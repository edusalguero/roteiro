package solver

import (
	"context"
	"testing"

	"github.com/edusalguero/roteiro.git/internal/distanceestimator"
	mock_distanceestimator "github.com/edusalguero/roteiro.git/internal/distanceestimator/mock"
	"github.com/edusalguero/roteiro.git/internal/logger"
	"github.com/edusalguero/roteiro.git/internal/model"
	"github.com/edusalguero/roteiro.git/internal/point"
	"github.com/edusalguero/roteiro.git/internal/problem"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func Test_service_SolveProblem(t *testing.T) {
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
		Load:      1,
	}
	req2 := problem.Request{
		RequestID: "As Pontes 2",
		PickUp:    aspontesLoc,
		DropOff:   sadaLoc,
		Load:      1,
	}
	req3 := problem.Request{
		RequestID: "As Pontes 3",
		PickUp:    aspontesLoc,
		DropOff:   sadaLoc,
		Load:      1,
	}
	req4 := problem.Request{
		RequestID: "As Pontes 4",
		PickUp:    aspontesLoc,
		DropOff:   sadaLoc,
		Load:      1,
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
							RequestID: model.Ref(req1.RequestID),
							PickUp:    req1.PickUp,
							DropOff:   req1.DropOff,
							Load:      1,
						},
						{
							RequestID: model.Ref(req2.RequestID),
							PickUp:    req2.PickUp,
							DropOff:   req2.DropOff,
							Load:      1,
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
									Ref:          model.Ref(req2.RequestID),
								},
								{
									ActivityType: model.ActivityTypePickUp,
									Ref:          model.Ref(req1.RequestID),
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
							RequestID: model.Ref(req3.RequestID),
							PickUp:    req1.PickUp,
							DropOff:   req1.DropOff,
							Load:      1,
						},
						{
							RequestID: model.Ref(req4.RequestID),
							PickUp:    req2.PickUp,
							DropOff:   req2.DropOff,
							Load:      1,
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
									Ref:          model.Ref(req4.RequestID),
								},
								{
									ActivityType: model.ActivityTypePickUp,
									Ref:          model.Ref(req3.RequestID),
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
	s := NewSolver(logger.NewLogger(), e, Config{})

	t.Run("Solve problem", func(t *testing.T) {
		got, err := s.SolveProblem(context.Background(), *p)
		assert.NoError(t, err)
		assert.NotNil(t, got)
		got.Metrics.SolvedTime = SolvedTime
		assert.Equal(t, solution, *got)
	})
}

func Test_service_SolveProblem_WithErrorBuildingMatrix(t *testing.T) {
	var aspontesLoc = point.NewPoint(43.450218, -7.853109)
	var sadaLoc = point.NewPoint(43.347306, -8.276904)

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

	p := problem.NewProblem(
		uuid.New(),
		[]problem.Asset{
			aspontesAsset,
		},
		[]problem.Request{
			req1,
		},
		problem.Constraints{
			MaxJourneyTimeFactor: 1.5,
		})

	t.Run("Solve problem", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		distanceEstimatorMock := mock_distanceestimator.NewMockService(ctrl)
		distanceEstimatorMock.
			EXPECT().
			GetCost(gomock.Any(), gomock.Any(), gomock.Any()).
			Return(nil, ErrBuildingDistanceMatrix).AnyTimes()

		s := NewSolver(logger.NewNopLogger(), distanceEstimatorMock, Config{})

		got, err := s.SolveProblem(context.Background(), *p)
		assert.Error(t, err, ErrBuildingDistanceMatrix)
		assert.Nil(t, got)
	})
}
