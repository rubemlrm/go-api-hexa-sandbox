package app

import (
	"log/slog"

	"github.com/rubemlrm/go-api-bootstrap/internal/user/app/command"
	"github.com/rubemlrm/go-api-bootstrap/internal/user/app/query"
)

type UserModule struct {
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
