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

func NewApplication(ctx context.Context, l *slog.Logger, db *sql.DB) app.UserModule {
	return newApplication(ctx, l, db)
}

func newApplication(_ context.Context, l *slog.Logger, db *sql.DB) app.UserModule {
	repo := adapters.NewUserRepository(db, l)
	return app.UserModule{
		Commands: app.Commands{
			CreateUser: command.NewCreateUserHandler(repo, l),
		},
		Queries: app.Queries{
			GetUser:  query.NewGetUserHandler(repo, l),
			GetUsers: query.NewListUsersHandler(repo, l),
		},
	}
}
