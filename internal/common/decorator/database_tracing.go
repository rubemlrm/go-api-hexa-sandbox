package decorator

import (
	"context"
)

type RecordTracer interface {
	RecordTrace(ctx context.Context, actionName string, tracerName string)
}
type DatabaseTracingCommandDecorator[C any, Q any] struct {
	base   CommandHandler[C, Q]
	tracer RecordTracer
}

func (d DatabaseTracingCommandDecorator[C, Q]) Handle(ctx context.Context, cmd C) (result Q, err error) {
	d.tracer.RecordTrace(ctx, generateActionName(cmd), "database-tracing-command")
	return d.base.Handle(ctx, cmd)
}

type DatabaseTracingQueryDecorator[Q any, R any] struct {
	base   QueryHandler[Q, R]
	tracer RecordTracer
}

func (d DatabaseTracingQueryDecorator[Q, R]) Handle(ctx context.Context, q Q) (result R, err error) {
	d.tracer.RecordTrace(ctx, generateActionName(q), "database-tracing-query")
	return d.base.Handle(ctx, q)
}
