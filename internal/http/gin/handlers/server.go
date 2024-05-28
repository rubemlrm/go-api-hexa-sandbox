package handlers

import (
	"github.com/rubemlrm/go-api-bootstrap/internal/http/gin/openapi"
	"github.com/rubemlrm/go-api-bootstrap/user"
	"golang.org/x/exp/slog"
)

type server struct {
	UserService *user.Service
	Logger      *slog.Logger
}

func NewServer(userService *user.Service, logger *slog.Logger) openapi.ServerInterface {
	return &server{
		UserService: userService,
		Logger:      logger,
	}
}
