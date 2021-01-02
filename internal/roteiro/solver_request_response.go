package roteiro

import (
	"time"

	"github.com/edusalguero/roteiro.git/internal/model"
	"github.com/edusalguero/roteiro.git/internal/point"
	"github.com/edusalguero/roteiro.git/internal/problem"
	"github.com/google/uuid"
)

type problemRequest struct {
	Assets      []asset     `json:"assets" binding:"required"`
	Requests    []request   `json:"requests" binding:"required"`
	Constraints constraints `json:"constraints"`
}

type asset struct {
	AssetID  string `json:"asset_id"`
	Location Point  `json:"location"`
	Capacity int    `json:"capacity"`
}

type request struct {
	RequesterID string `json:"requester_id"`
	PickUp      Point  `json:"pick_up"`
	DropOff     Point  `json:"drop_off"`
	Load        uint8  `json:"load"`
}

type constraints struct {
	MaxJourneyTimeFactor float64 `json:"max_journey_time_factor"` // Max multiplier on the direct route. Used to calculate the dropoff time offset
}

type Point struct {
	Lat float64 `json:"lat"`
	Lon float64 `json:"lon"`
}

type problemResponse struct {
	ProblemID  string    `json:"problem_id"`
	Metrics    metrics   `json:"metrics"`
	Routes     []route   `json:"routes"`
	Unassigned []request `json:"unassigned"`
}

type metrics struct {
	NumAssets     int           `json:"num_assets"`
	NumRequests   int           `json:"num_requests"`
	NumUnassigned int           `json:"num_unassigned"`
	Duration      time.Duration `json:"duration"`
	Distance      float64       `json:"distance"`
	SolvedTime    time.Duration `json:"solved_time"`
}

type route struct {
	Asset     asset        `json:"asset"`
	Requests  []request    `json:"requests"`
	Waypoints []waypoint   `json:"waypoints"`
	Metrics   routeMetrics `json:"metrics"`
}

type activity struct {
	ActivityType string `json:"activity_type"`
	Ref          string `json:"ref"`
}

type waypoint struct {
	Location   Point      `json:"location"`
	Load       int        `json:"load"`
	Activities []activity `json:"activities"`
}

type routeMetrics struct {
	Requests int           `json:"requests"`
	Duration time.Duration `json:"duration"`
	Distance float64       `json:"distance"`
}

func newSolutionResponseFromSol(solution *problem.Solution) problemResponse {
	routes := make([]route, len(solution.Routes))
	for i, r := range solution.Routes {
		reqs := make([]request, len(r.Requests))
		for iReq, req := range r.Requests {
			reqs[iReq] = newResponseRequestFromModelRequest(req)
		}

		waypoints := make([]waypoint, len(r.Waypoints))
		for iW, w := range r.Waypoints {
			activities := make([]activity, len(w.Activities))
			for iA, a := range w.Activities {
				activities[iA] = activity{
					ActivityType: string(a.ActivityType),
					Ref:          string(a.Ref),
				}
			}

			waypoints[iW] = waypoint{
				Location: Point{
					Lat: w.Location.Lat(),
					Lon: w.Location.Lon(),
				},
				Load:       int(w.Load),
				Activities: activities,
			}
		}
		ro := route{
			Asset: newResponseAssetFromSolutionRoute(r.Asset),
			Metrics: routeMetrics{
				Requests: len(reqs),
				Duration: r.Metrics.Duration,
				Distance: r.Metrics.Distance,
			},
			Requests:  reqs,
			Waypoints: waypoints,
		}

		routes[i] = ro
	}

	unassignedReqs := make([]request, len(solution.Unassigned))
	for i, req := range solution.Unassigned {
		unassignedReqs[i] = newResponseRequestFromModelRequest(req)
	}

	return problemResponse{
		ProblemID: solution.ID.String(),
		Metrics: metrics{
			NumAssets:     solution.Metrics.NumAssets,
			NumRequests:   solution.Metrics.NumRequests,
			NumUnassigned: solution.Metrics.NumUnassigned,
			Duration:      solution.Metrics.Duration,
			Distance:      solution.Metrics.Distance,
			SolvedTime:    solution.Metrics.SolvedTime,
		},
		Routes:     routes,
		Unassigned: unassignedReqs,
	}
}

func newResponseAssetFromSolutionRoute(a model.Asset) asset {
	return asset{
		AssetID: string(a.AssetID),
		Location: Point{
			Lat: a.Location.Lat(),
			Lon: a.Location.Lon(),
		},
		Capacity: int(a.Capacity),
	}
}

func newResponseRequestFromModelRequest(req model.Request) request {
	return request{
		RequesterID: string(req.RequestID),
		PickUp: Point{
			Lat: req.PickUp.Lat(),
			Lon: req.PickUp.Lon(),
		},
		DropOff: Point{
			Lat: req.DropOff.Lat(),
			Lon: req.DropOff.Lon(),
		},
		Load: uint8(req.Load),
	}
}

func newProblemFromRequest(req problemRequest, id uuid.UUID) problem.Problem {
	var fleet []problem.Asset
	for _, a := range req.Assets {
		fleet = append(fleet, problem.Asset{
			AssetID:  problem.AssetID(a.AssetID),
			Location: point.NewPoint(a.Location.Lat, a.Location.Lon),
			Capacity: problem.Capacity(a.Capacity),
		})
	}

	var reqs []problem.Request
	for _, r := range req.Requests {
		reqs = append(reqs, problem.Request{
			RequestID: problem.RequestID(r.RequesterID),
			PickUp:    point.NewPoint(r.PickUp.Lat, r.PickUp.Lon),
			DropOff:   point.NewPoint(r.DropOff.Lat, r.DropOff.Lon),
			Load:      problem.Load(r.Load),
		})
	}
	return problem.Problem{
		ID:       problem.ID{UUID: id},
		Fleet:    fleet,
		Requests: reqs,
		Constraints: problem.Constraints{
			MaxJourneyTimeFactor: req.Constraints.MaxJourneyTimeFactor,
		},
	}
}
