package adapters_test

import (
	"context"
	"database/sql"
	"errors"
	"github.com/rubemlrm/go-api-bootstrap/internal/user/adapters"
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	_ "github.com/lib/pq"
	"github.com/rubemlrm/go-api-bootstrap/internal/common/config"
	"github.com/rubemlrm/go-api-bootstrap/internal/common/logger"
	user "github.com/rubemlrm/go-api-bootstrap/internal/user/models"
	"github.com/stretchr/testify/assert"
)

func TestUserList(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	lg := logger.NewLogger(config.Logger{
		Level: "Debug",
	})
	repo := adapters.NewConnection(db, lg)
	mock.ExpectPrepare("SELECT id, name, password, is_enabd from users").ExpectQuery().WillReturnError(errors.New("error"))
	ctx := context.Background()
	users, err := repo.All(ctx)
	if err != nil {
		assert.Error(t, err)
		assert.Nil(t, users)
	}
}

func TestUserGetUser(t *testing.T) {
	tests := []struct {
		name             string
		userID           user.ID
		expectedError    bool
		expectedMockFunc func() *sql.DB
		want             string
	}{
		{
			name:          "Fail on prepare",
			userID:        user.ID(1),
			expectedError: true,
			want:          "ttstrinest",
			expectedMockFunc: func() *sql.DB {
				db, mock, err := sqlmock.New()
				if err != nil {
					t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
				}
				mock.ExpectPrepare("SELECT id, name, password, is_enabd from users where id = $1").ExpectQuery().WillReturnError(errors.New("error"))
				return db
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := tt.expectedMockFunc()
			defer db.Close()

			lg := logger.NewLogger(config.Logger{
				Level: "Debug",
			})
			repo := adapters.NewConnection(db, lg)
			ctx := context.Background()
			users, err := repo.Get(ctx, tt.userID)
			if err != nil {
				assert.Error(t, err)
				assert.Nil(t, users)
			}
		})
	}
}
