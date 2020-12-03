package roteiro

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type StatusController struct {
}

func NewStatusController() *StatusController {
	return &StatusController{}
}

func (c *StatusController) AddRoutes(g *gin.Engine) {
	g.GET("status", c.status)
}

func (c *StatusController) status(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"status": "ok",
	})
}
