package rpc

import (
	"context"

	"github.com/pkg/errors"
	"github.com/streadway/amqp"

	"github.com/sknv/microrabbit/app/lib/rmq"
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

	// add headers table if exists
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
			err := rmq.WrapError(rmq.StatusDeadlineExceeded, ctx.Err())
			return nil, errors.Wrap(err, "failed to handle a reply")
		}
	}
}
