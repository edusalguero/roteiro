package roteiro

import (
	"errors"
	"net/http"

	"github.com/edusalguero/roteiro.git/internal/logger"
	"github.com/edusalguero/roteiro.git/internal/problem"
	"github.com/edusalguero/roteiro.git/internal/store"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ProblemController struct {
	logger logger.Logger
	repo   store.Repository
}

func NewProblemController(log logger.Logger, repo store.Repository) *ProblemController {
	return &ProblemController{
		repo:   repo,
		logger: log,
	}
}

func (c *ProblemController) AddRoutes(g *gin.Engine) {
	v1 := g.Group("/api/v1/")
	v1.GET("problem/:problem_id", c.getProblem)
}

func (c *ProblemController) getProblem(ctx *gin.Context) {
	problemID := ctx.Param("problem_id")
	log := c.logger.WithField("problem_id", problemID)

	id := problem.ID{UUID: uuid.MustParse(problemID)}
	sol, err := c.repo.GetSolutionByProblemID(ctx, id)
	if err != nil {
		log.Errorf("Error getting problem solution: %v", err)
		if errors.Is(err, store.ErrNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Problem not found!"})
			return
		}
		if errors.Is(err, store.ErrInProcess) {
			ctx.JSON(http.StatusConflict, gin.H{"error": "Solution is being processed!"})
			return
		}
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Error retrieving solution!"})
		return
	}
	res := newSolutionResponseFromSol(sol)
	log.Infof("Problem solution for [%s]... [%v]", id, res)
	ctx.JSON(http.StatusOK, res)
}
