package algorithms

import (
	"fmt"
	"time"

	"github.com/edusalguero/roteiro.git/internal/point"
)

type Problem struct {
	Fleet       []Asset
	Requests    []Request
	Constraints Constraints
}

func (p Problem) GetMaxJourneyTimeFactor() float64 {
	return p.Constraints.MaxJourneyTimeFactor
}

type Asset struct {
	AssetID  AssetID
	Location point.Point
	Capacity Capacity
}

type AssetID string
type Capacity int

type Constraints struct {
	MaxJourneyTimeFactor float64 // Max multiplier on the direct route. Used to calculate the dropoff time offset
}

type Solution struct {
	Metrics    SolutionMetrics
	Routes     []SolutionRoute
	Unassigned []Request
}

type SolutionRoute struct {
	Route   []point.Point
	Metrics RouteMetrics
}
type RouteMetrics struct {
	Duration time.Duration
	Distance float64
}
type SolutionMetrics struct {
	NumAssets  uint16
	NumRiders  uint16
	Duration   time.Duration
	Distance   float64
	SolvedTime time.Duration
}

type Requests []*Request

type Request struct {
	RiderID            RiderID
	PickUp             point.Point
	DropOff            point.Point
	PickUpServiceTime  time.Duration
	DropOffServiceTime time.Duration
}

type Route []*Stop

func (r Route) Swap(i, j int) {
	r[i], r[j] = r[j], r[i]
}

func (r Route) GetPoints() []point.Point {
	var points []point.Point

	for _, stop := range r {
		points = append(points, stop.Point)
	}

	return points
}

type Stop struct {
	Name        string
	RiderID     *RiderID
	Point       point.Point
	ServiceTime time.Duration
}

type RiderID string

func (s Stop) IsDepot() bool {
	return s.RiderID == nil
}

func (s Stop) GetServiceTime() time.Duration {
	return s.ServiceTime
}

func (s Stop) String() string {
	t := "Depot"
	if !s.IsDepot() {
		t = "Request"
	}
	return fmt.Sprintf(`[%s '%s' (%s)]`, t, s.Name, s.Point)
}
