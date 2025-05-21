package service

import (
	"context"
	"database/sql"
	"log/slog"

	"github.com/rubemlrm/go-api-bootstrap/internal/user/adapters"
	"github.com/rubemlrm/go-api-bootstrap/internal/user/app"
	"github.com/rubemlrm/go-api-bootstrap/internal/user/app/command"
	"github.com/rubemlrm/go-api-bootstrap/internal/user/app/query"
)

func NewApplication(ctx context.Context, l *slog.Logger, db *sql.DB) app.Application {
	return newApplication(ctx, l, db)
}

func newApplication(_ context.Context, l *slog.Logger, db *sql.DB) app.Application {
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
