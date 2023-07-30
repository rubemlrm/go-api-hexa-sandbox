package main

import (
	"github.com/rubemlrm/go-api-bootstrap/config"
)

func main() {
	cfg, err := config.LoadConfig("config")
	if err != nil {
		panic(err)
	}
	print(cfg.App.Name)
}
