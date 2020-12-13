package algorithms

import (
	"context"
	"fmt"
	"testing"

	"github.com/edusalguero/roteiro.git/internal/config"
	"github.com/edusalguero/roteiro.git/internal/distanceestimator"
	"github.com/edusalguero/roteiro.git/internal/logger"
	"github.com/edusalguero/roteiro.git/internal/model"
	"github.com/edusalguero/roteiro.git/internal/point"
	"github.com/edusalguero/roteiro.git/internal/routeestimator"
	"github.com/stretchr/testify/assert"
)

var pontedeumeLoc = point.NewPoint(43.407259, -8.171882)
var minoLoc = point.NewPoint(43.3475, -8.206389)
var aspontesLoc = point.NewPoint(43.450218, -7.853109)
var sadaLoc = point.NewPoint(43.347306, -8.276904)
var pontevedraLoc = point.NewPoint(42.4336114, -8.6475)
var vilalbaLoc = point.NewPoint(43.296272, -7.67861)
var oneOrigin = point.NewPoint(4.68295, -74.04965)

type Route []point.Point

func TestSequentialConstruction_Solve(t *testing.T) {
	cnf, err := config.Get()
	if err != nil {
		t.Errorf("Invalid config = %v", err)
	}

	log, err := logger.New(cnf.Log)
	if err != nil {
		t.Errorf("Invalid loger = %v", err)
	}

	e := distanceestimator.NewHaversineDistanceEstimator(80)
	routeE := routeestimator.NewEstimator(e)

	tests := []struct {
		name       string
		problem    model.Problem
		routes     []Route
		unassigned []model.Request
		wantErr    bool
		skip       bool
	}{
		{
			"Same pickup same dropoff",
			model.Problem{
				Fleet: []model.Asset{
					{
						AssetID:  "Miño Asset",
						Location: minoLoc,
						Capacity: 4,
					},
				},
				Requests: []model.Request{
					{
						RequestID: "As Pontes 1",
						PickUp:    aspontesLoc,
						DropOff:   sadaLoc,
						Load:      1,
					},
					{
						RequestID: "As Pontes 2",
						PickUp:    aspontesLoc,
						DropOff:   sadaLoc,
						Load:      1,
					},
					{
						RequestID: "As Pontes 3",
						PickUp:    aspontesLoc,
						DropOff:   sadaLoc,
						Load:      1,
					},
					{
						RequestID: "As Pontes 4",
						PickUp:    aspontesLoc,
						DropOff:   sadaLoc,
						Load:      1,
					},
				},
				Constraints: model.Constraints{
					MaxJourneyTimeFactor: 1.5,
				},
			},
			[]Route{[]point.Point{minoLoc, aspontesLoc, sadaLoc}},
			[]model.Request{},
			false,
			false,
		},
		{
			"With unassigned requests",
			model.Problem{
				Fleet: []model.Asset{
					{
						AssetID:  "Miño Asset",
						Location: minoLoc,
						Capacity: 2,
					},
				},
				Requests: []model.Request{
					{
						RequestID: "As Pontes 1",
						PickUp:    aspontesLoc,
						DropOff:   sadaLoc,
						Load:      1,
					},
					{
						RequestID: "As Pontes 2",
						PickUp:    aspontesLoc,
						DropOff:   sadaLoc,
						Load:      1,
					},
					{
						RequestID: "As Pontes 3",
						PickUp:    aspontesLoc,
						DropOff:   sadaLoc,
						Load:      1,
					},
					{
						RequestID: "As Pontes 4",
						PickUp:    aspontesLoc,
						DropOff:   sadaLoc,
						Load:      1,
					},
				},
				Constraints: model.Constraints{
					MaxJourneyTimeFactor: 1.5,
				},
			},
			[]Route{[]point.Point{minoLoc, aspontesLoc, sadaLoc}},
			[]model.Request{
				{
					RequestID: "As Pontes 3",
					PickUp:    aspontesLoc,
					DropOff:   sadaLoc,
					Load:      1,
				},
				{
					RequestID: "As Pontes 4",
					PickUp:    aspontesLoc,
					DropOff:   sadaLoc,
					Load:      1,
				},
			},
			false,
			false,
		},
		{
			"Different assets, same pickup and same dropoff",
			model.Problem{
				Fleet: []model.Asset{
					{
						AssetID:  "Miño Asset",
						Location: minoLoc,
						Capacity: 2,
					},
					{
						AssetID:  "As Pontes Asset",
						Location: aspontesLoc,
						Capacity: 2,
					},
				},
				Requests: []model.Request{
					{
						RequestID: "As Pontes 1",
						PickUp:    aspontesLoc,
						DropOff:   sadaLoc,
						Load:      1,
					},
					{
						RequestID: "As Pontes 2",
						PickUp:    aspontesLoc,
						DropOff:   sadaLoc,
						Load:      1,
					},
					{
						RequestID: "As Pontes 3",
						PickUp:    aspontesLoc,
						DropOff:   sadaLoc,
						Load:      1,
					},
					{
						RequestID: "As Pontes 4",
						PickUp:    aspontesLoc,
						DropOff:   sadaLoc,
						Load:      1,
					},
				},

				Constraints: model.Constraints{
					MaxJourneyTimeFactor: 1.5,
				},
			},
			[]Route{
				[]point.Point{minoLoc, aspontesLoc, sadaLoc},
				[]point.Point{aspontesLoc, sadaLoc},
			},
			[]model.Request{},
			false,
			false,
		},
		{
			"Insufficient capacity",
			model.Problem{
				Fleet: []model.Asset{
					{
						AssetID:  "As Pontes Asset",
						Location: aspontesLoc,
						Capacity: 2,
					},
					{
						AssetID:  "Miño Asset",
						Location: minoLoc,
						Capacity: 1,
					},
				},
				Requests: []model.Request{
					{
						RequestID: "As Pontes 1",
						PickUp:    aspontesLoc,
						DropOff:   sadaLoc,
						Load:      1,
					},
					{
						RequestID: "As Pontes 2",
						PickUp:    aspontesLoc,
						DropOff:   sadaLoc,
						Load:      1,
					},
					{
						RequestID: "As Pontes 3",
						PickUp:    aspontesLoc,
						DropOff:   sadaLoc,
						Load:      1,
					},
					{
						RequestID: "As Pontes 4",
						PickUp:    aspontesLoc,
						DropOff:   sadaLoc,
						Load:      1,
					},
				},

				Constraints: model.Constraints{
					MaxJourneyTimeFactor: 1.5,
				},
			},
			[]Route{
				[]point.Point{aspontesLoc, sadaLoc},
				[]point.Point{minoLoc, aspontesLoc, sadaLoc},
			},
			[]model.Request{{
				RequestID: "As Pontes 4",
				PickUp:    aspontesLoc,
				DropOff:   sadaLoc,
				Load:      1,
			}},
			false,
			false,
		},
		{
			"Insufficient capacity by a given request",
			model.Problem{
				Fleet: []model.Asset{
					{
						AssetID:  "As Pontes Asset",
						Location: aspontesLoc,
						Capacity: 2,
					},
					{
						AssetID:  "Miño Asset",
						Location: minoLoc,
						Capacity: 1,
					},
				},
				Requests: []model.Request{
					{
						RequestID: "As Pontes 1",
						PickUp:    aspontesLoc,
						DropOff:   sadaLoc,
						Load:      1,
					},
					{
						RequestID: "As Pontes 2",
						PickUp:    aspontesLoc,
						DropOff:   sadaLoc,
						Load:      1,
					},
					{
						RequestID: "As Pontes 3",
						PickUp:    aspontesLoc,
						DropOff:   sadaLoc,
						Load:      3,
					},
					{
						RequestID: "As Pontes 4",
						PickUp:    aspontesLoc,
						DropOff:   sadaLoc,
						Load:      1,
					},
				},

				Constraints: model.Constraints{
					MaxJourneyTimeFactor: 1.5,
				},
			},
			[]Route{
				[]point.Point{aspontesLoc, sadaLoc},
				[]point.Point{minoLoc, aspontesLoc, sadaLoc},
			},
			[]model.Request{{
				RequestID: "As Pontes 3",
				PickUp:    aspontesLoc,
				DropOff:   sadaLoc,
				Load:      3,
			}},
			false,
			false,
		},
		{
			"Different capacity assets, same pickup and same dropoff",
			model.Problem{
				Fleet: []model.Asset{
					{
						AssetID:  "As Pontes Asset",
						Location: aspontesLoc,
						Capacity: 3,
					},
					{
						AssetID:  "Miño Asset",
						Location: minoLoc,
						Capacity: 2,
					},
				},
				Requests: []model.Request{
					{
						RequestID: "As Pontes 1",
						PickUp:    aspontesLoc,
						DropOff:   sadaLoc,
						Load:      1,
					},
					{
						RequestID: "As Pontes 2",
						PickUp:    aspontesLoc,
						DropOff:   sadaLoc,
						Load:      1,
					},
					{
						RequestID: "As Pontes 3",
						PickUp:    aspontesLoc,
						DropOff:   sadaLoc,
						Load:      1,
					},
					{
						RequestID: "As Pontes 4",
						PickUp:    aspontesLoc,
						DropOff:   sadaLoc,
						Load:      1,
					},
				},

				Constraints: model.Constraints{
					MaxJourneyTimeFactor: 1.5,
				},
			},
			[]Route{
				[]point.Point{aspontesLoc, sadaLoc},
				[]point.Point{minoLoc, aspontesLoc, sadaLoc},
			},
			[]model.Request{},
			false,
			false,
		},
		{
			"Asset start where finish",
			model.Problem{
				Fleet: []model.Asset{
					{
						AssetID:  "Sada Asset",
						Location: sadaLoc,
						Capacity: 4,
					},
					{
						AssetID:  "Sada Asset",
						Location: sadaLoc,
						Capacity: 4,
					},
					{
						AssetID:  "Sada Asset",
						Location: sadaLoc,
						Capacity: 4,
					},
				},
				Requests: []model.Request{
					{
						RequestID: "Pontevedra - Sada",
						PickUp:    pontevedraLoc, // Pontevedra
						DropOff:   sadaLoc,       // Sada
						Load:      1,
					},
					{
						RequestID: "Vilalba - Sada",
						PickUp:    vilalbaLoc, // Vilalba
						DropOff:   sadaLoc,    // Sada
						Load:      1,
					},
					{
						RequestID: "As Pontes - Sada",
						PickUp:    aspontesLoc, // As Pontes
						DropOff:   sadaLoc,     // Sada
						Load:      1,
					},
					{
						RequestID: "As Pontes - Miño",
						PickUp:    aspontesLoc, // As Pontes
						DropOff:   minoLoc,     // Miño
						Load:      1,
					},
				},
				Constraints: model.Constraints{
					MaxJourneyTimeFactor: 1.5,
				},
			},
			[]Route{
				[]point.Point{sadaLoc, aspontesLoc, minoLoc, sadaLoc},
				[]point.Point{sadaLoc, pontevedraLoc, sadaLoc},
				[]point.Point{sadaLoc, vilalbaLoc, sadaLoc},
			},
			[]model.Request{},
			false,
			false,
		},
		{
			"Similar points reverse order",
			model.Problem{
				Fleet: []model.Asset{
					{
						AssetID:  "Pontedeume Asset",
						Location: pontedeumeLoc,
						Capacity: 4,
					},
				},
				Requests: []model.Request{
					{
						RequestID: "As Pontes - Sada",
						PickUp:    aspontesLoc, // As Pontes
						DropOff:   sadaLoc,     // Sada
						Load:      1,
					},
					{
						RequestID: "Miño - As Pontes",
						PickUp:    minoLoc,     // Miño
						DropOff:   aspontesLoc, // As Pontes
						Load:      1,
					},
				},
				Constraints: model.Constraints{
					MaxJourneyTimeFactor: 1.5,
				},
			},
			[]Route{
				[]point.Point{pontedeumeLoc, minoLoc, aspontesLoc, sadaLoc},
			},
			[]model.Request{},
			false,
			false,
		},
		{
			"Routific example",
			model.Problem{
				Fleet: []model.Asset{
					{
						AssetID:  "Asset",
						Location: point.NewPoint(49.2553636, -123.0873365),
						Capacity: 4,
					},
				},
				Requests: []model.Request{
					{
						RequestID: "Order 1",
						PickUp:    point.NewPoint(49.227107, -123.1163085),
						DropOff:   point.NewPoint(49.2474624, -123.1532338),
						Load:      1,
					},
					{
						RequestID: "Order 2",
						PickUp:    point.NewPoint(49.2474624, -123.1532338),
						DropOff:   point.NewPoint(49.287107, -122.1163085),
						Load:      1,
					},
				},
				Constraints: model.Constraints{
					MaxJourneyTimeFactor: 1.5,
				},
			},
			[]Route{
				[]point.Point{
					point.NewPoint(49.2553636, -123.0873365), // Depot
					point.NewPoint(49.227107, -123.1163085),  // Order 1 PickUP
					point.NewPoint(49.2474624, -123.1532338), // Order 1 DropOff - Order 2 PickUp
					point.NewPoint(49.287107, -122.1163085),  // Order 2 Drop off
				},
			},
			[]model.Request{},
			false,
			false,
		},
		{
			"One to many",
			model.Problem{
				Fleet: []model.Asset{
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
				Requests: manyRequests(t, oneOrigin),
				Constraints: model.Constraints{
					MaxJourneyTimeFactor: 1.5,
				},
			},
			[]Route{
				[]point.Point{
					point.NewPoint(4.682950, -74.049650),
					point.NewPoint(4.625660, -74.137240),
					point.NewPoint(4.606340, -74.146500),
					point.NewPoint(4.579290, -74.154500),
					point.NewPoint(4.529240, -74.088940),
				},
				[]point.Point{
					point.NewPoint(4.682950, -74.049650),
					point.NewPoint(4.635770, -74.156810),
					point.NewPoint(4.650390, -74.170110),
					point.NewPoint(4.629630, -74.175660),
					point.NewPoint(4.589630, -74.174800),
				},
				[]point.Point{
					point.NewPoint(4.682950, -74.049650),
					point.NewPoint(4.632900, -74.062820),
					point.NewPoint(4.562940, -74.067940),
					point.NewPoint(4.574190, -74.106780),
					point.NewPoint(4.564180, -74.087050),
				},
				[]point.Point{
					point.NewPoint(4.682950, -74.049650),
					point.NewPoint(4.685480, -74.070040),
					point.NewPoint(4.689680, -74.112080),
					point.NewPoint(4.618860, -74.128050),
					point.NewPoint(4.602650, -74.125650),
				},
				[]point.Point{
					point.NewPoint(4.682950, -74.049650),
					point.NewPoint(4.708840, -74.114480),
					point.NewPoint(4.722910, -74.131050),
					point.NewPoint(4.691230, -74.147890),
				},
				[]point.Point{
					point.NewPoint(4.682950, -74.049650),
					point.NewPoint(4.753690, -74.100280),
					point.NewPoint(4.747060, -74.112300),
					point.NewPoint(4.758210, -74.100110),
					point.NewPoint(4.758280, -74.104930),
				},
				[]point.Point{
					point.NewPoint(4.682950, -74.049650),
					point.NewPoint(4.757960, -74.037860),
					point.NewPoint(4.758860, -74.090600),
				},
				[]point.Point{
					point.NewPoint(4.682950, -74.049650),
					point.NewPoint(4.721290, -74.055900),
				},
			},
			[]model.Request{},
			false,
			false,
		},
	}

	algo := NewSequentialConstruction(log, routeE, e)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := algo.Solve(context.Background(), tt.problem)
			if (err != nil) != tt.wantErr {
				t.Errorf("Solve() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.NotNil(t, got)
			if got == nil {
				t.Errorf("Solve() is nil")
				return
			}
			assert.Equal(t, tt.routes, getTestRoutes(t, got.Routes))
			assert.Equal(t, tt.unassigned, got.Unassigned)
		})
	}
}

func getTestRoutes(t *testing.T, routes []model.SolutionRoute) []Route {
	t.Helper()
	var testRoutes []Route
	for _, r := range routes {
		var route Route
		for _, p := range r.Waypoints {
			route = append(route, p.Location)
		}
		testRoutes = append(testRoutes, route)
	}
	return testRoutes
}

func manyRequests(t *testing.T, oneOrigin point.Point) []model.Request {
	t.Helper()
	var many []model.Request
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
		many = append(many, model.Request{
			RequestID: model.Ref(fmt.Sprintf("Rider %d", i)),
			PickUp:    oneOrigin,
			DropOff:   p,
			Load:      1,
		})
	}
	return many
}

func Test_buildRouteWaypoints(t *testing.T) {
	tests := []struct {
		name      string
		asset     model.Asset
		route     model.Route
		waypoints []model.Waypoint
	}{
		{
			"Multiple points same location",
			model.Asset{
				AssetID:  "Pontedeume Asset",
				Location: pontedeumeLoc,
				Capacity: 4,
			},
			model.Route{
				&model.Stop{
					Ref:      "Pontedeume Asset",
					Point:    pontedeumeLoc,
					Load:     0,
					Activity: model.ActivityTypeStart,
				},
				&model.Stop{
					Ref:      "As Pontes - Sada",
					Point:    aspontesLoc,
					Load:     1,
					Activity: model.ActivityTypePickUp,
				},
				&model.Stop{
					Ref:      "As Pontes - Miño",
					Point:    aspontesLoc,
					Load:     1,
					Activity: model.ActivityTypePickUp,
				},
				&model.Stop{
					Ref:      "As Pontes - Miño",
					Point:    minoLoc,
					Load:     -1,
					Activity: model.ActivityTypeDropOff,
				},
				&model.Stop{
					Ref:      "As Pontes - Sada",
					Point:    sadaLoc,
					Load:     -1,
					Activity: model.ActivityTypeDropOff,
				},
			},
			[]model.Waypoint{
				{
					Location: pontedeumeLoc,
					Load:     0,
					Activities: []model.Activity{
						model.NewActivity(model.ActivityTypeStart, "Pontedeume Asset"),
					},
				},
				{
					Location: aspontesLoc,
					Load:     2,
					Activities: []model.Activity{
						model.NewActivity(model.ActivityTypePickUp, "As Pontes - Sada"),
						model.NewActivity(model.ActivityTypePickUp, "As Pontes - Miño"),
					},
				},
				{
					Location: minoLoc,
					Load:     1,
					Activities: []model.Activity{
						model.NewActivity(model.ActivityTypeDropOff, "As Pontes - Miño"),
					},
				},
				{
					Location: sadaLoc,
					Load:     0,
					Activities: []model.Activity{
						model.NewActivity(model.ActivityTypeDropOff, "As Pontes - Sada"),
					},
				},
			},
		},
		{
			"Loop",
			model.Asset{
				AssetID:  "Pontedeume Asset",
				Location: pontedeumeLoc,
				Capacity: 4,
			},
			model.Route{
				&model.Stop{
					Ref:      "Pontedeume Asset",
					Point:    pontedeumeLoc,
					Load:     0,
					Activity: model.ActivityTypeStart,
				},
				&model.Stop{
					Ref:      "As Pontes - Pontedeume",
					Point:    aspontesLoc,
					Load:     1,
					Activity: model.ActivityTypePickUp,
				},
				&model.Stop{
					Ref:      "As Pontes - Miño",
					Point:    aspontesLoc,
					Load:     1,
					Activity: model.ActivityTypePickUp,
				},
				&model.Stop{
					Ref:      "As Pontes - Pontedeume",
					Point:    pontedeumeLoc,
					Load:     -1,
					Activity: model.ActivityTypeDropOff,
				},
				&model.Stop{
					Ref:      "As Pontes - Miño",
					Point:    minoLoc,
					Load:     -1,
					Activity: model.ActivityTypeDropOff,
				},
			},
			[]model.Waypoint{
				{
					Location: pontedeumeLoc,
					Load:     0,
					Activities: []model.Activity{
						model.NewActivity(model.ActivityTypeStart, "Pontedeume Asset"),
					},
				},
				{
					Location: aspontesLoc,
					Load:     2,
					Activities: []model.Activity{
						model.NewActivity(model.ActivityTypePickUp, "As Pontes - Pontedeume"),
						model.NewActivity(model.ActivityTypePickUp, "As Pontes - Miño"),
					},
				},
				{
					Location: pontedeumeLoc,
					Load:     1,
					Activities: []model.Activity{
						model.NewActivity(model.ActivityTypeDropOff, "As Pontes - Pontedeume"),
					},
				},
				{
					Location: minoLoc,
					Load:     0,
					Activities: []model.Activity{
						model.NewActivity(model.ActivityTypeDropOff, "As Pontes - Miño"),
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := buildRouteWaypoints(tt.route, tt.asset)
			assert.Equal(t, tt.waypoints, got)
		})
	}
}
