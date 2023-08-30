package gin_handler

import (
	"github.com/gin-gonic/gin"
	"github.com/rubemlrm/go-api-bootstrap/internal/http/gin/handlers"
	"github.com/rubemlrm/go-api-bootstrap/internal/http/gin/openapi"
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

func (s *Engine) SetHandlers() {
	opt := openapi.GinServerOptions{
		BaseURL: "/api/v1",
	}
	openapi.RegisterHandlersWithOptions(s.Engine, handlers.NewServer(), opt)
}

func (s *Engine) StartHTTP() http.Handler {
	return s.Engine.Handler()
}
