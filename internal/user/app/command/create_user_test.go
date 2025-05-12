package command_test

import (
	"context"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-faker/faker/v4"
	user_mocks "github.com/rubemlrm/go-api-bootstrap/internal/user/domain/user/mocks"

	"github.com/rubemlrm/go-api-bootstrap/internal/user/app/command"
	"github.com/rubemlrm/go-api-bootstrap/internal/user/domain/user"

	"github.com/stretchr/testify/assert"
)

func TestCreateUserHandler_Handle(t *testing.T) {
	mockRepo := user_mocks.NewMockRepository(t)
	_, mockDB, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	tests := []struct {
		name          string
		input         *user.UserCreate
		mockUserID    user.ID
		mockError     error
		expectedID    user.ID
		expectedError error
		expectUser    bool
	}{
		{
			name: "successfully create user",
			input: &user.UserCreate{
				Name:     "John Doe",
				Email:    "john.doe@example.com",
				Password: faker.Password(),
			},
			mockUserID:    1,
			mockError:     nil,
			expectedID:    1,
			expectedError: nil,
			expectUser:    true,
		},
		{
			name: "fail to create user due to repository error",
			input: &user.UserCreate{
				Name:     "Jane Doe",
				Email:    "jane.doe@example.com",
				Password: faker.Password(),
			},
			mockUserID:    0,
			mockError:     errors.New("repository error"),
			expectedID:    0,
			expectedError: errors.New("repository error"),
			expectUser:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.expectUser {
				mockDB.ExpectQuery("INSERT into users (name, email, password) values($1,$2,$3) RETURNING id").
					WithArgs(tt.input.Name, tt.input.Email, tt.input.Password).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(tt.mockUserID))
			} else {
				mockDB.ExpectQuery("INSERT into users (name, email, password) values($1,$2,$3) RETURNING id").
					WithArgs(tt.input.Name, tt.input.Email, tt.input.Password).
					WillReturnError(errors.New("repository error"))
			}

			cmd := command.NewCreateUserHandler(mockRepo)
			mockRepo.On("Create", context.Background(), tt.input).Return(tt.mockUserID, tt.mockError)

			id, err := cmd.Handle(context.Background(), tt.input)
			assert.Equal(t, tt.expectedID, id)
			assert.Equal(t, tt.expectedError, err)
		})
	}
}
