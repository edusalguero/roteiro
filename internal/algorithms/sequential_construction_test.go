package algorithms

import (
	"context"
	"fmt"
	"testing"

	"github.com/edusalguero/roteiro.git/internal/config"
	"github.com/edusalguero/roteiro.git/internal/distanceestimator"
	"github.com/edusalguero/roteiro.git/internal/logger"
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
		name    string
		problem Problem
		want    []SolutionRoute
		wantErr bool
		skip    bool
	}{
		{
			"Same pickup same dropoff",
			Problem{
				Fleet: []Asset{
					{
						AssetID:  "Miño Asset",
						Location: minoLoc,
						Capacity: 4,
					},
				},
				Requests: []RiderRequest{
					{
						RiderID: "As Pontes 1",
						PickUp:  aspontesLoc,
						DropOff: sadaLoc,
					},
					{
						RiderID: "As Pontes 2",
						PickUp:  aspontesLoc,
						DropOff: sadaLoc,
					},
					{
						RiderID: "As Pontes 3",
						PickUp:  aspontesLoc,
						DropOff: sadaLoc,
					},
					{
						RiderID: "As Pontes 4",
						PickUp:  aspontesLoc,
						DropOff: sadaLoc,
					},
				},

				Constraints: Constraints{
					MaxJourneyTimeFactor: 1.5,
				},
			},
			[]SolutionRoute{
				{
					Route: []point.Point{minoLoc, aspontesLoc, sadaLoc},
					Metrics: RouteMetrics{
						Duration: 3007990710701,
						Distance: 66844,
					},
				},
			},
			false,
			false,
		},
		{
			"Simple routing",
			Problem{
				Fleet: []Asset{
					{
						AssetID:  "Pontedeume Asset",
						Location: pontedeumeLoc,
						Capacity: 4,
					},
				},
				Requests: []RiderRequest{
					{
						RiderID: "Pontevedra - Sada",
						PickUp:  pontevedraLoc, // Pontevedra
						DropOff: sadaLoc,       // Sada
					},
					{
						RiderID: "Vilalba - Sada",
						PickUp:  vilalbaLoc, // Vilalba
						DropOff: sadaLoc,    // Sada
					},
					{
						RiderID: "As Pontes - Sada",
						PickUp:  aspontesLoc, // As Pontes
						DropOff: sadaLoc,     // Sada
					},
					{
						RiderID: "As Pontes - Miño",
						PickUp:  aspontesLoc, // As Pontes
						DropOff: minoLoc,     // Miño
					},
				},
				Constraints: Constraints{
					MaxJourneyTimeFactor: 1.5,
				},
			},
			[]SolutionRoute{
				{
					Route: []point.Point{pontedeumeLoc, pontevedraLoc, sadaLoc},
					Metrics: RouteMetrics{
						Duration: 9943722533515,
						Distance: 220972,
					},
				},
				{
					Route: []point.Point{pontedeumeLoc, vilalbaLoc, sadaLoc},
					Metrics: RouteMetrics{
						Duration: 4071613441518,
						Distance: 90480,
					},
				},
				{
					Route: []point.Point{pontedeumeLoc, aspontesLoc, minoLoc, sadaLoc},
					Metrics: RouteMetrics{
						Duration: 2818183252084,
						Distance: 62626,
					},
				},
			},
			false,
			false,
		},
		{
			"Routific example",
			Problem{
				Fleet: []Asset{
					{
						AssetID:  "Asset",
						Location: point.NewPoint(49.2553636, -123.0873365),
						Capacity: 4,
					},
				},
				Requests: []RiderRequest{
					{
						RiderID: "Order 1",
						PickUp:  point.NewPoint(49.227107, -123.1163085),
						DropOff: point.NewPoint(49.2474624, -123.1532338),
					},
					{
						RiderID: "Order 2",
						PickUp:  point.NewPoint(49.2474624, -123.1532338),
						DropOff: point.NewPoint(49.287107, -122.1163085),
					},
				},
				Constraints: Constraints{
					MaxJourneyTimeFactor: 1.5,
				},
			},
			[]SolutionRoute{
				{
					Route: []point.Point{
						point.NewPoint(49.2553636, -123.0873365), // Depot
						point.NewPoint(49.227107, -123.1163085),  // Order 1 PickUP
						point.NewPoint(49.2474624, -123.1532338), // Order 1 DropOff - Order 2 PickUp
						point.NewPoint(49.287107, -122.1163085),  // Order 2 Drop off
					},
					Metrics: RouteMetrics{
						Duration: 3719492063273,
						Distance: 82656,
					},
				},
			},
			false,
			false,
		},
		{
			"One to many",
			Problem{
				Fleet: []Asset{
					{
						AssetID:  "Asset",
						Location: oneOrigin,
						Capacity: 4,
					},
				},
				Requests: manyRequests(t, oneOrigin),
				Constraints: Constraints{
					MaxJourneyTimeFactor: 1.5,
				},
			},
			[]SolutionRoute{
				{
					Route: []point.Point{
						point.NewPoint(4.682950, -74.049650),
						point.NewPoint(4.625660, -74.137240),
						point.NewPoint(4.606340, -74.146500),
						point.NewPoint(4.579290, -74.154500),
						point.NewPoint(4.529240, -74.088940),
					},
					Metrics: RouteMetrics{
						Duration: 1182639057395,
						Distance: 26281,
					},
				},
				{
					Route: []point.Point{
						point.NewPoint(4.682950, -74.049650),
						point.NewPoint(4.635770, -74.156810),
						point.NewPoint(4.650390, -74.170110),
						point.NewPoint(4.629630, -74.175660),
						point.NewPoint(4.589630, -74.174800),
					},
					Metrics: RouteMetrics{
						Duration: 990701924097,
						Distance: 22015,
					},
				},
				{
					Route: []point.Point{
						point.NewPoint(4.682950, -74.049650),
						point.NewPoint(4.632900, -74.062820),
						point.NewPoint(4.562940, -74.067940),
						point.NewPoint(4.574190, -74.106780),
						point.NewPoint(4.564180, -74.087050),
					},
					Metrics: RouteMetrics{
						Duration: 922068500653,
						Distance: 20491,
					},
				},
				{
					Route: []point.Point{
						point.NewPoint(4.682950, -74.049650),
						point.NewPoint(4.685480, -74.070040),
						point.NewPoint(4.689680, -74.112080),
						point.NewPoint(4.618860, -74.128050),
						point.NewPoint(4.602650, -74.125650),
					},
					Metrics: RouteMetrics{
						Duration: 758373847392,
						Distance: 16852,
					},
				},
				{
					Route: []point.Point{
						point.NewPoint(4.682950, -74.049650),
						point.NewPoint(4.708840, -74.114480),
						point.NewPoint(4.722910, -74.131050),
						point.NewPoint(4.691230, -74.147890),
					},
					Metrics: RouteMetrics{
						Duration: 636242010369,
						Distance: 14138,
					},
				},
				{
					Route: []point.Point{
						point.NewPoint(4.682950, -74.049650),
						point.NewPoint(4.753690, -74.100280),
						point.NewPoint(4.747060, -74.112300),
						point.NewPoint(4.758210, -74.100110),
						point.NewPoint(4.758280, -74.104930),
					},
					Metrics: RouteMetrics{
						Duration: 609840857983,
						Distance: 13552,
					},
				},
				{
					Route: []point.Point{
						point.NewPoint(4.682950, -74.049650),
						point.NewPoint(4.757960, -74.037860),
						point.NewPoint(4.758860, -74.090600),
					},
					Metrics: RouteMetrics{
						Duration: 642937858539,
						Distance: 14287,
					},
				},
				{
					Route: []point.Point{
						point.NewPoint(4.682950, -74.049650),
						point.NewPoint(4.721290, -74.055900),
					},
					Metrics: RouteMetrics{
						Duration: 194360019503,
						Distance: 4319,
					},
				},
			},
			false,
			false,
		},
	}

	algo := NewSequentialConstruction(log, routeE, e)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.skip {
				t.Skip()
			}
			got, err := algo.Solve(context.Background(), tt.problem)
			if (err != nil) != tt.wantErr {
				t.Errorf("Solve() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.NotNil(t, got)
			assert.Equal(t, tt.want, got.Routes)
		})
	}
}

func manyRequests(t *testing.T, oneOrigin point.Point) []RiderRequest {
	t.Helper()
	var many []RiderRequest
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
		many = append(many, RiderRequest{
			RiderID: RiderID(fmt.Sprintf("Rider %d", i)),
			PickUp:  oneOrigin,
			DropOff: p,
		})
	}
	return many
}
