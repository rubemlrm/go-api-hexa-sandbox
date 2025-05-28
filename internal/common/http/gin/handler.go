package gin

import (
	"errors"
	"github.com/go-playground/validator/v10"
	"github.com/rubemlrm/go-api-bootstrap/internal/common/validations"
	"log/slog"
	"net/http"
	"time"

	"github.com/google/uuid"

	"github.com/gin-gonic/gin"
	"github.com/go-openapi/runtime/middleware"
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

func (s *Engine) SetHandlers(logger *slog.Logger, openapiHandler func(), appName string) {
	s.Engine.StaticFile("/swagger", "./spec/user.yaml")
	// s.Engine.Use(otelgin.Middleware(appName))
	s.Engine.Use(SetRequestID())
	s.Engine.Use(RequestLogger(logger))
	// s.Engine.Use(otelgin.Middleware(appName))
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

func RequestLogger(log *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID, _ := c.Get("requestID")
		la := RequestLog{
			Request: RequestMetadata{
				Method:   c.Request.Method,
				URL:      c.Request.URL.String(),
				Path:     c.Request.URL.Path,
				Query:    c.Request.URL.Query(),
				RemoteIP: c.ClientIP(),
				UserID:   c.GetString("userID"),
				Body:     c.GetString("requestBody"),
			},
		}

		start := time.Now()
		c.Next()
		la.Response = ResponseMetadata{
			Status:    c.Writer.Status(),
			LatencyMS: time.Since(start).Milliseconds(),
		}
		logData := []any{
			"req", la,
			"requestID", requestID.(string),
			"context", "http",
		}
		log.Log(c, convertToLogLevel(c.Writer.Status()), "Request received", logData...)
	}
}

func ValidateRequestBody[T any](log *slog.Logger, handler func(c *gin.Context, dto T)) gin.HandlerFunc {
	return func(c *gin.Context) {
		var payload T
		requestID, _ := c.Get("requestID")
		if err := c.ShouldBindJSON(&payload); err != nil {
			var validationErrors validator.ValidationErrors
			ok := errors.As(err, &validationErrors)
			if ok {
				log.Warn("validation", "creation", "error", slog.Any("error", err), slog.String("requestID", requestID.(string)), slog.Any("context", "Validation"))
				c.JSON(http.StatusUnprocessableEntity, gin.H{
					"message": "Validation failed",
					"errors":  validations.ConvertToMap(validationErrors),
				})
				return
			}
			log.Error("validation", "creation", "error", slog.Any("error", err), slog.String("requestID", requestID.(string)), slog.Any("context", "Validation"))
			c.JSON(http.StatusBadRequest, gin.H{"error": err})
			return
		}
		handler(c, payload)
	}
}

func convertToLogLevel(level int) slog.Level {
	if level > 500 {
		return slog.LevelError
	}
	if level >= 400 && level < 500 {
		return slog.LevelWarn
	}
	return slog.LevelInfo
}
