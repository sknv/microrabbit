package rpc

import (
	"context"

	"github.com/streadway/amqp"
)

type contextKey string

const (
	headersKey = contextKey("rmq.headers")
)

func WithHeadersTable(ctx context.Context, headers amqp.Table) context.Context {
	return context.WithValue(ctx, headersKey, headers)
}

func HeadersTable(ctx context.Context) (amqp.Table, bool) {
	headers, ok := ctx.Value(headersKey).(amqp.Table)
	return headers, ok
}
