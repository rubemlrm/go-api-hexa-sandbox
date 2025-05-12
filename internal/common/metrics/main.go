package metrics

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/prometheus"
	apimetric "go.opentelemetry.io/otel/metric"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
)

func InitMeter() {
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
