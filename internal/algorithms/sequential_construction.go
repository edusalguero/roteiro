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
	"fmt"
	"math"
	"sort"
	"time"

	"github.com/edusalguero/roteiro.git/internal/distanceestimator"
	"github.com/edusalguero/roteiro.git/internal/logger"
	"github.com/edusalguero/roteiro.git/internal/model"
	"github.com/edusalguero/roteiro.git/internal/point"
	"github.com/edusalguero/roteiro.git/internal/routeestimator"
)

type SequentialConstruction struct {
	logger            logger.Logger
	routeEstimator    routeestimator.Estimator
	distanceEstimator distanceestimator.Service
}

func NewSequentialConstruction(l logger.Logger, e routeestimator.Estimator, de distanceestimator.Service) *SequentialConstruction {
	return &SequentialConstruction{logger: l, routeEstimator: e, distanceEstimator: de}
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
			Name:        "Asset",
			RequestID:   nil,
			Point:       assetLocation,
			ServiceTime: 0,
		})

		for i := range unassignedRequests {
			ur := &unassignedRequests[i]
			if *ur == nil {
				continue
			}
			req := *ur
			a.logger.Debugf("###  Adding request a new route: %s", req.RequestID)
			r = a.addRequestStops(ctx, r, assetLocation, req, p.GetMaxJourneyTimeFactor())
			r, err = a.hillClimbingRoutingAlgorithmV2(ctx, r, asset)
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
			Waypoints: buildWaypoints(point.UniquePoints(r.GetPoints()), routeReqs, asset),
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
			NumAssets:   usedAssets,
			NumRequests: insertedRequests,
			Duration:    totalDuration,
			Distance:    totalDistance,
			SolvedTime:  algoDuration,
		},
		Routes:     solutionRoutes,
		Unassigned: unassigned,
	}
	a.logger.Debugf("Solution: %v", s)

	return &s, nil
}

func buildWaypoints(points []point.Point, reqs []model.Request, asset model.Asset) []model.Waypoint {
	var waypoints []model.Waypoint

	var load model.Capacity = 0
	for i, p := range points {
		var activities []model.Activity
		if i == 0 {
			activities = append(activities, model.NewActivity(model.ActivityTypeStart, model.Ref(asset.AssetID)))
		}
		for _, req := range reqs {
			var activityType model.ActivityType
			if req.PickUp == p {
				activityType = model.ActivityTypePickUp
				load++
				activities = append(activities, model.NewActivity(activityType, model.Ref(req.RequestID)))
			}

			if req.DropOff == p {
				activityType = model.ActivityTypeDropOff
				load--
				activities = append(activities, model.NewActivity(activityType, model.Ref(req.RequestID)))
			}
		}

		waypoints = append(waypoints, model.Waypoint{
			Location:   p,
			Load:       load,
			Activities: activities,
		})
	}

	return waypoints
}

func getNotAssignedRequest(requests model.Requests) []model.Request {
	unassigned := make([]model.Request, 0)
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
		Name:        fmt.Sprintf("%s PickUp", req.RequestID),
		RequestID:   &req.RequestID,
		Point:       req.PickUp,
		ServiceTime: req.PickUpServiceTime,
	})

	r = append(r, &model.Stop{
		Name:        fmt.Sprintf("%s DropOff", req.RequestID),
		RequestID:   &req.RequestID,
		Point:       req.DropOff,
		ServiceTime: req.DropOffServiceTime,
	})
	return r
}

// Based on algorithm 1: The HC routing algorithm.
// https://www.sciencedirect.com/science/article/pii/S131915781100036X#n0035
func (a *SequentialConstruction) hillClimbingRoutingAlgorithmV2(ctx context.Context, r model.Route, asset model.Asset) (model.Route, error) {
	l := len(r)
	cursor := l - 1

	for {
		current := r[cursor]
		compareCursor := l - 1
		if current.IsDepot() {
			break
		}
		for {
			neighbor := r[compareCursor]
			if neighbor.IsDepot() {
				break
			}

			if compareCursor == cursor {
				// Same point
				compareCursor--
				continue
			}

			a.logger.Debugf("Comparing ServiceTimes [%s] vs [%s]: (%v, %v)%v", current.Name, neighbor.Name, current.GetServiceTime(), neighbor.GetServiceTime(), current.GetServiceTime() < neighbor.GetServiceTime())
			if current.GetServiceTime() > neighbor.GetServiceTime() {
				rPrima := r
				costR, err := a.cost(ctx, r, asset)
				if err != nil {
					return nil, err
				}
				rPrima.Swap(cursor, compareCursor)
				costRPrima, err := a.cost(ctx, rPrima, asset)
				if err != nil {
					return nil, err
				}
				if costRPrima-costR < 0 {
					r = rPrima
					a.logger.Debugf("Updated r [Cost(%d) vs NewCost(%d)]", costR, costRPrima)
				}
			}
			compareCursor--
		}
		cursor--
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

	twv := countTimeWindowViolations(r, estimation)
	cv := countCapacityViolations(int(asset.Capacity), r)
	return W1*estimation.TotalDuration.Minutes() + W2*float64(twv) + W3*float64(cv), nil
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
		distanceToI, _ := a.distanceEstimator.EstimateDistance(ctx, depot, requests[i].DropOff)
		distanceToJ, _ := a.distanceEstimator.EstimateDistance(ctx, depot, requests[j].DropOff)

		return distanceToI.Distance > distanceToJ.Distance
	})
	return requests
}

func removeFromRoute(r model.Route, req model.Request) model.Route {
	var route model.Route
	route = append(route, r[0])
	for _, stop := range r[1:] {
		skip := *stop.RequestID == req.RequestID
		if !skip {
			route = append(route, stop)
		}
	}
	return route
}

func (a *SequentialConstruction) isFeasibleRoute(ctx context.Context, r model.Route, asset model.Asset) bool {
	// time window constraint violations
	points := r.GetPoints()
	for i := range points {
		if i == 0 {
			// no time from depot to depot
			continue
		}

		e, err := a.routeEstimator.GetRouteEstimation(ctx, points[0:i+1])
		if err != nil {
			return false
		}
		ahead := r[i].GetServiceTime() - e.TotalDuration
		if ahead < 0 {
			return false
		}
	}

	//  capacity constraint violations
	occupied := countAssignedRequests(r)
	feasible := occupied <= int(asset.Capacity)
	a.logger.Debugf("Is Feasible %b [%d/%d]", feasible, occupied, asset.Capacity)
	return feasible
}

func (a *SequentialConstruction) updateRequestServiceTime(
	ctx context.Context,
	assetLocation point.Point,
	request *model.Request,
	timeFactor float64,
) {
	toPickUp, _ := a.distanceEstimator.EstimateDistance(ctx, assetLocation, request.PickUp)
	request.PickUpServiceTime = increaseDurationInAFactor(toPickUp.Duration, timeFactor)

	directRoute, _ := a.routeEstimator.GetRouteEstimation(ctx, []point.Point{assetLocation, request.PickUp, request.DropOff})
	request.DropOffServiceTime = increaseDurationInAFactor(directRoute.TotalDuration, timeFactor)
}

func increaseDurationInAFactor(duration time.Duration, factor float64) time.Duration {
	d := float64(duration.Nanoseconds()) * factor
	return time.Duration(d)
}

func countCapacityViolations(capacity int, route model.Route) int {
	reqs := countAssignedRequests(route)
	if reqs > capacity {
		return int(math.Abs(float64(reqs - capacity)))
	}
	return 0
}

func countAssignedRequests(route model.Route) int {
	return (len(route) - 1) / 2
}

func countTimeWindowViolations(route model.Route, estimation *routeestimator.Estimation) int {
	twv := 0
	for _, stop := range route {
		duration := time.Duration(0)
		for _, leg := range estimation.Legs {
			if leg.To == stop.Point {
				duration += leg.Duration
			}
		}
		if duration > stop.GetServiceTime() {
			twv++
		}
	}

	return twv
}
