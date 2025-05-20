package logger_test

import (
	"context"
	"log/slog"
	"testing"

	"github.com/stretchr/testify/assert"

	internalLogger "github.com/rubemlrm/go-api-bootstrap/internal/common/logger"
)

type testHandler struct {
	records []slog.Record
}

func (h *testHandler) Enabled(_ context.Context, _ slog.Level) bool { return true }
func (h *testHandler) Handle(_ context.Context, r slog.Record) error {
	h.records = append(h.records, r)
	return nil
}
func (h *testHandler) WithAttrs(attrs []slog.Attr) slog.Handler { return h }
func (h *testHandler) WithGroup(name string) slog.Handler       { return h }

func TestNewLogger(t *testing.T) {
	var tests = []struct {
		name   string
		level  string
		format string
	}{
		{
			name: "Create Logger with default options",
		},
		{
			name:   "Create Logger with json options",
			level:  "Debug",
			format: "json",
		},
		{
			name:   "Create Logger with text options",
			level:  "Debug",
			format: "text",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger := internalLogger.NewLogger(
				internalLogger.WithLogLevel(tt.level),
				internalLogger.WithLogFormat(tt.format))
			if logger == nil {
				t.Errorf("NewLogger() = nil, want non-nil")
			}
		})
	}
}

func TestInternalLogger_Info(t *testing.T) {
	h := &testHandler{}
	l := &internalLogger.SlogWrapper{
		Logger: slog.New(h),
	}
	l.Info("test message", "key", "value")

	assert.Len(t, h.records, 1)
	assert.Equal(t, "test message", h.records[0].Message)
}

func TestInternalLogger_Error(t *testing.T) {
	h := &testHandler{}
	l := &internalLogger.SlogWrapper{
		Logger: slog.New(h),
	}
	l.Error("test message", "key", "value")

	assert.Len(t, h.records, 1)
	assert.Equal(t, "test message", h.records[0].Message)
}

func TestInternalLogger_Debug(t *testing.T) {
	h := &testHandler{}
	l := &internalLogger.SlogWrapper{
		Logger: slog.New(h),
	}
	l.Debug("test message", "key", "value")

	assert.Len(t, h.records, 1)
	assert.Equal(t, "test message", h.records[0].Message)
}
