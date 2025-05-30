package gin_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	gin_handler "github.com/rubemlrm/go-api-bootstrap/internal/common/http/gin"

	"github.com/rubemlrm/go-api-bootstrap/internal/common/logger"
	"github.com/stretchr/testify/assert"
)

func TestSetHandlers(t *testing.T) {
	engine := gin_handler.NewEngine()
	l := logger.NewLogger(logger.WithLogFormat("json"), logger.WithLogLevel("Debug"))

	engine.SetHandlers(l.Logger, func() {}, "test")

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

func TestEngine_StartHTTP(t *testing.T) {
	en := gin_handler.NewEngine()
	h := en.StartHTTP()
	assert.NotNil(t, h)
}
