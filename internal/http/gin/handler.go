package gin_handler

import (
	"github.com/gin-gonic/gin"
	"github.com/go-openapi/runtime/middleware"
	"github.com/rubemlrm/go-api-bootstrap/internal/http/gin/handlers"
	"github.com/rubemlrm/go-api-bootstrap/internal/http/gin/openapi"
	"github.com/rubemlrm/go-api-bootstrap/user"
	"log/slog"
	"net/http"
)

type Engine struct {
	Engine *gin.Engine
}

func NewEngine() *Engine {
	eng := gin.New()
	return &Engine{
		Engine: eng,
	}
}

func (s *Engine) SetHandlers(userService *user.Service, logger *slog.Logger) {
	s.Engine.StaticFile("/swagger", "./spec/openapi.yaml")

	opts := middleware.SwaggerUIOpts{SpecURL: "/swagger", Path: "/swagger-ui"}
	sh := middleware.SwaggerUI(opts, nil)
	s.Engine.GET("/swagger-ui", func(ctx *gin.Context) {
		sh.ServeHTTP(ctx.Writer, ctx.Request)
	})
	opt := openapi.GinServerOptions{
		BaseURL: "/api/v1",
	}
	openapi.RegisterHandlersWithOptions(s.Engine, handlers.NewServer(userService, logger), opt)
}

func (s *Engine) StartHTTP() http.Handler {
	return s.Engine.Handler()
}
