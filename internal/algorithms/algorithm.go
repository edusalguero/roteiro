package algorithms

import (
	"context"

	"github.com/edusalguero/roteiro.git/internal/model"
)

type Algorithm interface {
	Solve(ctx context.Context, problem model.Problem) (*model.Solution, error)
}
