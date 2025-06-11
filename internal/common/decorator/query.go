package decorator

import (
	"context"
	"log/slog"
)

func ApplyQueryDecorators[Q any, R any](handler QueryHandler[Q, R], logger *slog.Logger, tracer RecordTracer) QueryHandler[Q, R] {
	return DatabaseLoggingQueryDecorator[Q, R]{
		base: DatabaseTracingQueryDecorator[Q, R]{
			base:   handler,
			tracer: tracer,
		},
		logger: logger,
	}
}

type QueryHandler[Q any, R any] interface {
	Handle(ctx context.Context, q Q) (R, error)
}
