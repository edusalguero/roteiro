package algorithms

import (
	"context"

	"github.com/edusalguero/roteiro.git/internal/model"
)

//go:generate mockgen -source=./algorithm.go -destination=./mock/algorithm.go
type Algorithm interface {
	Solve(ctx context.Context, problem model.Problem) (*model.Solution, error)
}
