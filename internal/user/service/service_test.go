package service_test

import (
	"context"
	"database/sql"
	"github.com/rubemlrm/go-api-bootstrap/internal/user/service"
	"github.com/stretchr/testify/assert"
	"io"
	"log/slog"
	"testing"
)

func TestNewApplication(t *testing.T) {
	ctx := context.Background()
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	db, _ := sql.Open("sqlite3", ":memory:")

	app := service.NewApplication(ctx, logger, db)

	assert.NotNil(t, app.Commands.CreateUser)
	assert.NotNil(t, app.Queries.GetUser)
	assert.NotNil(t, app.Queries.GetUsers)
}
