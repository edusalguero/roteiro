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
func (a *SequentialConstruction) Solve(ctx context.Context, p Problem) (*Solution, error) {
	algoStart := time.Now()
	usedAssets := 0
	insertedRequests := 0

	asset := p.Fleet[0]
	assetLocation := asset.Location
	requests := a.sortRequestFromDepotToDropOffFarthestFirst(ctx, assetLocation, p.Requests)

	var unassignedRequests Requests
	for i := range requests {
		unassignedRequests = append(unassignedRequests, &requests[i])
	}

	totalRequests := len(requests)

	var routes []Route
	for {
		a.logger.Debugf("##  Creating a new route....")
		usedAssets++
		var r Route
		var err error

		r = append(r, &Stop{
			Name:        "Asset",
			RiderID:     nil,
			Point:       assetLocation,
			ServiceTime: 0,
		})

		for i := range unassignedRequests {
			ur := &unassignedRequests[i]
			if *ur == nil {
				continue
			}
			req := *ur
			a.logger.Debugf("###  Adding request a new route: %s", req.RiderID)
			r = a.addRequestStops(ctx, r, assetLocation, req, p.GetMaxJourneyTimeFactor())

			r, err = a.hillClimbingRoutingAlgorithmV2(ctx, r, asset)
			if err != nil {
				return nil, err
			}

			if a.isFeasibleRoute(ctx, r, asset) {
				// Remove from unassignedRequests
				unassignedRequests[i] = nil
				insertedRequests++
			} else {
				// Remove req from r
				r = removeFromRoute(r, *req)
			}
		}
		routes = append(routes, r)

		if insertedRequests == totalRequests {
			break
		}
	}

	var solutionRoutes []SolutionRoute
	var distance float64 = 0
	var duration time.Duration

	for _, r := range routes {
		re, _ := a.routeEstimator.GetRouteEstimation(ctx, r.GetPoints())
		duration += re.TotalDuration
		distance += re.TotalDistance
		solutionRoutes = append(solutionRoutes, SolutionRoute{
			Route:   point.UniquePoints(r.GetPoints()),
			Metrics: RouteMetrics{Duration: re.TotalDuration, Distance: re.TotalDistance},
		})
	}

	algoDuration := time.Since(algoStart)
	s := Solution{
		Metrics: SolutionMetrics{
			NumAssets:  uint16(usedAssets),
			NumRiders:  uint16(insertedRequests),
			Duration:   duration,
			Distance:   distance,
			SolvedTime: algoDuration,
		},
		Routes: solutionRoutes,
	}
	a.logger.Debugf("Solution: %v", s)

	return &s, nil
}

func (a *SequentialConstruction) addRequestStops(ctx context.Context, r Route, assetLocation point.Point, req *RiderRequest, timeFactor float64) Route {
	a.updateRequestServiceTime(ctx, assetLocation, req, timeFactor)
	r = append(r, &Stop{
		Name:        fmt.Sprintf("%s PickUp", req.RiderID),
		RiderID:     &req.RiderID,
		Point:       req.PickUp,
		ServiceTime: req.PickUpServiceTime,
	})

	r = append(r, &Stop{
		Name:        fmt.Sprintf("%s DropOff", req.RiderID),
		RiderID:     &req.RiderID,
		Point:       req.DropOff,
		ServiceTime: req.DropOffServiceTime,
	})
	return r
}

// Based on algorithm 1: The HC routing algorithm.
// https://www.sciencedirect.com/science/article/pii/S131915781100036X#n0035
func (a *SequentialConstruction) hillClimbingRoutingAlgorithmV2(ctx context.Context, r Route, asset Asset) (Route, error) {
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

func (a *SequentialConstruction) cost(ctx context.Context, r Route, asset Asset) (float64, error) {
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
	requests []RiderRequest,
) []RiderRequest {
	sort.SliceStable(requests, func(i, j int) bool {
		distanceToI, _ := a.distanceEstimator.EstimateDistance(ctx, depot, requests[i].DropOff)
		distanceToJ, _ := a.distanceEstimator.EstimateDistance(ctx, depot, requests[j].DropOff)

		return distanceToI.Distance > distanceToJ.Distance
	})
	return requests
}

func removeFromRoute(r Route, req RiderRequest) Route {
	var route Route
	route = append(route, r[0])
	for _, stop := range r[1:] {
		skip := *stop.RiderID == req.RiderID
		if !skip {
			route = append(route, stop)
		}
	}
	return route
}

func (a *SequentialConstruction) isFeasibleRoute(ctx context.Context, r Route, asset Asset) bool {
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
	request *RiderRequest,
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

func countCapacityViolations(capacity int, route Route) int {
	reqs := countAssignedRequests(route)
	if reqs > capacity {
		return int(math.Abs(float64(reqs - capacity)))
	}
	return 0
}

func countAssignedRequests(route Route) int {
	return (len(route) - 1) / 2
}

func countTimeWindowViolations(route Route, estimation *routeestimator.Estimation) int {
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
