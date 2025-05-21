package factories

import (
	"errors"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/rubemlrm/go-api-bootstrap/internal/user/domain/user"
	"github.com/stretchr/testify/assert"
)

func TestGenerateUsersOnDB(t *testing.T) {
	uf := UserFactory{}
	tests := []struct {
		name        string
		users       []user.User
		expectError bool
	}{
		{
			name: "Users generated with success",
			users: []user.User{
				*uf.CreateUser(),
				*uf.CreateUser(),
			},
			expectError: false,
		},

		{
			name: "Failed to generate users",
			users: []user.User{
				*uf.CreateUser(),
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			assert.NoError(t, err)
			defer db.Close()
			if tt.expectError {
				mock.ExpectExec("INSERT INTO users").
					WithArgs(tt.users[0].ID, tt.users[0].Name, tt.users[0].Email, tt.users[0].Password).
					WillReturnError(errors.New("db error"))
			} else {
				mock.ExpectExec(regexp.QuoteMeta(
					"INSERT INTO users (id, name, email, password) VALUES ($1, $2, $3, $4),($5, $6, $7, $8)",
				)).
					WithArgs(
						tt.users[0].ID, tt.users[0].Name, tt.users[0].Email, tt.users[0].Password,
						tt.users[1].ID, tt.users[1].Name, tt.users[1].Email, tt.users[1].Password,
					).
					WillReturnResult(sqlmock.NewResult(1, 2))
			}
			err = GenerateUsersOnDB(db, tt.users)
			if tt.expectError {
				assert.Error(t, err)
				assert.NoError(t, mock.ExpectationsWereMet())
			} else {
				assert.NoError(t, err)
				assert.NoError(t, mock.ExpectationsWereMet())
			}
		})
	}
}
func TestGenerateUsers(t *testing.T) {
	tests := []struct {
		name        string
		total       int
		expectCount int
	}{
		{
			name:        "Generate zero users",
			total:       0,
			expectCount: 0,
		},
		{
			name:        "Generate one user",
			total:       1,
			expectCount: 1,
		},
		{
			name:        "Generate multiple users",
			total:       5,
			expectCount: 5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			users := GenerateUsers(tt.total)
			assert.Len(t, users, tt.expectCount)
			ids := make(map[user.ID]struct{})
			for _, u := range users {
				ids[u.ID] = struct{}{}
			}
			assert.Equal(t, tt.expectCount, len(ids))
		})
	}
}
