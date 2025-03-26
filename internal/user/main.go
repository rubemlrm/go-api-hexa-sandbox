package main

import (
	"context"
	"github.com/rubemlrm/go-api-bootstrap/internal/common/api"
	"github.com/rubemlrm/go-api-bootstrap/internal/common/config"
	ginhandler "github.com/rubemlrm/go-api-bootstrap/internal/common/http/gin"
	"github.com/rubemlrm/go-api-bootstrap/internal/common/logger"
	"github.com/rubemlrm/go-api-bootstrap/internal/common/tracing"
	"github.com/rubemlrm/go-api-bootstrap/internal/user/app"
	"github.com/rubemlrm/go-api-bootstrap/internal/user/ports"
	"github.com/rubemlrm/go-api-bootstrap/internal/user/service"
	"log"
	"log/slog"
)

func main() {
	tp, err := tracing.InitTracer()
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := tp.Shutdown(context.Background()); err != nil {
			log.Printf("Error shutting down tracer provider: %v", err)
		}
	}()
	cfg, err := config.LoadConfig("config")

	if err != nil {
		panic(err)
	}

	//metrics.InitMeter()

	l := logger.NewLogger(cfg.Logger)

	l.Info("app starting")
	app := service.NewApplication(context.Background(), cfg, l)
	err = startWeb(cfg.HTTP, l, app)

	if err != nil {
		panic(err)
	}
}

func startWeb(httpConfig config.HTTP, logger *slog.Logger, app app.Application) error {
	ne := ginhandler.NewEngine()
	ne.SetHandlers(logger, func() {
		opt := ports.GinServerOptions{
			BaseURL: "/api/v1",
		}
		ports.RegisterHandlersWithOptions(ne.Engine, ports.NewHttpServer(app, logger), opt)
	})
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
