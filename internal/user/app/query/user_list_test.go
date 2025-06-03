package query_test

import (
	"context"
	"errors"
	"testing"

	"github.com/rubemlrm/go-api-bootstrap/internal/common/decorator/mocks"
	"github.com/stretchr/testify/mock"

	"github.com/rubemlrm/go-api-bootstrap/internal/common/logger"

	"github.com/rubemlrm/go-api-bootstrap/internal/user/app/query"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-faker/faker/v4"
	user_mocks "github.com/rubemlrm/go-api-bootstrap/internal/user/domain/user/mocks"

	"github.com/rubemlrm/go-api-bootstrap/internal/user/domain/user"

	"github.com/stretchr/testify/assert"
)

func TestListUsersHandler_Handle(t *testing.T) {
	tests := []struct {
		name          string
		searchedUsers *[]user.User
		searchInput   query.UserSearchFilters
		mockError     error
		expectedError error
		expectUsers   bool
	}{
		{
			name: "successfully list users",
			searchedUsers: &[]user.User{
				{
					Name:     "John Doe",
					Email:    "john.doe@example.com",
					Password: faker.Password(),
				},
				{
					Name:     "Jane Doe2",
					Email:    "jane.doe2@example.com",
					Password: faker.Password(),
				},
			},
			searchInput:   query.UserSearchFilters{},
			mockError:     nil,
			expectedError: nil,
			expectUsers:   true,
		},
		{
			name:          "repository error",
			searchedUsers: &[]user.User{},
			searchInput:   query.UserSearchFilters{},
			mockError:     errors.New("repository error"),
			expectedError: errors.New("repository error"),
			expectUsers:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := logger.NewLogger(logger.WithLogFormat("json"), logger.WithLogLevel("Debug"))
			mockRepo := user_mocks.NewMockUserRepository(t)
			mockTracer := mocks.NewMockRecordTracer(t)
			_, mockDB, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}
			if tt.expectUsers {
				mockDB.ExpectQuery("SELECT id, name, password, is_enabled FROM users").
					WillReturnRows(
						sqlmock.NewRows([]string{"id", "name "}).AddRow((*tt.searchedUsers)[0].ID, (*tt.searchedUsers)[0].Name).AddRow((*tt.searchedUsers)[1].ID, (*tt.searchedUsers)[1].Name))
			} else {
				mockDB.ExpectQuery("SELECT id, name, password, is_enabled FROM users").
					WillReturnError(errors.New("repository error"))
			}

			cmd := query.NewListUsersHandler(mockRepo, l.Logger, mockTracer)
			mockRepo.On("All", context.Background()).Return(tt.searchedUsers, tt.mockError)
			mockTracer.On("RecordTrace", context.Background(), mock.Anything, mock.Anything).Return(nil)

			id, err := cmd.Handle(context.Background(), tt.searchInput)
			assert.Equal(t, tt.searchedUsers, id)
			assert.Equal(t, tt.expectedError, err)
		})
	}
}
