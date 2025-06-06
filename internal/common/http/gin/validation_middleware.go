package gin

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rubemlrm/go-api-bootstrap/internal/common/validations"
)

func ValidateRequestBody[T validations.Validater[any]](factory func() T, log *slog.Logger, key string) gin.HandlerFunc {
	return func(c *gin.Context) {
		obj := factory()
		requestID, _ := c.Get("requestID")
		decoder := json.NewDecoder(c.Request.Body)
		if err := decoder.Decode(&obj); err != nil {
			log.Error("validation", key, "error", slog.Any("error", "Invalid request body"), slog.String("requestID", requestID.(string)), slog.Any("context", "Validation"))
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
			c.Abort()
			return
		}
		failedValidations, err := obj.Check(validations.New)
		if err != nil {
			log.Error("validation", key, "error", slog.Any("error", "Unhandled exception for input validation"), slog.String("requestID", requestID.(string)), slog.Any("context", "Validation"))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Unhandled exception for input validation"})
			c.Abort()
			return
		}
		if len(failedValidations) > 0 {
			log.Warn("validation", key, "error", slog.Any("error", err), slog.String("requestID", requestID.(string)), slog.Any("context", "Validation"))
			c.JSON(http.StatusUnprocessableEntity, gin.H{
				"message": "Validation failed",
				"errors":  failedValidations,
			})
			c.Abort()
			return
		}
		c.Set(key, obj)
		c.Next()
	}
}
