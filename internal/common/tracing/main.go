package tracing

import (
	"context"
	"fmt"

	"github.com/rubemlrm/go-api-bootstrap/internal/common/config"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.25.0"
)

func InitTracer(tracing config.Tracing) (*sdktrace.TracerProvider, error) {
	host := fmt.Sprintf("%s:%s", tracing.AgentHost, tracing.AgentPort)
	exporter, err := otlptrace.New(
		context.Background(),
		otlptracegrpc.NewClient(
			otlptracegrpc.WithInsecure(),
			otlptracegrpc.WithEndpoint(host),
		),
	)
	if err != nil {
		return nil, err
	}
	ctx := context.Background()
	traceRes, err := resource.New(ctx,
		resource.WithAttributes(semconv.ServiceName(tracing.ServiceName), semconv.ServiceVersion("0.1.0")),
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
