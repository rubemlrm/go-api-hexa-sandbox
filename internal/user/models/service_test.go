package models_test

import (
	"context"
	"errors"
	"github.com/rubemlrm/go-api-bootstrap/internal/user/adapters/factories"
	"github.com/rubemlrm/go-api-bootstrap/internal/user/models"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUserCreation(t *testing.T) {
	uf := &factories.UserFactory{}
	tests := []struct {
		name           string
		user           *models.UserCreate
		mockUserID     int
		mockError      error
		expectedError  error
		expectedUserID ID
	}{
		{
			name:           "create user with success",
			user:           uf.CreateUserCreate(),
			mockUserID:     1,
			mockError:      nil,
			expectedError:  nil,
			expectedUserID: 1,
		},
		{
			name:           "failed to create user with success",
			user:           uf.CreateUserCreate(),
			mockUserID:     1,
			mockError:      errors.New("something went wrong"),
			expectedError:  errors.New("error creating user"),
			expectedUserID: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := models.NewMockRepository(t)
			ctx := context.Background()
			repo.On("Create", ctx, tt.user).Return(ID(tt.mockUserID), tt.mockError).Once()
			service := models.user.NewService(repo, nil)
			// act
			userCreate, err := service.Create(ctx, tt.user)
			// assert
			if tt.mockError != nil {
				assert.NotNil(t, err)
				assert.Equal(t, tt.expectedError, err)
				assert.Equal(t, tt.expectedUserID, userCreate)
			} else {
				assert.Nil(t, err)
				assert.Equal(t, ID(tt.mockUserID), userCreate)
			}
			repo.AssertExpectations(t)
		})
	}
}

func TestServiceGet(t *testing.T) {
	uf := &factories.UserFactory{}
	tests := []struct {
		name          string
		user          *models.User
		mockUserID    int
		mockError     error
		expectedError error
	}{
		{
			name:       "user found",
			user:       uf.CreateUser(),
			mockUserID: 1,
			mockError:  nil,
		},
		{
			name:          "user not found",
			user:          nil,
			mockUserID:    1,
			mockError:     errors.New("not found"),
			expectedError: errors.New("not found"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// arrange
			repo := models.NewMockRepository(t)
			ctx := context.Background()
			repo.On("Get", ctx, ID(tt.mockUserID)).Return(tt.user, tt.mockError).Once()

			service := models.user.NewService(repo, nil)
			// act
			userFound, err := service.Get(ctx, ID(tt.mockUserID))

			// assert
			if tt.mockError != nil {
				assert.NotNil(t, err)
				assert.Equal(t, tt.expectedError, err)
			} else {
				assert.Nil(t, err)
				assert.Equal(t, tt.user, userFound)
			}
			repo.AssertExpectations(t)
		})
	}
}

func TestAll(t *testing.T) {
	uu := factories.GenerateUsers(1)
	tests := []struct {
		name          string
		users         *[]User
		mockError     error
		expectedError error
	}{
		{
			name:      "users found",
			users:     &uu,
			mockError: nil,
		},
		{
			name:          "failed to fetch users",
			users:         nil,
			mockError:     errors.New("failed to fetch users"),
			expectedError: errors.New("failed to fetch users"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// arrange
			repo := usermocks.NewMockRepository(t)
			ctx := context.Background()
			repo.On("All", ctx).Return(tt.users, tt.mockError).Once()

			service := models.user.NewService(repo, nil)
			// act
			userFound, err := service.All(ctx)

			// assert
			if tt.mockError != nil {
				assert.NotNil(t, err)
				assert.Equal(t, tt.expectedError, err)
			} else {
				assert.Nil(t, err)
				assert.Equal(t, tt.users, userFound)
			}
			repo.AssertExpectations(t)
		})
	}
}
