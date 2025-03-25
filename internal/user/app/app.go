package app

import (
	"github.com/rubemlrm/go-api-bootstrap/internal/user/models"
	"log/slog"
)

type Application struct {
	Commands    Commands
	Queries     Queries
	Logger      *slog.Logger
	UserService *models.Service
}
type Commands struct{}
type Queries struct{}
