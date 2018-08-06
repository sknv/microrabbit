package rpc

import (
	"context"

	"github.com/streadway/amqp"
)

type contextKey string

const (
	headersKey contextKey = "rmq.headers"
)

func WithHeaders(ctx context.Context, headers amqp.Table) context.Context {
	return context.WithValue(ctx, headersKey, headers)
}

func Headers(ctx context.Context) (amqp.Table, bool) {
	headers, exist := ctx.Value(headersKey).(amqp.Table)
	return headers, exist
}
