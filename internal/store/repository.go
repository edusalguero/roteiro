package store

import (
	"context"

	"github.com/edusalguero/roteiro.git/internal/problem"
)

//go:generate mockgen -source=./repository.go -destination=./mock/repository.go
type Repository interface {
	AddProblem(ctx context.Context, problem *problem.Problem) error
	SetSolution(ctx context.Context, id problem.ID, solution *problem.Solution) error
	SetError(ctx context.Context, id problem.ID, err error) error
	GetSolutionByProblemID(ctx context.Context, id problem.ID) (*problem.Solution, error)
}
