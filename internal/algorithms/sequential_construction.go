// Algorithm based on:
//
// Manar I. Hosny, Christine L. Mumford,
// Constructing initial solutions for the multiple vehicle pickup and delivery problem with time windows,
// Journal of King Saud University - Computer and Information Sciences, Volume 24, Issue 1, 2012, Pages 59-69, ISSN 1319-1578,
// https://doi.org/10.1016/j.jksuci.2011.10.006.
//
package algorithms

import (
	"context"
	"math"
	"sort"
	"time"

	"github.com/edusalguero/roteiro.git/internal/cost"
	"github.com/edusalguero/roteiro.git/internal/logger"
	"github.com/edusalguero/roteiro.git/internal/model"
	"github.com/edusalguero/roteiro.git/internal/point"
	"github.com/edusalguero/roteiro.git/internal/routeestimator"
)

type SequentialConstruction struct {
	logger         logger.Logger
	routeEstimator routeestimator.Estimator
	costEstimator  cost.Service
}

func NewSequentialConstruction(l logger.Logger, e routeestimator.Estimator, de cost.Service) *SequentialConstruction {
	return &SequentialConstruction{logger: l, routeEstimator: e, costEstimator: de}
}

// Based on algorithm 2: The Sequential Construction
// https://www.sciencedirect.com/science/article/pii/S131915781100036X#n0040
func (a *SequentialConstruction) Solve(ctx context.Context, p model.Problem) (*model.Solution, error) {
	algoStart := time.Now()
	usedAssets := 0
	insertedRequests := 0

	assets := p.Fleet
	availableAssets := len(assets)
	totalRequests := len(p.Requests)

	var unassignedRequests model.Requests
	for i := range p.Requests {
		unassignedRequests = append(unassignedRequests, &p.Requests[i])
	}

	var solutionRoutes []model.SolutionRoute
	var totalDistance float64 = 0
	var totalDuration time.Duration
	assets = assetsWithMoreCapacityFirst(assets)

	for _, asset := range assets {
		assetLocation := asset.Location
		unassignedRequests := a.sortRequestFromDepotToDropOffFarthestFirst(ctx, assetLocation, unassignedRequests)
		a.logger.Debugf("##  Creating a new route....")
		var r model.Route
		var routeReqs []model.Request
		var err error

		r = append(r, &model.Stop{
			Ref:      model.Ref(asset.AssetID),
			Point:    assetLocation,
			Activity: model.ActivityTypeStart,
		})

		for i := range unassignedRequests {
			ur := &unassignedRequests[i]
			if *ur == nil {
				continue
			}
			req := *ur
			a.logger.Debugf("###  Adding request a new route: %s", req.RequestID)
			r = a.addRequestStops(ctx, r, assetLocation, req, p.GetMaxJourneyTimeFactor())
			r, err = a.hillClimbingRoutingAlgorithmV3(ctx, r, asset)
			if err != nil {
				return nil, err
			}

			if a.isFeasibleRoute(ctx, r, asset) {
				// Remove from unassignedRequests
				unassignedRequests[i] = nil
				insertedRequests++
				routeReqs = append(routeReqs, model.Request{
					RequestID: req.RequestID,
					PickUp:    req.PickUp,
					DropOff:   req.DropOff,
					Load:      req.Load,
				})
			} else {
				// Remove req from r
				r = removeFromRoute(r, *req)
			}
		}

		re, _ := a.routeEstimator.GetRouteEstimation(ctx, r.GetPoints())
		totalDuration += re.TotalDuration
		totalDistance += re.TotalDistance
		solutionRoutes = append(solutionRoutes, model.SolutionRoute{
			Asset:     asset,
			Requests:  routeReqs,
			Waypoints: buildRouteWaypoints(r, asset),
			Metrics:   model.RouteMetrics{Duration: re.TotalDuration, Distance: re.TotalDistance},
		})

		usedAssets++
		if insertedRequests == totalRequests {
			break
		}

		if usedAssets == availableAssets {
			break
		}
	}

	unassigned := getNotAssignedRequest(unassignedRequests)
	algoDuration := time.Since(algoStart)
	s := model.Solution{
		Metrics: model.SolutionMetrics{
			NumAssets:     usedAssets,
			NumRequests:   insertedRequests,
			NumUnassigned: len(unassigned),
			Duration:      totalDuration,
			Distance:      totalDistance,
			SolvedTime:    algoDuration,
		},
		Routes:     solutionRoutes,
		Unassigned: unassigned,
	}
	a.logger.Debugf("Solution: %v", s)

	return &s, nil
}

func buildRouteWaypoints(r model.Route, asset model.Asset) []model.Waypoint {
	var waypoints []model.Waypoint
	var load model.Load = 0
	l := len(r)
	for i := 0; i < l; {
		stop := r[i]
		p := stop.Point
		var activities []model.Activity
		if stop.Activity == model.ActivityTypeStart {
			activities = append(activities, model.NewActivity(stop.Activity, model.Ref(asset.AssetID)))
		}

		j := i
		for {
			if j == l {
				break
			}
			s := r[j]
			if s.Point != p {
				break
			}
			if s.Activity != model.ActivityTypeStart {
				load += s.Load
				activities = append(activities, model.NewActivity(s.Activity, s.Ref))
			}
			j++
		}
		i = j

		waypoints = append(waypoints, model.Waypoint{
			Location:   p,
			Load:       load,
			Activities: activities,
		})
	}

	return waypoints
}

func getNotAssignedRequest(requests model.Requests) []model.Request {
	var unassigned = make([]model.Request, 0)
	for i := range requests {
		ur := requests[i]
		if ur == nil {
			continue
		}

		// Do no copy calculated service times
		unassigned = append(unassigned, model.Request{
			RequestID: ur.RequestID,
			PickUp:    ur.PickUp,
			DropOff:   ur.DropOff,
			Load:      ur.Load,
		})
	}

	return unassigned
}

func assetsWithMoreCapacityFirst(assets []model.Asset) []model.Asset {
	sort.SliceStable(assets, func(i, j int) bool {
		return assets[i].Capacity > assets[j].Capacity
	})

	return assets
}

func (a *SequentialConstruction) addRequestStops(ctx context.Context, r model.Route, assetLocation point.Point, req *model.Request, timeFactor float64) model.Route {
	a.updateRequestServiceTime(ctx, assetLocation, req, timeFactor)
	r = append(r, &model.Stop{
		Ref:         req.RequestID,
		Point:       req.PickUp,
		ServiceTime: req.PickUpServiceTime,
		Load:        req.Load,
		Activity:    model.ActivityTypePickUp,
	})

	r = append(r, &model.Stop{
		Ref:         req.RequestID,
		Point:       req.DropOff,
		ServiceTime: req.DropOffServiceTime,
		Load:        -req.Load,
		Activity:    model.ActivityTypeDropOff,
	})
	return r
}

// Based on algorithm 1: The HC routing algorithm.
// https://www.sciencedirect.com/science/article/pii/S131915781100036X#n0035
func (a *SequentialConstruction) hillClimbingRoutingAlgorithmV3(ctx context.Context, r model.Route, asset model.Asset) (model.Route, error) {
	l := len(r)
	for i := range r {
		i := l - 1 - i
		current := r[i]
		if current.IsAssetDeparture() {
			break
		}
		for j := l - 1; j > 0; j-- {
			neighbor := r[j]
			if current.GetServiceTime() > neighbor.GetServiceTime() {
				rPrima := r
				costR, err := a.cost(ctx, r, asset)
				if err != nil {
					return nil, err
				}
				rPrima.Swap(i, j)
				costRPrima, err := a.cost(ctx, rPrima, asset)
				if err != nil {
					return nil, err
				}
				if costRPrima-costR < 0 {
					r = rPrima
					a.logger.Debugf("Updated r [Cost(%d) vs NewCost(%d)]", costR, costRPrima)
				}
			}
		}
	}
	return r, nil
}

// The cost function of a route
// It is described by the following equation:
//
// F(r) =w1×D(r) +w2×T W V(r) +w3×CV(r)
//
// where D(r)is the total route duration, including the waiting time and the service time at each location.
// TWV(r)is the total number of time window violations in the route.
// CV(r)is the total number of capacity violations.
// The constants w1,w2, and  w3 are  weights in the range [0,1], and w1+w2+w3= 1.0.
//
// The largest penalty should be imposed on the time window violations, in order to direct the search towards more feasible routes.
// We used the following weights for the route cost function:w1= 0.201,w2= 0.7 and w3= 0.0992.
func (a *SequentialConstruction) cost(ctx context.Context, r model.Route, asset model.Asset) (float64, error) {
	const (
		W1 float64 = 0.201
		W2 float64 = 0.7
		W3 float64 = 0.0992
	)

	points := r.GetPoints()
	estimation, err := a.routeEstimator.GetRouteEstimation(ctx, points)
	if err != nil {
		return math.Inf(0), err
	}

	// time window constraint violations
	twv, err := a.countTimeWindowViolations(ctx, r)
	if err != nil {
		a.logger.Debugf("Error counting time window violations: %s", err)
		return math.Inf(0), err
	}
	cv := countCapacityViolations(asset.Capacity, r)

	vs := []float64{estimation.TotalDuration.Minutes(), float64(twv), float64(cv)}
	d := normalize(estimation.TotalDuration.Minutes(), vs)
	t := normalize(float64(twv), vs)
	c := normalize(float64(cv), vs)
	return W1*d + W2*t + W3*c, nil
}

func normalize(v float64, vs []float64) float64 {
	return v - smallest(vs)/biggest(vs) - smallest(vs)
}

func smallest(values []float64) float64 {
	var n, s float64
	for _, v := range values {
		if v < n {
			n = v
			s = n
		}
	}
	return s
}

func biggest(values []float64) float64 {
	var n, b float64
	for _, v := range values {
		if v > n {
			n = v
			b = n
		}
	}
	return b
}

func (a *SequentialConstruction) sortRequestFromDepotToDropOffFarthestFirst(
	ctx context.Context,
	depot point.Point,
	requests model.Requests,
) model.Requests {
	sort.SliceStable(requests, func(i, j int) bool {
		if requests[i] == nil || requests[j] == nil {
			return false
		}
		distanceToI, _ := a.costEstimator.GetCost(ctx, depot, requests[i].DropOff)
		distanceToJ, _ := a.costEstimator.GetCost(ctx, depot, requests[j].DropOff)

		return distanceToI.Distance > distanceToJ.Distance
	})
	return requests
}

func removeFromRoute(r model.Route, req model.Request) model.Route {
	var route model.Route
	route = append(route, r[0])
	for _, stop := range r[1:] {
		skip := stop.Ref == req.RequestID
		if !skip {
			route = append(route, stop)
		}
	}
	return route
}

func (a *SequentialConstruction) isFeasibleRoute(ctx context.Context, r model.Route, asset model.Asset) bool {
	// time window constraint violations
	timeWindowViolations, err := a.countTimeWindowViolations(ctx, r)
	if err != nil {
		a.logger.Debugf("Error counting time window violations: %s", err)
		return false
	}
	//  capacity constraint violations
	violations := countCapacityViolations(asset.Capacity, r)

	feasible := violations == 0 && timeWindowViolations == 0
	a.logger.Debugf("Is Feasible %b [CV: %d / TWV %d ]", feasible, violations, timeWindowViolations)
	return feasible
}

func (a *SequentialConstruction) countTimeWindowViolations(ctx context.Context, r model.Route) (int, error) {
	points := r.GetPoints()
	twv := 0
	for i, stop := range r {
		if i == 0 {
			// no time from depot to depot
			continue
		}

		e, err := a.routeEstimator.GetRouteEstimation(ctx, points[0:i+1])
		if err != nil {
			return 0, err
		}
		ahead := stop.GetServiceTime() - e.TotalDuration
		if ahead < 0 {
			twv++
		}
	}
	return twv, nil
}

func (a *SequentialConstruction) updateRequestServiceTime(
	ctx context.Context,
	assetLocation point.Point,
	request *model.Request,
	timeFactor float64,
) {
	toPickUp, _ := a.costEstimator.GetCost(ctx, assetLocation, request.PickUp)
	request.PickUpServiceTime = increaseDurationInAFactor(toPickUp.Duration, timeFactor)

	directRoute, _ := a.routeEstimator.GetRouteEstimation(ctx, []point.Point{assetLocation, request.PickUp, request.DropOff})
	request.DropOffServiceTime = increaseDurationInAFactor(directRoute.TotalDuration, timeFactor)
}

func increaseDurationInAFactor(duration time.Duration, factor float64) time.Duration {
	d := float64(duration.Nanoseconds()) * factor
	return time.Duration(d)
}

func countCapacityViolations(capacity model.Capacity, route model.Route) int {
	violations := 0
	load := 0
	for _, s := range route {
		load += int(s.Load)
		if load > int(capacity) {
			violations++
		}
	}

	return violations
}
