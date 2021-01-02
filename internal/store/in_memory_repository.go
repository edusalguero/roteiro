package store

import (
	"context"
	"errors"

	"github.com/edusalguero/roteiro.git/internal/problem"
)

var ErrInProcess = errors.New("solution in process")
var ErrNotFound = errors.New("problem not found")
var ErrAlreadyExist = errors.New("problem already exist")

type InMemoryRepository struct {
	problems map[problem.ID]*Record
}

func NewInMemoryRepository() *InMemoryRepository {
	return &InMemoryRepository{problems: make(map[problem.ID]*Record)}
}

func (r *InMemoryRepository) AddProblem(_ context.Context, p *problem.Problem) error {
	_, ok := r.problems[p.ID]
	if ok {
		return ErrAlreadyExist
	}
	r.problems[p.ID] = &Record{
		Problem:        p,
		SolutionStatus: StatusProcessing,
		Solution:       nil,
	}

	return nil
}

func (r *InMemoryRepository) SetError(_ context.Context, id problem.ID, err error) error {
	record, ok := r.problems[id]
	if !ok {
		return ErrNotFound
	}
	record.Error = err
	record.SolutionStatus = StatusError

	return nil
}

func (r *InMemoryRepository) SetSolution(_ context.Context, id problem.ID, solution *problem.Solution) error {
	record, ok := r.problems[id]
	if !ok {
		return ErrNotFound
	}
	record.Solution = solution
	record.SolutionStatus = StatusDone

	return nil
}

func (r *InMemoryRepository) GetSolutionByProblemID(_ context.Context, id problem.ID) (*problem.Solution, error) {
	record, ok := r.problems[id]
	if !ok {
		return nil, ErrNotFound
	}

	if record.SolutionStatus == StatusProcessing {
		return nil, ErrInProcess
	}

	if record.SolutionStatus == StatusDone {
		return record.Solution, nil
	}

	return nil, record.Error
}
