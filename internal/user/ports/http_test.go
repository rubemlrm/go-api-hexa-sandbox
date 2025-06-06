package ports_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	gin_handler "github.com/rubemlrm/go-api-bootstrap/internal/common/http/gin"

	"github.com/rubemlrm/go-api-bootstrap/internal/user/app"
	command_mocks "github.com/rubemlrm/go-api-bootstrap/internal/user/app/command/mocks"
	"github.com/rubemlrm/go-api-bootstrap/internal/user/app/query"
	query_mocks "github.com/rubemlrm/go-api-bootstrap/internal/user/app/query/mocks"
	"github.com/rubemlrm/go-api-bootstrap/internal/user/domain/user"
	"github.com/rubemlrm/go-api-bootstrap/internal/user/factories"
	"github.com/rubemlrm/go-api-bootstrap/internal/user/ports"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var _ ports.ServerInterface = (*MockServerInterface)(nil)

type MockServerInterface struct {
}

func (m *MockServerInterface) ListUsers(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "ListUsers called"})
}

func (m *MockServerInterface) AddUser(c *gin.Context) {
	//TODO implement me
	panic("implement me")
}

func (m *MockServerInterface) GetUser(c *gin.Context, userId int) {
	//TODO implement me
	panic("implement me")
}

func TestGetUser(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name             string
		userID           int
		mockUser         *user.User
		expectedStatus   int
		expectedResponse string
		mockError        error
	}{
		{
			userID:         1,
			name:           "get user with success",
			expectedStatus: http.StatusOK,
			mockUser: &user.User{
				Name: "test",
			},
			expectedResponse: `{"data":{"id":0,"name":"test","email":"","is_enabled":false,"created_at":"0001-01-01T00:00:00Z","updated_at":"0001-01-01T00:00:00Z"}}`,
			mockError:        nil,
		},
		{
			userID:           12,
			name:             "User not found",
			expectedStatus:   http.StatusNotFound,
			mockUser:         nil,
			expectedResponse: "{}",
			mockError:        nil,
		},
		{
			userID:           12,
			name:             "getting error fetching user",
			expectedStatus:   http.StatusInternalServerError,
			mockUser:         nil,
			expectedResponse: `{"error":"internal error"}`,
			mockError:        errors.New("internal error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger := slog.New(slog.NewTextHandler(io.Discard, nil))
			mockHandler := query_mocks.NewMockGetUserHandler(t)
			application := app.UserModule{
				Queries: app.Queries{
					GetUser: mockHandler,
				},
			}
			s := ports.NewHTTPServer(application, logger)
			mockHandler.On("Handle", mock.Anything, query.UserSearch{
				ID: user.ID(tt.userID),
			}).Return(tt.mockUser, tt.mockError)
			router := gin.Default()
			router.GET("/api/v1/users/:id", func(c *gin.Context) {
				id := c.Param("id")
				userID, _ := strconv.Atoi(id)
				s.GetUser(c, userID)
			})

			// Create a request body

			req := httptest.NewRequest(http.MethodGet, "/api/v1/users/"+strconv.Itoa(tt.userID), nil)
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			// Assertions
			assert.Equal(t, tt.expectedStatus, w.Code)
			assert.Equal(t, tt.expectedResponse, w.Body.String())
			mockHandler.AssertExpectations(t)
		})
	}
}

func TestListUsers(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name            string
		mockResultUsers *[]user.User
		expectedStatus  int
		expectedBody    string
		mockError       error
	}{
		{
			name:           "get users with success",
			expectedStatus: http.StatusOK,
			mockResultUsers: &[]user.User{
				{
					Name: "test",
				},
			},
			expectedBody: `{"data":[{"id":0,"name":"test","email":"","is_enabled":false,"created_at":"0001-01-01T00:00:00Z","updated_at":"0001-01-01T00:00:00Z"}]}`,
			mockError:    nil,
		},
		{
			name:            "get empty user list",
			expectedStatus:  http.StatusOK,
			mockResultUsers: nil,
			expectedBody:    `{"data":[]}`,
			mockError:       nil,
		},
		{
			name:            "getting error fetching users",
			expectedStatus:  http.StatusInternalServerError,
			mockResultUsers: nil,
			expectedBody:    `{"error": "internal error"}}`,
			mockError:       errors.New("internal error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger := slog.New(slog.NewTextHandler(io.Discard, nil))
			mockHandler := query_mocks.NewMockListUsersHandler(t)
			application := app.UserModule{
				Queries: app.Queries{
					GetUsers: mockHandler,
				},
			}
			s := ports.NewHTTPServer(application, logger)
			mockHandler.On("Handle", mock.Anything, query.UserSearchFilters{}).Return(tt.mockResultUsers, tt.mockError)
			router := gin.Default()
			router.GET("/api/v1/users/", func(c *gin.Context) {
				s.ListUsers(c)
			})

			// Create a request body

			req := httptest.NewRequest(http.MethodGet, "/api/v1/users/", nil)
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			// Assertions
			assert.Equal(t, tt.expectedStatus, w.Code)
			mockHandler.AssertExpectations(t)
		})
	}
}

func TestAddUser(t *testing.T) {
	gin.SetMode(gin.TestMode)
	uf := factories.UserFactory{}
	tests := []struct {
		name             string
		expectedRequest  *user.UserCreate
		expectedStatus   int
		expectedResponse string
		mockError        error
		mockUserID       int
	}{
		{
			name:             "user created with success",
			expectedStatus:   http.StatusCreated,
			expectedRequest:  uf.CreateUserCreate(),
			expectedResponse: `{"id":1}`,
			mockError:        nil,
			mockUserID:       1,
		},
		{
			name:             "failed to create user",
			expectedStatus:   http.StatusUnprocessableEntity,
			expectedRequest:  uf.CreateInvalidUserCreate(),
			expectedResponse: `{"errors":[{"error":"Key: 'UserCreate.email' Error:Field validation for 'email' failed on the 'email' tag","field":"email"}],"message":"Validation failed"}`,
			mockError:        errors.New(`"{"errors":[{"error":"Key: 'UserCreate.email' Error:Field validation for 'email' failed on the 'email' tag","field":"email"}],"message":"Validation failed"}"`),
			mockUserID:       0,
		},
		{
			name:             "failed to create user",
			expectedStatus:   http.StatusInternalServerError,
			expectedRequest:  uf.CreateUserCreate(),
			expectedResponse: `{"errors":"Internal error"}`,
			mockError:        errors.New(`Internal error`),
			mockUserID:       0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger := slog.New(slog.NewTextHandler(io.Discard, nil))
			mockHandler := command_mocks.NewMockCreateUserHandler(t)
			application := app.UserModule{
				Commands: app.Commands{
					CreateUser: mockHandler,
				},
			}
			s := ports.NewHTTPServer(application, logger)
			mockHandler.On("Handle", mock.Anything, mock.Anything).Return(user.ID(tt.mockUserID), tt.mockError).Maybe()
			router := gin.Default()
			router.Use(func(c *gin.Context) {
				requestID := "test-request-id"
				c.Set("requestID", requestID)
				c.Next()
			})
			router.POST("/api/v1/users/", gin_handler.ValidateRequestBody[*user.UserCreate](func() *user.UserCreate {
				return &user.UserCreate{}
			}, logger, "userCreate"), s.AddUser)

			var requestBody []byte
			var err error
			if tt.expectedRequest != nil {
				requestBody, err = json.Marshal(tt.expectedRequest)
				assert.NoError(t, err)
			}

			req := httptest.NewRequest(http.MethodPost, "/api/v1/users/", bytes.NewBuffer(requestBody))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			// Assertions
			assert.Equal(t, tt.expectedStatus, w.Code)
			assert.Equal(t, tt.expectedResponse, w.Body.String())
			mockHandler.AssertExpectations(t)
		})
	}
}

func TestRegisterHandlersWithOptionsAndValidations(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	mockSI := &MockServerInterface{}
	options := ports.GinServerOptions{BaseURL: "/api/v1"}
	logger := slog.Default()

	ports.RegisterHandlersWithOptionsAndValidations(router, mockSI, options, logger)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/users", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	// Add more assertions as needed
}

func TestRegisterHandlersWithOptionsAndValidationsWithoutErrorHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	mockSI := &MockServerInterface{}
	var options ports.GinServerOptions

	logger := slog.Default()

	ports.RegisterHandlersWithOptionsAndValidations(router, mockSI, options, logger)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/users", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	// Add more assertions as needed
}
