package logger

import (
	"os"

	"github.com/rubemlrm/go-api-bootstrap/config"
	"log/slog"
)

func NewLogger(cfg config.Logger) *slog.Logger {
	opts := slog.HandlerOptions{
		Level: setLogLevel(cfg.Level),
	}

	if cfg.Handler == "textHandler" {
		return setTextLogger(opts)

	}
	return setJsonLogger(opts)
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

func setJsonLogger(cfg slog.HandlerOptions) *slog.Logger {
	handler := slog.NewJSONHandler(os.Stdout, &cfg)
	return slog.New(handler)
}
