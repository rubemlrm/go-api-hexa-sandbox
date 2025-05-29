package query_test

import (
	"context"
	"errors"
	"testing"

	"github.com/rubemlrm/go-api-bootstrap/internal/common/logger"

	"github.com/rubemlrm/go-api-bootstrap/internal/user/app/query"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-faker/faker/v4"
	user_mocks "github.com/rubemlrm/go-api-bootstrap/internal/user/domain/user/mocks"

	"github.com/rubemlrm/go-api-bootstrap/internal/user/domain/user"

	"github.com/stretchr/testify/assert"
)

func TestGetUserHandler_Handle(t *testing.T) {
	mockRepo := user_mocks.NewMockUserRepository(t)
	_, mockDB, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	tests := []struct {
		name          string
		searchedUser  *user.User
		searchInput   query.UserSearch
		mockError     error
		expectedError error
		expectUser    bool
	}{
		{
			name: "successfully create user",
			searchedUser: &user.User{
				Name:     "John Doe",
				Email:    "john.doe@example.com",
				Password: faker.Password(),
			},
			searchInput: query.UserSearch{
				ID: 1,
			},
			mockError:     nil,
			expectedError: nil,
			expectUser:    true,
		},
		{
			name: "fail to create user due to repository error",
			searchedUser: &user.User{
				Name:     "Jane Doe",
				Email:    "jane.doe@example.com",
				Password: faker.Password(),
			},
			searchInput: query.UserSearch{
				ID: 0,
			},
			mockError:     errors.New("repository error"),
			expectedError: errors.New("repository error"),
			expectUser:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := logger.NewLogger(logger.WithLogFormat("json"), logger.WithLogLevel("Debug"))
			if tt.expectUser {
				mockDB.ExpectQuery("SELECT id, name, password, is_enabled FROM users where id = $1").
					WithArgs(tt.searchInput.ID).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(tt.searchInput.ID))
			} else {
				mockDB.ExpectQuery("SELECT id, name, password, is_enabled FROM users where id = $1").
					WithArgs(tt.searchInput.ID).
					WillReturnError(errors.New("repository error"))
			}

			cmd := query.NewGetUserHandler(mockRepo, l.Logger)
			mockRepo.On("Get", context.Background(), tt.searchInput.ID).Return(tt.searchedUser, tt.mockError)

			id, err := cmd.Handle(context.Background(), tt.searchInput)
			assert.Equal(t, tt.searchedUser, id)
			assert.Equal(t, tt.expectedError, err)
		})
	}
}
