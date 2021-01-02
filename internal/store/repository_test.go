package store

import (
	"context"
	"fmt"
	"testing"

	"github.com/edusalguero/roteiro.git/internal/model"
	"github.com/edusalguero/roteiro.git/internal/problem"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestInMemoryRepository_AddProblem(t *testing.T) {
	t.Run("Add if not exist", func(t *testing.T) {
		problemID := problem.ID{UUID: uuid.New()}
		p := &problem.Problem{ID: problemID}
		record := &Record{
			Problem:        p,
			SolutionStatus: StatusProcessing,
			Solution:       nil,
			Error:          nil,
			Time:           0,
		}
		r := NewInMemoryRepository()
		if err := r.AddProblem(context.Background(), p); err != nil {
			t.Errorf("AddProblem() error = %v", err)
		}
		assert.NotNil(t, r.problems[p.ID])
		assert.Equal(t, r.problems[p.ID], record)
	})

	t.Run("Error if already exist", func(t *testing.T) {
		problemID := problem.ID{UUID: uuid.New()}
		p := &problem.Problem{ID: problemID}
		r := InMemoryRepository{problems: map[problem.ID]*Record{
			problemID: {
				Problem:        p,
				SolutionStatus: StatusProcessing,
				Solution:       nil,
				Error:          nil,
				Time:           0,
			},
		}}
		err := r.AddProblem(context.Background(), p)
		assert.Error(t, err, ErrAlreadyExist)
	})
}

func TestInMemoryRepository_GetSolutionByProblemID(t *testing.T) {
	t.Run("Err if not exist", func(t *testing.T) {
		problemID := problem.ID{UUID: uuid.New()}
		r := NewInMemoryRepository()
		_, err := r.GetSolutionByProblemID(context.Background(), problemID)
		assert.Error(t, err, ErrNotFound)
	})

	t.Run("Error if not solved", func(t *testing.T) {
		problemID := problem.ID{UUID: uuid.New()}
		p := &problem.Problem{ID: problemID}
		r := InMemoryRepository{problems: map[problem.ID]*Record{
			problemID: {
				Problem:        p,
				SolutionStatus: StatusProcessing,
				Solution:       nil,
				Error:          nil,
				Time:           0,
			},
		}}

		sol, err := r.GetSolutionByProblemID(context.Background(), problemID)
		assert.Error(t, err, ErrInProcess)
		assert.Nil(t, sol)
	})

	t.Run("Return solution if exist", func(t *testing.T) {
		problemID := problem.ID{UUID: uuid.New()}
		p := &problem.Problem{ID: problemID}
		s := &problem.Solution{
			ID:       problemID,
			Solution: model.Solution{},
		}
		r := InMemoryRepository{problems: map[problem.ID]*Record{
			problemID: {
				Problem:        p,
				SolutionStatus: StatusDone,
				Solution:       s,
				Error:          nil,
				Time:           0,
			},
		}}

		sol, err := r.GetSolutionByProblemID(context.Background(), problemID)
		assert.NoError(t, err)
		assert.NotNil(t, sol)
		assert.Equal(t, s, sol)
	})

	t.Run("Return error if status error", func(t *testing.T) {
		problemID := problem.ID{UUID: uuid.New()}
		p := &problem.Problem{ID: problemID}
		s := &problem.Solution{
			ID:       problemID,
			Solution: model.Solution{},
		}
		r := InMemoryRepository{problems: map[problem.ID]*Record{
			problemID: {
				Problem:        p,
				SolutionStatus: StatusError,
				Solution:       s,
				Error:          ErrInProcess,
				Time:           0,
			},
		}}

		sol, err := r.GetSolutionByProblemID(context.Background(), problemID)
		assert.Error(t, err)
		assert.Nil(t, sol)
		assert.EqualError(t, err, ErrInProcess.Error())
	})
}

func TestInMemoryRepository_SetError(t *testing.T) {
	t.Run("Err if not exist", func(t *testing.T) {
		problemID := problem.ID{UUID: uuid.New()}
		r := NewInMemoryRepository()
		err := r.SetError(context.Background(), problemID, fmt.Errorf("some error"))
		assert.Error(t, err, ErrNotFound)
	})

	t.Run("Set the error", func(t *testing.T) {
		problemID := problem.ID{UUID: uuid.New()}
		p := &problem.Problem{ID: problemID}
		r := InMemoryRepository{problems: map[problem.ID]*Record{
			problemID: {
				Problem:        p,
				SolutionStatus: StatusError,
				Solution:       nil,
				Error:          fmt.Errorf("some error"),
				Time:           0,
			},
		}}

		err := r.SetError(context.Background(), problemID, fmt.Errorf("some error"))
		assert.NoError(t, err)
		assert.Error(t, r.problems[problemID].Error, fmt.Errorf("some error"))
	})
}

func TestInMemoryRepository_SetSolution(t *testing.T) {
	t.Run("Err if not exist", func(t *testing.T) {
		problemID := problem.ID{UUID: uuid.New()}
		r := NewInMemoryRepository()
		err := r.SetSolution(context.Background(), problemID, &problem.Solution{})
		assert.Error(t, err, ErrNotFound)
	})

	t.Run("Ok if exist", func(t *testing.T) {
		problemID := problem.ID{UUID: uuid.New()}
		p := &problem.Problem{ID: problemID}
		s := &problem.Solution{
			ID:       problemID,
			Solution: model.Solution{},
		}
		r := InMemoryRepository{problems: map[problem.ID]*Record{
			problemID: {
				Problem:        p,
				SolutionStatus: StatusProcessing,
				Solution:       nil,
				Error:          ErrInProcess,
				Time:           0,
			},
		}}
		err := r.SetSolution(context.Background(), problemID, &problem.Solution{ID: problemID})
		assert.NoError(t, err)
		assert.Equal(t, r.problems[problemID].Solution, s)
	})
}
