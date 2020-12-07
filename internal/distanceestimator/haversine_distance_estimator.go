package distanceestimator

import (
	"context"
	"math"
	"time"

	"github.com/edusalguero/roteiro.git/internal/cost"
	"github.com/edusalguero/roteiro.git/internal/point"
)

type HaversineDistanceEstimator struct {
	Velocity float64 // km per h
}

func NewHaversineDistanceEstimator(velocity float64) Service {
	return &HaversineDistanceEstimator{Velocity: velocity}
}

// degreesToRadians converts from degrees to radians.
func degreesToRadians(d float64) float64 {
	return d * math.Pi / 180
}

// Distance calculates the shortest path between two coordinates on the surface
// of the Earth. This function returns the distance in meters.
func (e *HaversineDistanceEstimator) GetCost(_ context.Context, from, to point.Point) (*cost.Cost, error) {
	const earthRadiusKm = 6371 // radius of the earth in kilometers.
	lat1 := degreesToRadians(from.Lat())
	lon1 := degreesToRadians(from.Lon())
	lat2 := degreesToRadians(to.Lat())
	lon2 := degreesToRadians(to.Lon())

	diffLat := lat2 - lat1
	diffLon := lon2 - lon1

	a := math.Pow(math.Sin(diffLat/2), 2) + math.Cos(lat1)*math.Cos(lat2)*
		math.Pow(math.Sin(diffLon/2), 2)

	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	km := c * earthRadiusKm
	meters := math.Round(km * 1000)
	duration := time.Duration(km / e.Velocity * float64(time.Hour))
	return &cost.Cost{
		Distance: meters,
		Duration: duration,
	}, nil
}
