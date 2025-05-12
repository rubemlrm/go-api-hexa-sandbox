package service

import (
	"context"
	"database/sql"
	"log/slog"

	"github.com/rubemlrm/go-api-bootstrap/internal/common/config"
	"github.com/rubemlrm/go-api-bootstrap/internal/user/adapters"
	"github.com/rubemlrm/go-api-bootstrap/internal/user/app"
	"github.com/rubemlrm/go-api-bootstrap/internal/user/app/command"
	"github.com/rubemlrm/go-api-bootstrap/internal/user/app/query"
)

func NewApplication(ctx context.Context, cfg *config.Config, l *slog.Logger, db *sql.DB) app.Application {
	return newApplication(ctx, cfg, l, db)
}

func newApplication(_ context.Context, cfg *config.Config, l *slog.Logger, db *sql.DB) app.Application {
	return app.Application{
		Commands: app.Commands{
			CreateUser: command.NewCreateUserHandler(adapters.NewUserRepository(db, l)),
		},
		Queries: app.Queries{
			GetUser:  query.NewGetUserHandler(db, l),
			GetUsers: query.NewListUsersHandler(db, l),
		},
	}
}
