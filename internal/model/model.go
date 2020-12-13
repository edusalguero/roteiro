package model

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
type Load int
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
	Asset     Asset
	Requests  []Request
	Waypoints []Waypoint
	Metrics   RouteMetrics
}

type Waypoint struct {
	Location   point.Point
	Load       Load
	Activities []Activity
}

type Ref string
type ActivityType string

type Activity struct {
	ActivityType ActivityType
	Ref          Ref
}

func NewActivity(activityType ActivityType, ref Ref) Activity {
	return Activity{ActivityType: activityType, Ref: ref}
}

const (
	ActivityTypePickUp  ActivityType = "PickUp"
	ActivityTypeDropOff ActivityType = "DropOff"
	ActivityTypeStart   ActivityType = "Start"
)

type RouteMetrics struct {
	Duration time.Duration
	Distance float64
}
type SolutionMetrics struct {
	NumAssets     int
	NumRequests   int
	NumUnassigned int
	Duration      time.Duration
	Distance      float64
	SolvedTime    time.Duration
}

type Requests []*Request

type Request struct {
	RequestID          Ref
	Load               Load
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
	Ref         Ref
	Point       point.Point
	ServiceTime time.Duration
	Load        Load
	Activity    ActivityType
}

func (s Stop) IsAssetDeparture() bool {
	return s.Activity == ActivityTypeStart
}

func (s Stop) GetServiceTime() time.Duration {
	return s.ServiceTime
}

func (s Stop) String() string {
	t := "Depot"
	if !s.IsAssetDeparture() {
		t = "Request"
	}
	return fmt.Sprintf(`[%s '%s' (%s)]`, t, s.Activity, s.Point)
}
