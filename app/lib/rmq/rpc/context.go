package rpc

import (
	"context"
	"time"

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
	headers, ok := ctx.Value(headersKey).(amqp.Table)
	return headers, ok
}

func Timeout(ctx context.Context) (time.Duration, bool) {
	deadline, ok := ctx.Deadline()
	if !ok {
		return 0, false
	}
	return deadline.Sub(time.Now()), true
}
