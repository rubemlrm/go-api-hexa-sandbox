package tracing

import (
	"context"
	"github.com/rubemlrm/go-api-bootstrap/internal/common/config"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestInitializesTracerProviderSuccessfully(t *testing.T) {
	tracingConfig := config.Tracing{
		AgentHost:   "localhost",
		AgentPort:   "4317",
		ServiceName: "test-service",
	}
	tp, err := InitTracer(tracingConfig)

	assert.NoError(t, err)
	assert.NotNil(t, tp)
}

func TestRecordsTraceSuccessfully(t *testing.T) {
	tracerProvider := &TracerProvider{}
	ctx := context.Background()
	tracerProvider.RecordTrace(ctx, "test-action", "test-tracer")
}
