package main

import (
	"context"
	"log"
	"log/slog"

	"github.com/rubemlrm/go-api-bootstrap/internal/common/api"
	"github.com/rubemlrm/go-api-bootstrap/internal/common/app"
	"github.com/rubemlrm/go-api-bootstrap/internal/common/config"
	ginhandler "github.com/rubemlrm/go-api-bootstrap/internal/common/http/gin"
	"github.com/rubemlrm/go-api-bootstrap/internal/common/logger"
	"github.com/rubemlrm/go-api-bootstrap/internal/common/postgres"
	"github.com/rubemlrm/go-api-bootstrap/internal/common/tracing"
	"github.com/rubemlrm/go-api-bootstrap/internal/user/ports"
	user_service "github.com/rubemlrm/go-api-bootstrap/internal/user/service"
)

func main() {
	cfg, err := config.LoadConfig("config")

	if err != nil {
		panic(err)
	}

	tp, err := tracing.InitTracer(cfg.Tracing)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := tp.Shutdown(context.Background()); err != nil {
			log.Printf("Error shutting down tracer provider: %v", err)
		}
	}()

	//metrics.InitMeter()

	l := logger.NewLogger(logger.WithLogFormat(cfg.Logger.Handler), logger.WithLogLevel(cfg.Logger.Level))

	l.Info("app starting")

	db, err := postgres.NewConnection(
		l.Logger,
		postgres.WithUsername(cfg.Database.User),
		postgres.WithPassword(cfg.Database.Password),
		postgres.WithHost(cfg.Database.Host),
		postgres.WithPort(cfg.Database.Port),
		postgres.WithSchema(cfg.Database.Schema),
		postgres.WithSSLMode(cfg.Database.SSLMode))
	if err != nil {
		panic(err)
	}
	um := user_service.NewApplication(context.Background(), l.Logger, db)
	err = startWeb(cfg.HTTP, l.Logger, app.Application{UserModule: um}, cfg.App.Name)

	if err != nil {
		panic(err)
	}
}

func startWeb(httpConfig config.HTTP, logger *slog.Logger, app app.Application, appName string) error {
	ne := ginhandler.NewEngine()
	ne.SetHandlers(logger, func() {
		opt := ports.GinServerOptions{
			BaseURL: "/api/v1",
		}
		ports.RegisterHandlersWithOptions(ne.Engine, ports.NewHTTPServer(app, logger), opt)
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
