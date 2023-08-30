package main

import (
	"github.com/rubemlrm/go-api-bootstrap/config"
	"github.com/rubemlrm/go-api-bootstrap/internal/api"
	gin_handler "github.com/rubemlrm/go-api-bootstrap/internal/http/gin"
	"github.com/rubemlrm/go-api-bootstrap/pkg/postgres"
	"github.com/rubemlrm/go-api-bootstrap/user"
	user_postgres "github.com/rubemlrm/go-api-bootstrap/user/postgres"
)

func main() {
	cfg, err := config.LoadConfig("config")
	if err != nil {
		panic(err)
	}
	db := postgres.StartConnection(cfg)
	repo := user_postgres.NewConnection(db)
	_ = user.NewService(repo)

	err = startWeb()

	if err != nil {
		panic(err)
	}
}

func startWeb() error {
	ne := gin_handler.NewEngine()
	ne.SetHandlers()
	err := api.Start(ne.StartHTTP())
	if err != nil {
		panic(err)
	}
	return nil
}
