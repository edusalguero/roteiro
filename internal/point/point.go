package point

import "fmt"

// Point represents a point encoded as a two-element array: latitude & longitude
type Point [2]float64

func NewPoint(lat, lon float64) Point {
	return Point{lat, lon}
}

// Lat returns the latitude of the point
func (p Point) Lat() float64 {
	return p[0]
}

// Lon returns the longitude of the point
func (p Point) Lon() float64 {
	return p[1]
}

// Check if the same point
func (p Point) Equal(p2 Point) bool {
	return p.Lon() == p2.Lon() && p.Lat() == p2.Lat()
}

func (p Point) String() string {
	return fmt.Sprintf("%v,%v", p.Lat(), p.Lon())
}

func UniquePoints(s []Point) []Point {
	if len(s) == 0 {
		return s
	}
	seen := make([]Point, 0, len(s))
slice:
	for i, n := range s {
		if i == 0 {
			s = s[:0]
		}
		for _, t := range seen {
			if n.Equal(t) {
				continue slice
			}
		}
		seen = append(seen, n)
		s = append(s, n)
	}
	return s
}
