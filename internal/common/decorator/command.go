package decorator

import (
	"context"
	"log/slog"
)

func ApplyCommandDecorators[Q any, R any](handler QueryHandler[Q, R], logger *slog.Logger) QueryHandler[Q, R] {
	return DatabaseLoggingCommandDecorator[Q, R]{
		base:   handler,
		logger: logger,
	}
}

type CommandHandler[C any, Q any] interface {
	Handle(ctx context.Context, cmd C) (Q, error)
}
