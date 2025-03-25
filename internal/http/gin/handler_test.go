package gin_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/rubemlrm/go-api-bootstrap/config"
	"github.com/rubemlrm/go-api-bootstrap/internal/http/gin"
	"github.com/rubemlrm/go-api-bootstrap/pkg/logger"
	"github.com/rubemlrm/go-api-bootstrap/user"
	"github.com/stretchr/testify/assert"
)

func TestSetHandlers(t *testing.T) {
	engine := gin.NewEngine()
	userService := &user.Service{}
	logger := logger.NewLogger(
		config.Logger{
			Level:   "Debug",
			Handler: "textHandler",
		})

	engine.SetHandlers(userService, logger)

	// Test /swagger-ui route
	req, _ := http.NewRequest("GET", "/swagger-ui", nil)
	w := httptest.NewRecorder()
	engine.Engine.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// Test /api/v1 route (assuming there is a handler registered)
	req, _ = http.NewRequest("GET", "/api/v1", nil)
	w = httptest.NewRecorder()
	engine.Engine.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNotFound, w.Code)
}
