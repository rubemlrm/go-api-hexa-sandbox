package service

import (
	"context"
	"github.com/rubemlrm/go-api-bootstrap/internal/common/config"
	"github.com/rubemlrm/go-api-bootstrap/internal/common/postgres"
	"github.com/rubemlrm/go-api-bootstrap/internal/user/app"
	"github.com/rubemlrm/go-api-bootstrap/internal/user/app/command"
	"github.com/rubemlrm/go-api-bootstrap/internal/user/app/query"
	"log/slog"
)

func NewApplication(ctx context.Context, cfg *config.Config, l *slog.Logger) app.Application {
	return newApplication(ctx, cfg, l)
}

func newApplication(_ context.Context, cfg *config.Config, l *slog.Logger) app.Application {
	db := postgres.StartConnection(cfg, l)

	return app.Application{
		Commands: app.Commands{
			CreateUser: command.NewCreateUserHandler(db, l),
		},
		Queries: app.Queries{
			GetUser:  query.NewGetUserHandler(db, l),
			GetUsers: query.NewListUsersHandler(db, l),
		},
	}
}
