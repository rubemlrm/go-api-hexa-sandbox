package decorator

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
	"time"
)

type DatabaseLoggingCommandDecorator[Q any, R any] struct {
	base   CommandHandler[Q, R]
	logger *slog.Logger
}

func (d DatabaseLoggingCommandDecorator[C, R]) Handle(ctx context.Context, cmd C) (result R, err error) {
	start := time.Now()
	loggerAttrs := []any{
		"operation_type", "database",
		"command", generateActionName(cmd),
		"command_body", fmt.Sprintf("%#v", cmd),
		"requestID", ctx.Value("requestID"),
	}

	d.logger.Debug("Executing command")
	defer func(attrs []any, start time.Time) {
		loggerAttrs := append(attrs, "duration", time.Since(start))
		if err == nil {
			d.logger.Info("Command executed successfully", loggerAttrs...)
		} else {
			loggerAttrs = append(loggerAttrs, "error", err)
			d.logger.Error("Command execution failed", loggerAttrs...)
		}
	}(loggerAttrs, start)

	return d.base.Handle(ctx, cmd)
}

type DatabaseLoggingQueryDecorator[Q any, R any] struct {
	base   QueryHandler[Q, R]
	logger *slog.Logger
}

func (d DatabaseLoggingQueryDecorator[C, R]) Handle(ctx context.Context, cmd C) (result R, err error) {
	start := time.Now()
	loggerAttrs := []any{
		"operation_type", "database",
		"query", generateActionName(cmd),
		"query_body", fmt.Sprintf("%#v", cmd),
		"requestID", ctx.Value("requestID"),
	}

	d.logger.Debug("Executing query")
	defer func(attrs []any, start time.Time) {
		loggerAttrs := append(attrs, "duration", time.Since(start))
		if err == nil {
			d.logger.Info("Query executed successfully", loggerAttrs...)
		} else {
			loggerAttrs = append(loggerAttrs, "error", err)
			d.logger.Error("Query execution failed", loggerAttrs...)
		}
	}(loggerAttrs, start)

	return d.base.Handle(ctx, cmd)
}

func generateActionName(handler any) string {
	return strings.Split(fmt.Sprintf("%T", handler), ".")[1]
}
