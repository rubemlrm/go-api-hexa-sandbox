package gin

import (
	"log/slog"
	"net/http"

	"github.com/google/uuid"

	"github.com/gin-gonic/gin"
	"github.com/go-openapi/runtime/middleware"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
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

func (s *Engine) SetHandlers(logger *slog.Logger, openapiHandler func()) {
	s.Engine.Use(otelgin.Middleware("sandbox"))
	s.Engine.Use(SetRequestID())
	s.Engine.Use(otelgin.Middleware("my-server"))
	opts := middleware.SwaggerUIOpts{SpecURL: "/swagger", Path: "/swagger-ui"}
	sh := middleware.SwaggerUI(opts, nil)
	s.Engine.GET("/swagger-ui", func(ctx *gin.Context) {
		sh.ServeHTTP(ctx.Writer, ctx.Request)
	})
	openapiHandler()
}

func (s *Engine) StartHTTP() http.Handler {
	return s.Engine.Handler()
}

func SetRequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		rid, err := uuid.NewV6()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		c.Set("requestID", rid.String())
		c.Next()
	}
}
