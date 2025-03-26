package app

import (
	"github.com/rubemlrm/go-api-bootstrap/internal/user/app/command"
	"github.com/rubemlrm/go-api-bootstrap/internal/user/app/query"
	"log/slog"
)

type Application struct {
	Commands Commands
	Queries  Queries
	Logger   *slog.Logger
}
type Commands struct {
	CreateUser command.CreateUserHandler
}
type Queries struct {
	GetUser  query.GetUserHandler
	GetUsers query.ListUsersHandler
}
