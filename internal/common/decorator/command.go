package decorator

import "context"

type CommandHandler[C any, Q any] interface {
	Handle(ctx context.Context, cmd C) (Q, error)
}
