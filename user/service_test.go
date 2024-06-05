package user_test

import (
	"errors"
	"github.com/rubemlrm/go-api-bootstrap/user"
	"github.com/rubemlrm/go-api-bootstrap/user/factories"
	user_mocks "github.com/rubemlrm/go-api-bootstrap/user/mocks"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUserCreation(t *testing.T) {
	uf := &factories.UserFactory{}
	tests := []struct {
		name           string
		user           *user.UserCreate
		mockUserId     int
		mockError      error
		expectedError  error
		expectedUserId user.ID
	}{
		{
			name:           "create user with success",
			user:           uf.CreateUserCreate(),
			mockUserId:     1,
			mockError:      nil,
			expectedError:  nil,
			expectedUserId: 1,
		},
		{
			name:           "failed to create user with success",
			user:           uf.CreateUserCreate(),
			mockUserId:     1,
			mockError:      errors.New("something went wrong"),
			expectedError:  errors.New("error creating user"),
			expectedUserId: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			repo := user_mocks.NewMockRepository(t)
			repo.On("Create", tt.user).Return(user.ID(tt.mockUserId), tt.mockError).Once()
			service := user.NewService(repo, nil)
			// act
			userCreate, err := service.Create(tt.user)
			// assert
			if tt.mockError != nil {
				assert.NotNil(t, err)
				assert.Equal(t, tt.expectedError, err)
				assert.Equal(t, tt.expectedUserId, userCreate)
			} else {
				assert.Nil(t, err)
				assert.Equal(t, user.ID(tt.mockUserId), userCreate)
			}
			repo.AssertExpectations(t)
		})
	}
}

func TestServiceGet(t *testing.T) {
	uf := &factories.UserFactory{}
	tests := []struct {
		name          string
		user          *user.User
		mockUserId    int
		mockError     error
		expectedError error
	}{
		{
			name:       "user found",
			user:       uf.CreateUser(),
			mockUserId: 1,
			mockError:  nil,
		},
		{
			name:          "user not found",
			user:          nil,
			mockUserId:    1,
			mockError:     errors.New("not found"),
			expectedError: errors.New("not found"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// arrange
			repo := user_mocks.NewMockRepository(t)

			repo.On("Get", user.ID(tt.mockUserId)).Return(tt.user, tt.mockError).Once()

			service := user.NewService(repo, nil)
			// act
			userFound, err := service.Get(user.ID(tt.mockUserId))

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
		users         *[]user.User
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
			repo := user_mocks.NewMockRepository(t)

			repo.On("All").Return(tt.users, tt.mockError).Once()

			service := user.NewService(repo, nil)
			// act
			userFound, err := service.All()

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
