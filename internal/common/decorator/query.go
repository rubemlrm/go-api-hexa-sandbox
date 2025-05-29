package decorator

import (
	"context"
	"log/slog"
)

func ApplyQueryDecorators[Q any, R any](handler QueryHandler[Q, R], logger *slog.Logger) QueryHandler[Q, R] {
	return DatabaseLoggingQueryDecorator[Q, R]{
		base:   handler,
		logger: logger,
	}
}

type QueryHandler[Q any, R any] interface {
	Handle(ctx context.Context, q Q) (R, error)
}
