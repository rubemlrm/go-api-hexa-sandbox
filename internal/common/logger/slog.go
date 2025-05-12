package logger

import (
	"os"

	"log/slog"

	"github.com/rubemlrm/go-api-bootstrap/internal/common/config"
)

func NewLogger(cfg config.Logger) *slog.Logger {
	opts := slog.HandlerOptions{
		Level: setLogLevel(cfg.Level),
	}

	if cfg.Handler == "textHandler" {
		return setTextLogger(opts)
	}
	return setJSONLogger(opts)
}

func setLogLevel(ht string) slog.Level {
	if ht == "Debug" {
		return slog.LevelDebug
	}
	return slog.LevelInfo
}

func setTextLogger(cfg slog.HandlerOptions) *slog.Logger {
	handler := slog.NewTextHandler(os.Stdout, &cfg)
	return slog.New(handler)
}

func setJSONLogger(cfg slog.HandlerOptions) *slog.Logger {
	handler := slog.NewJSONHandler(os.Stdout, &cfg)
	return slog.New(handler)
}
