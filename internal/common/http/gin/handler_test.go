package gin_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/rubemlrm/go-api-bootstrap/internal/common/http/gin"

	"github.com/rubemlrm/go-api-bootstrap/internal/common/logger"
	"github.com/stretchr/testify/assert"
)

func TestSetHandlers(t *testing.T) {
	engine := gin.NewEngine()
	l := logger.NewLogger(logger.WithLogFormat("json"), logger.WithLogLevel("Debug"))

	engine.SetHandlers(l.Logger, func() {})

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
