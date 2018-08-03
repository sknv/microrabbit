package rpc

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
	"github.com/streadway/amqp"

	"github.com/sknv/microrabbit/app/lib/rmq"
	"github.com/sknv/microrabbit/app/lib/rmq/status"
	"github.com/sknv/microrabbit/app/lib/xcontext"
	"github.com/sknv/microrabbit/app/lib/xstrings"
)

const (
	rpcReplyQueue = "amq.rabbitmq.reply-to"
)

type RemoteClient struct {
	Conn *amqp.Connection
}

func NewRemoteClient(conn *amqp.Connection) *RemoteClient {
	return &RemoteClient{Conn: conn}
}

func (c *RemoteClient) Call(ctx context.Context, method string, message *amqp.Publishing) (*amqp.Delivery, error) {
	ch, err := c.Conn.Channel()
	if err != nil {
		return nil, errors.Wrap(err, "failed to open a channel")
	}
	channel := rmq.NewChannel(ch)
	defer channel.Close()

	// consume from the rpc reply queue
	messages, err := channel.Consume(rpcReplyQueue, true)
	if err != nil {
		return nil, errors.Wrap(err, "failed to register a consumer")
	}

	// prepare a message and publish it
	correlationID := xstrings.RandomString(32)
	msg := prepareMessage(ctx, message, correlationID, rpcReplyQueue)
	if err = channel.Publish("", method, msg); err != nil {
		return nil, errors.Wrap(err, "failed to publish a message")
	}
	return handleReply(ctx, messages, correlationID)
}

// ----------------------------------------------------------------------------
// ----------------------------------------------------------------------------
// ----------------------------------------------------------------------------

func prepareMessage(ctx context.Context, message *amqp.Publishing, correlationID, replyTo string) *amqp.Publishing {
	message.CorrelationId = correlationID
	message.ReplyTo = replyTo
	message = expireIfNeeded(ctx, message)     // expire a message if a deadline specified for the context
	message = addHeadersIfNeeded(ctx, message) // add headers table if one exists in the context
	return message
}

func expireIfNeeded(ctx context.Context, message *amqp.Publishing) *amqp.Publishing {
	timeout, ok := xcontext.Timeout(ctx)
	if !ok {
		return message
	}

	exp := timeout.Seconds() * 1000 // expiration in ms
	message.Expiration = fmt.Sprint(exp)
	return message
}

func addHeadersIfNeeded(ctx context.Context, message *amqp.Publishing) *amqp.Publishing {
	headers, ok := Headers(ctx)
	if !ok {
		return message
	}

	message.Headers = headers
	return message
}

func handleReply(ctx context.Context, messages <-chan amqp.Delivery, correlationID string) (*amqp.Delivery, error) {
	for {
		select {
		case msg := <-messages:
			if correlationID == msg.CorrelationId {
				return &msg, nil
			}
		case <-ctx.Done():
			err := status.Error(status.DeadlineExceeded, ctx.Err().Error())
			return nil, errors.Wrap(err, "failed to handle a reply")
		}
	}
}
