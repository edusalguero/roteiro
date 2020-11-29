package main

import (
	"context"
	"fmt"

	"github.com/edusalguero/roteiro.git/internal/algorithms"
	"github.com/edusalguero/roteiro.git/internal/config"
	"github.com/edusalguero/roteiro.git/internal/distanceestimator"
	"github.com/edusalguero/roteiro.git/internal/logger"
	"github.com/edusalguero/roteiro.git/internal/point"
	"github.com/edusalguero/roteiro.git/internal/problem"
	"github.com/edusalguero/roteiro.git/internal/routeestimator"
	"github.com/edusalguero/roteiro.git/internal/solver"
	"github.com/edusalguero/roteiro.git/internal/staticmap"
	"github.com/fogleman/gg"
)

func main() {
	cnf, err := config.Get()
	if err != nil {
		panic("could not load config: " + err.Error())
	}

	log, err := logger.New(cnf.Log)
	if err != nil {
		panic("could not initialize logger: " + err.Error())
	}

	_, err = distanceestimator.NewGoogleMapsDistanceEstimator(cnf.DistanceEstimator.GoogleMaps.APIKey, logger.NewNopLogger())
	if err != nil {
		log.Fatalf("NewGoogleMapsDistanceEstimator() error = %v", err)
	}

	e := distanceestimator.NewHaversineDistanceEstimator(80)
	routeE := routeestimator.NewEstimator(e)
	algo := algorithms.NewSequentialConstruction(log, routeE, e)
	slvr := solver.NewService(algo)

	var oneOrigin = point.NewPoint(4.68295, -74.04965)

	prblm := problem.Problem{
		Fleet: []problem.Asset{
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
		Requests: manyRequests(oneOrigin),
		Constraints: problem.Constraints{
			MaxJourneyTimeFactor: 1.5,
		},
	}

	sol, err := slvr.SolveProblem(context.Background(), prblm)

	if err != nil {
		log.Panicf("could not solve problem: ", err.Error())
	}

	render := staticmap.NewService()
	img, err := render.Render(sol)
	if err != nil {
		log.Panicf("could not render solution: ", err.Error())
	}
	if err := gg.SavePNG("my-map.png", img); err != nil {
		log.Info("Image saved: my-map.png")
	}
	log.Info("Done")
}

func manyRequests(oneOrigin point.Point) []problem.Request {
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
