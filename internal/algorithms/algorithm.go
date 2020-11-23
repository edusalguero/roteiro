package algorithms

import (
	"context"
)

type Algorithm interface {
	Solve(ctx context.Context, problem Problem) (*Solution, error)
}
