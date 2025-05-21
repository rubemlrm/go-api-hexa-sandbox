package logger

import (
	"os"

	"log/slog"
)

type Logger interface {
	Info(msg string, args ...any)
	Error(msg string, args ...any)
	Debug(msg string, args ...any)
}

type SlogWrapper struct {
	*slog.Logger
	level  slog.Level
	format string
}

type LoggerOption func(*SlogWrapper)

func NewLogger(options ...LoggerOption) *SlogWrapper {
	// Default logger configuration
	l := &SlogWrapper{
		level:  slog.LevelInfo,
		format: "json",
	}

	// Apply options
	for _, option := range options {
		option(l)
	}

	var handler slog.Handler

	if l.format == "json" {
		handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: l.level,
		})
	} else {
		handler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: l.level,
		})
	}

	l.Logger = slog.New(handler)
	return l
}

func WithLogLevel(ht string) LoggerOption {
	return func(l *SlogWrapper) {
		if ht == "Debug" {
			l.level = slog.LevelDebug
		}
		l.level = slog.LevelInfo
	}
}

func WithLogFormat(format string) LoggerOption {
	return func(l *SlogWrapper) {
		l.format = format
	}
}

func (l *SlogWrapper) Info(msg string, args ...any) {
	l.Logger.Info(msg, args...)
}

func (l *SlogWrapper) Error(msg string, args ...any) {
	l.Logger.Error(msg, args...)
}

func (l *SlogWrapper) Debug(msg string, args ...any) {
	l.Logger.Debug(msg, args...)
}
