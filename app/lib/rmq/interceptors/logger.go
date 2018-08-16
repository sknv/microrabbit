package interceptors

import (
	"context"
	"log"
	"time"

	"github.com/streadway/amqp"

	"github.com/sknv/microrabbit/app/lib/rmq"
)

func WithLogger(next rmq.HandlerFunc) rmq.HandlerFunc {
	return func(ctx context.Context, msg *amqp.Delivery) error {
		start := time.Now()
		defer func() {
			log.Printf("[INFO] request %s processed in %s", msg.RoutingKey, time.Since(start))
		}()
		return next(ctx, msg)
	}
}
