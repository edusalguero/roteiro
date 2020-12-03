package roteiro

import (
	"context"
	"net/http"

	"github.com/edusalguero/roteiro.git/internal/logger"
	"github.com/edusalguero/roteiro.git/internal/problem"
	"github.com/edusalguero/roteiro.git/internal/solver"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type SolverController struct {
	solver          solver.Service
	logger          logger.Logger
	problems        map[string]*problem.Solution
	idGeneratorFunc IDGeneratorFunc
}

type IDGeneratorFunc func() uuid.UUID

func IDGenerator() uuid.UUID {
	return uuid.New()
}

func NewSolverController(solverService solver.Service, log logger.Logger, generatorFunc IDGeneratorFunc) *SolverController {
	return &SolverController{
		solver:          solverService,
		logger:          log,
		problems:        make(map[string]*problem.Solution),
		idGeneratorFunc: generatorFunc,
	}
}

func (c *SolverController) AddRoutes(g *gin.Engine) {
	v1 := g.Group("/api/v1/")
	v1.POST("problem", c.solveProblem)
}

func (c *SolverController) solveProblem(ctx *gin.Context) {
	var problemRequest problemRequest

	if err := ctx.ShouldBindJSON(&problemRequest); err != nil {
		c.logger.Errorf("Error processing request body: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Bad request!"})
		return
	}
	id := c.idGeneratorFunc()
	c.logger.Infof("Solving problem [%s]... [%v]", id, problemRequest)
	sol, err := c.solver.SolveProblem(
		context.Background(),
		newProblemFromRequest(problemRequest, id),
	)
	if err != nil {
		c.logger.Errorf("Error solving problem: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error processing solver!"})
		return
	}
	res := newSolutionResponseFromSol(sol)
	c.logger.Infof("Problem solved[%s]... [%v]", id, res)

	c.problems[sol.ID.String()] = sol
	ctx.JSON(http.StatusOK, res)
}
