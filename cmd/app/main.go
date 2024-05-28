package main

import (
	"github.com/rubemlrm/go-api-bootstrap/config"
	"github.com/rubemlrm/go-api-bootstrap/internal/api"
	gin_handler "github.com/rubemlrm/go-api-bootstrap/internal/http/gin"
	"github.com/rubemlrm/go-api-bootstrap/pkg/postgres"
	"github.com/rubemlrm/go-api-bootstrap/pkg/slog"
	"github.com/rubemlrm/go-api-bootstrap/user"
	user_postgres "github.com/rubemlrm/go-api-bootstrap/user/postgres"
	slogger "golang.org/x/exp/slog"
)

func main() {
	cfg, err := config.LoadConfig("config")
	if err != nil {
		panic(err)
	}
	logger := slog.NewLogger(cfg.Logger)

	logger.Info("app starting")
	db := postgres.StartConnection(cfg, logger)
	repo := user_postgres.NewConnection(db, logger)
	us := user.NewService(repo, logger)

	err = startWeb(cfg.HTTP, us, logger)

	if err != nil {
		panic(err)
	}
}

func startWeb(httpConfig config.HTTP, userService *user.Service, logger *slogger.Logger) error {
	ne := gin_handler.NewEngine()
	ne.SetHandlers(userService, logger)
	srv, err := api.NewServer(ne.StartHTTP(), httpConfig, logger)

	if err != nil {
		panic(err)
	}

	err = srv.Start()
	if err != nil {
		panic(err)
	}
	return nil
}
