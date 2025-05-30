package gin_test

import (
	"bytes"
	"encoding/json"
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
	"github.com/stretchr/testify/mock"
)

var _ gin.HandlerFunc = gin_handler.ValidateRequestBody[*MockTestPayload](nil, "")

type MockTestPayload struct {
	mock.Mock
	Name string `json:"name" validate:"required"`
}

func (t *MockTestPayload) Check(vf validations.ValidationFunc) ([]map[string]string, error) {
	return make([]map[string]string, 0), nil
}

func TestValidateRequestBody(t *testing.T) {
	var tests = []struct {
		name           string
		payload        MockTestPayload
		expectedStatus int
		expectedBody   string
		mockedBody     gin.H
		mockedError    error
		mockedReturns  []map[string]string
		mockedStatus   int
	}{
		{
			name:           "Valid payload",
			payload:        MockTestPayload{Name: "test"},
			expectedStatus: 200,
			mockedStatus:   200,
			expectedBody:   "{\"message\":\"ok\"}",
			mockedBody:     gin.H{"message": "ok"},
			mockedReturns:  nil,
			mockedError:    nil,
		},
		{
			name:           "Invalid payload",
			payload:        MockTestPayload{Name: ""},
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
			mockedError: nil,
		},
		{
			name:           "Internal server error",
			payload:        MockTestPayload{Name: ""},
			expectedStatus: 500,
			expectedBody:   "{\"error\":\"Unhandled exception for input validation\"}",
			mockedStatus:   500,
			mockedBody:     gin.H{"error": "Unhandled exception for input validation"},
			mockedReturns:  nil,
			mockedError:    errors.New("unhandled exception"),
		},
	}

	for i := range tests {
		tt := &tests[i]
		t.Run(tt.name, func(t *testing.T) {
			router := gin.Default()
			router.Use(func(c *gin.Context) {
				requestID := "test-request-id"
				c.Set("requestID", requestID)
				c.Next()
			})
			mockLogger := slog.New(slog.NewTextHandler(io.Discard, nil))
			mockTestPayload := new(MockTestPayload)
			mockTestPayload.On("Check", mock.Anything).Return(tt.mockedReturns, tt.mockedError)
			router.POST("/test", gin_handler.ValidateRequestBody[*MockTestPayload](mockLogger, "testKey"), func(c *gin.Context) {
				c.JSON(tt.mockedStatus, tt.mockedBody)
			})

			body, err := json.Marshal(&tt.payload)
			assert.NoError(t, err)
			req := httptest.NewRequest(http.MethodPost, "/test", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			// Assertions
			assert.Equal(t, tt.expectedStatus, w.Code)
			assert.Equal(t, tt.expectedBody, w.Body.String())
		})
	}
}
