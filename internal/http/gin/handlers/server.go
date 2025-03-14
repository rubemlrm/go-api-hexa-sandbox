package handlers

import (
	"github.com/rubemlrm/go-api-bootstrap/user"
	"log/slog"

	"github.com/rubemlrm/go-api-bootstrap/internal/http/gin/openapi"
)

type server struct {
	UserService user.UseCase
	Logger      *slog.Logger
}

func NewServer(userService user.UseCase, logger *slog.Logger) openapi.ServerInterface {
	return &server{
		UserService: userService,
		Logger:      logger,
	}
}
