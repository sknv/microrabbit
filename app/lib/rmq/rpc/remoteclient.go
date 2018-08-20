package rpc

import (
	"context"
	"fmt"
	"math"

	"github.com/pkg/errors"
	"github.com/streadway/amqp"

	"github.com/sknv/microrabbit/app/lib/rmq"
	"github.com/sknv/microrabbit/app/lib/rmq/status"
	"github.com/sknv/microrabbit/app/lib/xcontext"
	"github.com/sknv/microrabbit/app/lib/xstrings"
)

const (
	directReplyQueue = "amq.rabbitmq.reply-to"
)

type RemoteClient struct {
	RmqConn *rmq.Connection
}

func NewRemoteClient(rmqConn *rmq.Connection) *RemoteClient {
	return &RemoteClient{RmqConn: rmqConn}
}

func (c *RemoteClient) Call(ctx context.Context, method string, message *amqp.Publishing) (*amqp.Delivery, error) {
	// add the metadata from the context
	message = rmq.WithMetadata(message, rmq.ContextMetadata(ctx))

	reply, err := c.request(ctx, method, message)
	if err != nil { // handle network errors
		cause := errors.Cause(err)
		if cause == context.DeadlineExceeded { // handle timeout error if such exist
			cause = status.Error(status.DeadlineExceeded, err.Error())
		}
		return nil, errors.WithMessage(cause, fmt.Sprintf("failed to call %s", method))
	}
	return reply, nil
}

// ----------------------------------------------------------------------------
// ----------------------------------------------------------------------------
// ----------------------------------------------------------------------------

func (c *RemoteClient) request(ctx context.Context, routingKey string, message *amqp.Publishing) (*amqp.Delivery, error) {
	ch, err := c.RmqConn.OpenChannel()
	if err != nil {
		return nil, errors.WithMessage(err, "failed to open a sync consumer channel")
	}
	defer ch.Close()

	// prepare the message
	correlationID := xstrings.RandomString(32)
	msg := prepareMessage(ctx, message, correlationID)

	// consume from the reply queue
	msgs, err := ch.ConsumeFrom(msg.ReplyTo, true)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to register a sync consumer")
	}

	// publish the message
	if err = ch.PublishTo("", routingKey, msg); err != nil {
		return nil, errors.WithMessage(err, fmt.Sprintf("failed to publish the sync message to %s", routingKey))
	}
	return handleReply(ctx, msgs, correlationID) // wait for the reply
}

// ----------------------------------------------------------------------------
// ----------------------------------------------------------------------------
// ----------------------------------------------------------------------------

func prepareMessage(ctx context.Context, message *amqp.Publishing, correlationID string) *amqp.Publishing {
	message.ReplyTo = directReplyQueue
	message.CorrelationId = correlationID
	message = expireIfNeeded(ctx, message) // expire a message if a deadline specified for the context
	return message
}

func expireIfNeeded(ctx context.Context, message *amqp.Publishing) *amqp.Publishing {
	timeout, exist := xcontext.Timeout(ctx)
	if !exist {
		return message
	}

	exp := math.Ceil(timeout.Seconds() * 1000) // expiration in ms rounded to the nearest greater int
	message.Expiration = fmt.Sprint(exp)
	return message
}

func handleReply(ctx context.Context, messages <-chan amqp.Delivery, correlationID string) (*amqp.Delivery, error) {
	for {
		select {
		case <-ctx.Done():
			return nil, errors.WithMessage(ctx.Err(), "failed to wait for the reply")
		case msg := <-messages:
			if correlationID == msg.CorrelationId {
				return &msg, nil
			}
		}
	}
}
