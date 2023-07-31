package main

import (
	"github.com/rubemlrm/go-api-bootstrap/config"
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
	service := user.NewService(repo)

	id, err := service.Get(user.ID(4))
	if err != nil {
		panic(err)
	}
	print(id)
}
