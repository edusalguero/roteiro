package httpwrapper

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Controller is an application component that can define its own HTTP handlers for Gin
type Controller interface {
	AddRoutes(g *gin.Engine)
}

type Config struct {
	// The gin mode
	Mode string `default:"release"`
	// Port is the port where server will be listening on, can't be empty.
	Port string `required:"true"`
}

type ServerWrapper interface {
	Stop(ctx context.Context)
	Start() error
}

type HTTPServerWrapper struct {
	config Config
	server *http.Server
	engine *gin.Engine
}

func NewHTTPServerWrapper(config Config) *HTTPServerWrapper {
	srv := &http.Server{
		Addr: fmt.Sprintf(":%s", config.Port),
	}

	mode := config.Mode
	if mode == "" {
		// Don't let gin panic, use "release" as default mode
		mode = gin.ReleaseMode
	}
	gin.SetMode(mode)
	engine := gin.Default()
	return &HTTPServerWrapper{config: config, server: srv, engine: engine}
}

func (s *HTTPServerWrapper) Stop(ctx context.Context) {
	_ = s.server.Shutdown(ctx)
}

func (s *HTTPServerWrapper) Start() error {
	go func() {
		err := s.engine.Run(fmt.Sprintf(":%s", s.config.Port))
		if err != nil && err != http.ErrServerClosed {
			panic(err)
		}
	}()
	return nil
}

func (s *HTTPServerWrapper) AddController(controller Controller) {
	controller.AddRoutes(s.engine)
}

func (s *HTTPServerWrapper) GetGin() *gin.Engine {
	return s.engine
}
