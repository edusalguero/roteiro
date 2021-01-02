package store

import (
	"time"

	"github.com/edusalguero/roteiro.git/internal/problem"
)

const (
	StatusProcessing StatusType = "processing"
	StatusDone       StatusType = "done"
	StatusError      StatusType = "error"
)

type StatusType string

type Record struct {
	Problem        *problem.Problem
	SolutionStatus StatusType
	Solution       *problem.Solution
	Error          error
	Time           time.Duration
}
