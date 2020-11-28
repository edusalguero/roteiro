package problem

import (
	"time"

	"github.com/edusalguero/roteiro.git/internal/model"
	"github.com/edusalguero/roteiro.git/internal/point"
	"github.com/google/uuid"
)

type Asset struct {
	AssetID  AssetID
	Location point.Point
	Capacity Capacity
}

type AssetID string
type Capacity int8

type Request struct {
	RequestID          RequestID
	PickUp             point.Point
	DropOff            point.Point
	PickUpServiceTime  time.Duration
	DropOffServiceTime time.Duration
}

type RequestID string

type Problem struct {
	ID          uuid.UUID
	Fleet       []Asset
	Requests    []Request
	Constraints Constraints
}

func NewProblem(id uuid.UUID, fleet []Asset, requests []Request, constraints Constraints) *Problem {
	return &Problem{ID: id, Fleet: fleet, Requests: requests, Constraints: constraints}

}

type Constraints struct {
	MaxJourneyTimeFactor float64 // Max multiplier on the direct route. Used to calculate the dropoff time offset
}

type Solution struct {
	ID string
	model.Solution
}
