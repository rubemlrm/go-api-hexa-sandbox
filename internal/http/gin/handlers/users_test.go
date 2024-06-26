package handlers_test

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

	"github.com/gin-gonic/gin"
	"github.com/rubemlrm/go-api-bootstrap/internal/http/gin/handlers"
	"github.com/rubemlrm/go-api-bootstrap/user"
	"github.com/rubemlrm/go-api-bootstrap/user/factories"
	user_mocks "github.com/rubemlrm/go-api-bootstrap/user/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

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
			mockUserService := user_mocks.NewMockUseCase(t)
			logger := slog.New(slog.NewTextHandler(io.Discard, nil))
			s := handlers.NewServer(mockUserService, logger)
			mockUserService.On("Get", mock.Anything, user.ID(tt.userID)).Return(tt.mockUser, tt.mockError)
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
			mockUserService.AssertExpectations(t)
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
			mockUserService := user_mocks.NewMockUseCase(t)
			logger := slog.New(slog.NewTextHandler(io.Discard, nil))
			s := handlers.NewServer(mockUserService, logger)
			mockUserService.On("All", mock.Anything).Return(tt.mockResultUsers, tt.mockError)
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
			mockUserService.AssertExpectations(t)
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
			expectedStatus:   http.StatusInternalServerError,
			expectedRequest:  uf.CreateUserCreate(),
			expectedResponse: `{"error":"internal error"}`,
			mockError:        errors.New("internal error"),
			mockUserID:       0,
		},
		{
			name:             "failed to bind json",
			expectedStatus:   http.StatusBadRequest,
			expectedRequest:  uf.CreateInvalidUserCreate(),
			expectedResponse: `{"error":"failed to bind"}`,
			mockError:        nil,
			mockUserID:       0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUserService := user_mocks.NewMockUseCase(t)
			logger := slog.New(slog.NewTextHandler(io.Discard, nil))
			s := handlers.NewServer(mockUserService, logger)
			mockUserService.On("Create", mock.Anything, mock.Anything).Return(user.ID(tt.mockUserID), tt.mockError).Maybe()
			router := gin.Default()
			router.POST("/api/v1/users/", s.AddUser)

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
			mockUserService.AssertExpectations(t)
		})
	}
}
