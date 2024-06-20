package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rubemlrm/go-api-bootstrap/config"
	"github.com/rubemlrm/go-api-bootstrap/internal/api"
	ginhandler "github.com/rubemlrm/go-api-bootstrap/internal/http/gin"
	"github.com/rubemlrm/go-api-bootstrap/pkg/logger"
	"github.com/rubemlrm/go-api-bootstrap/pkg/postgres"
	"github.com/rubemlrm/go-api-bootstrap/user"
	userpostgres "github.com/rubemlrm/go-api-bootstrap/user/postgres"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/exporters/prometheus"
	apimetric "go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/propagation"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.25.0"
)

func main() {
	tp, err := initTracer()
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

	initMeter()

	l := logger.NewLogger(cfg.Logger)

	l.Info("app starting")
	db := postgres.StartConnection(cfg, l)
	repo := userpostgres.NewConnection(db, l)
	us := user.NewService(repo, l)

	err = startWeb(cfg.HTTP, us, l)

	if err != nil {
		panic(err)
	}
}

func startWeb(httpConfig config.HTTP, userService *user.Service, logger *slog.Logger) error {
	ne := ginhandler.NewEngine()
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

func initTracer() (*sdktrace.TracerProvider, error) {
	exporter, err := otlptrace.New(
		context.Background(),
		otlptracegrpc.NewClient(
			otlptracegrpc.WithInsecure(),
			otlptracegrpc.WithEndpoint("jaeger.box.rubemlrm.com:4317"),
		),
	)
	if err != nil {
		return nil, err
	}
	ctx := context.Background()
	traceRes, err := resource.New(ctx,
		resource.WithAttributes(semconv.ServiceName("my-service"), semconv.ServiceVersion("0.1.0")),
	)

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithResource(traceRes),
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithBatcher(exporter),
	)
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))
	return tp, nil
}

func initMeter() {
	ctx := context.Background()
	exporter, err := prometheus.New()
	if err != nil {
		log.Fatal(err)
	}
	provider := sdkmetric.NewMeterProvider(sdkmetric.WithReader(exporter))
	meter := provider.Meter("sandbox")
	counter, err := meter.Float64Counter("foo", apimetric.WithDescription("a simple counter"))
	if err != nil {
		log.Fatal(err)
	}
	counter.Add(ctx, 5, apimetric.WithAttributes(
		attribute.Key("A").String("B"),
		attribute.Key("C").String("D"),
	))
	ctx, _ = signal.NotifyContext(ctx, os.Interrupt)
	<-ctx.Done()
	// Start the prometheus HTTP server and pass the exporter Collector to it
	go func() {
		log.Printf("serving metrics at localhost:2223/metrics")
		http.Handle("/metrics", promhttp.Handler())
		err := http.ListenAndServe(":2223", nil) //nolint:gosec // Ignoring G114: Use of net/http serve function that has no support for setting timeouts.
		if err != nil {
			fmt.Printf("error serving http: %v", err)
			return
		}
	}()
}
