package postgres_test

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/rubemlrm/go-api-bootstrap/internal/common/logger"
	"github.com/rubemlrm/go-api-bootstrap/internal/common/postgres"
	"github.com/rubemlrm/go-api-bootstrap/internal/common/testcontainers"
	"github.com/stretchr/testify/assert"
)

func TestNewConnectionWithAuth(t *testing.T) {
	type accessData struct {
		username string
		password string
		schema   string
		port     string
		sslmode  string
		host     string
	}
	var tests = []struct {
		name         string
		requiresAuth bool
		expectsError bool
		accessData   accessData
	}{
		{
			name:         "Start new connection with auth and success",
			requiresAuth: true,
			expectsError: false,
			accessData: accessData{
				username: "user",
				password: "password",
				schema:   "postgres",
				port:     "5432",
				sslmode:  "disable",
				host:     "localhost",
			},
		},
		{
			name:         "Start new connection with auth and fail",
			requiresAuth: true,
			expectsError: true,
			accessData: accessData{
				username: "user",
				password: "password22",
				schema:   "postgres",
				port:     "5432",
				sslmode:  "disable",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var db *sql.DB

			var container *testcontainers.PostgresTestContainer
			l := logger.NewLogger(logger.WithLogFormat("json"), logger.WithLogLevel("Debug"))
			container, err := testcontainers.StartPostgresContainer(context.Background())

			assert.NoError(t, err)
			db, err = postgres.NewConnection(
				l.Logger,
				postgres.WithUsername(tt.accessData.username),
				postgres.WithPassword(tt.accessData.password),
				postgres.WithSchema(tt.accessData.schema),
				postgres.WithHost(container.Host),
				postgres.WithPort(container.Port),
				postgres.WithSSLMode(tt.accessData.sslmode),
				postgres.WithMaxOpenConns(20),
				postgres.WithMaxIddleTime(10),
				postgres.WithMaxIddleCons(5),
			)

			if tt.expectsError {
				assert.Nil(t, db)
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
				assert.NotNil(t, db)
			}

			ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer ctxCancel()
			container.Terminate(ctx)
		})
	}
}
