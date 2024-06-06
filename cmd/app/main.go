package main

import (
	"log/slog"

	"github.com/rubemlrm/go-api-bootstrap/config"
	"github.com/rubemlrm/go-api-bootstrap/internal/api"
	gin_handler "github.com/rubemlrm/go-api-bootstrap/internal/http/gin"
	"github.com/rubemlrm/go-api-bootstrap/pkg/logger"
	"github.com/rubemlrm/go-api-bootstrap/pkg/postgres"
	"github.com/rubemlrm/go-api-bootstrap/user"
	user_postgres "github.com/rubemlrm/go-api-bootstrap/user/postgres"
)

func main() {
	cfg, err := config.LoadConfig("config")
	if err != nil {
		panic(err)
	}
	l := logger.NewLogger(cfg.Logger)

	l.Info("app starting")
	db := postgres.StartConnection(cfg, l)
	repo := user_postgres.NewConnection(db, l)
	us := user.NewService(repo, l)

	err = startWeb(cfg.HTTP, us, l)

	if err != nil {
		panic(err)
	}
}

func startWeb(httpConfig config.HTTP, userService *user.Service, logger *slog.Logger) error {
	ne := gin_handler.NewEngine()
	ne.SetHandlers(userService, logger)
	srv, err := api.NewServer(ne.StartHTTP(), httpConfig, logger)

	if err != nil {
		return err
	}

	err = srv.Start()
	if err != nil {
		return err
	}
	return nil
}
