package decorator

import (
	"context"
	"log/slog"
)

func ApplyCommandDecorators[Q any, R any](handler QueryHandler[Q, R], logger *slog.Logger, tracer RecordTracer) QueryHandler[Q, R] {
	return DatabaseLoggingCommandDecorator[Q, R]{
		base: DatabaseTracingCommandDecorator[Q, R]{
			base:   handler,
			tracer: tracer,
		},
		logger: logger,
	}
}

type CommandHandler[C any, Q any] interface {
	Handle(ctx context.Context, cmd C) (Q, error)
}
