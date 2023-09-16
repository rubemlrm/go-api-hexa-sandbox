package main

import (
	"github.com/rubemlrm/go-api-bootstrap/config"
	"github.com/rubemlrm/go-api-bootstrap/internal/api"
	gin_handler "github.com/rubemlrm/go-api-bootstrap/internal/http/gin"
	"github.com/rubemlrm/go-api-bootstrap/pkg/postgres"
	"github.com/rubemlrm/go-api-bootstrap/pkg/slog"
	"github.com/rubemlrm/go-api-bootstrap/user"
	user_postgres "github.com/rubemlrm/go-api-bootstrap/user/postgres"
)

func main() {
	cfg, err := config.LoadConfig("config")
	if err != nil {
		panic(err)
	}
	logger := slog.NewLogger(cfg.Logger)

	logger.Info("app starting")
	db := postgres.StartConnection(cfg)
	repo := user_postgres.NewConnection(db)
	_ = user.NewService(repo)

	err = startWeb(cfg.HTTP)

	if err != nil {
		panic(err)
	}
}

func startWeb(httpConfig config.HTTP) error {
	ne := gin_handler.NewEngine()
	ne.SetHandlers()
	srv, err := api.NewServer(ne.StartHTTP(), httpConfig)

	if err != nil {
		panic(err)
	}

	err = srv.Start()
	if err != nil {
		panic(err)
	}
	return nil
}
