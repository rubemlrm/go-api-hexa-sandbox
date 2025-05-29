package gin

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/rubemlrm/go-api-bootstrap/internal/common/validations"
	"log/slog"
	"net/http"
)

func ValidateRequestBody[T any](log *slog.Logger, key string) gin.HandlerFunc {
	return func(c *gin.Context) {
		var payload T
		requestID, _ := c.Get("requestID")
		if err := c.ShouldBindJSON(&payload); err != nil {
			var validationErrors validator.ValidationErrors
			ok := errors.As(err, &validationErrors)
			if ok {
				log.Warn("validation", key, "error", slog.Any("error", err), slog.String("requestID", requestID.(string)), slog.Any("context", "Validation"))
				c.JSON(http.StatusUnprocessableEntity, gin.H{
					"message": "Validation failed",
					"errors":  validations.ConvertToMap(validationErrors),
				})
				c.Abort()
			}
			log.Error("validation", key, "error", slog.Any("error", err), slog.String("requestID", requestID.(string)), slog.Any("context", "Validation"))
			c.JSON(http.StatusBadRequest, gin.H{"error": err})
			c.Abort()
		}
		// Store in context for handler
		c.Set(key, payload)
		c.Next()
	}
}
