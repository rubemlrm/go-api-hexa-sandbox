package gin_test

import (
	"bytes"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/gin-gonic/gin"
	gin_handler "github.com/rubemlrm/go-api-bootstrap/internal/common/http/gin"
	"github.com/rubemlrm/go-api-bootstrap/internal/common/validations"
)

type MockTestPayload struct {
	Name      string                              `json:"name" validate:"required"`
	CheckFunc func() ([]map[string]string, error) `json:"-"`
}

func (m *MockTestPayload) Check(vf validations.ValidationFunc) ([]map[string]string, error) {
	if m.CheckFunc != nil {
		return m.CheckFunc()
	}
	return nil, nil
}

func TestValidateRequestBody(t *testing.T) {
	var tests = []struct {
		name           string
		factory        func() error
		payload        []byte
		expectedStatus int
		expectedBody   string
		mockedBody     gin.H
		mockedError    error
		mockedReturns  []map[string]string
		mockedStatus   int
		handlerCalled  bool
		checkFunc      func() ([]map[string]string, error)
	}{
		{
			name:    "Valid payload",
			payload: []byte(`{"name":"test"}`),
			checkFunc: func() ([]map[string]string, error) {
				return nil, nil
			},
			expectedStatus: 200,
			mockedStatus:   200,
			expectedBody:   "{\"message\":\"ok\"}",
			mockedBody:     gin.H{"message": "ok"},
			mockedReturns:  nil,
			mockedError:    nil,
			handlerCalled:  true,
		},
		{
			name:    "Invalid payload",
			payload: []byte(`{"name":""}`),
			checkFunc: func() ([]map[string]string, error) {
				return []map[string]string{
					{"field": "name", "error": "name must have a value!"},
				}, nil
			},
			expectedStatus: 422,
			mockedStatus:   422,
			expectedBody:   "{\"errors\":[{\"error\":\"name must have a value!\",\"field\":\"name\"}],\"message\":\"Validation failed\"}",
			mockedBody: gin.H{
				"message": "Validation failed",
				"errors": []gin.H{
					{
						"field": "name",
						"error": "name must have a value!",
					},
				},
			},

			mockedReturns: []map[string]string{
				{"field": "name", "error": "name must have a value!"},
			},
			mockedError:   nil,
			handlerCalled: false,
		},
		{
			name:    "Failed to validate payload because of unhandled exception",
			payload: []byte(`{"":"test"}`),
			checkFunc: func() ([]map[string]string, error) {
				return nil, errors.New("invalid request body")
			},
			expectedStatus: 500,
			expectedBody:   "{\"error\":\"Unhandled exception for input validation\"}",
			mockedStatus:   500,
			mockedBody:     gin.H{"error": "Unhandled exception for input validation"},
			mockedReturns:  nil,
			mockedError:    errors.New("unhandled exception"),
			handlerCalled:  false,
		},
		{
			name:    "Failed to decode request body",
			payload: []byte(`{invalid`),
			checkFunc: func() ([]map[string]string, error) {
				return nil, errors.New("invalid request body")
			},
			expectedStatus: 400,
			expectedBody:   "{\"error\":\"invalid request body\"}",
			mockedStatus:   400,
			mockedBody:     gin.H{"error": "invalid request body"},
			mockedReturns:  nil,
			mockedError:    errors.New("invalid request body"),
			handlerCalled:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gin.SetMode(gin.TestMode)
			router := gin.Default()
			router.Use(func(c *gin.Context) {
				requestID := "test-request-id"
				c.Set("requestID", requestID)
				c.Next()
			})
			mockLogger := slog.New(slog.NewTextHandler(io.Discard, nil))
			handlerCalled := false
			router.POST("/test", gin_handler.ValidateRequestBody[*MockTestPayload](func() *MockTestPayload {
				return &MockTestPayload{
					CheckFunc: tt.checkFunc,
				}
			},
				mockLogger, "testKey"), func(c *gin.Context) {
				handlerCalled = true
				c.JSON(tt.mockedStatus, tt.mockedBody)
			})

			req := httptest.NewRequest(http.MethodPost, "/test", bytes.NewBuffer(tt.payload))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)
			// Assertions
			assert.Equal(t, tt.expectedStatus, w.Code)
			assert.Equal(t, tt.expectedBody, w.Body.String())
			assert.Equal(t, tt.handlerCalled, handlerCalled)
			//mockerdPayload.AssertExpectations(t)
		})
	}
}
