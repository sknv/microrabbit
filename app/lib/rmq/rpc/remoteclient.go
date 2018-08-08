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
	ch, err := rmq.NewChannel(c.Conn)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to open a channel for the remote client")
	}
	defer ch.Close()

	// consume from the rpc reply queue
	msgs, err := ch.Consume(rpcReplyQueue, true)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to register a consumer for the remote client")
	}

	// prepare a message and publish it
	correlationID := xstrings.RandomString(32)
	msg := prepareMessage(ctx, message, correlationID, rpcReplyQueue)
	if err = ch.Publish("", method, msg); err != nil {
		return nil, errors.WithMessage(err, fmt.Sprintf("failed to publish a message from the remote client to %s", method))
	}
	return handleReply(ctx, msgs, correlationID)
}

// ----------------------------------------------------------------------------
// ----------------------------------------------------------------------------
// ----------------------------------------------------------------------------

func prepareMessage(ctx context.Context, message *amqp.Publishing, correlationID, replyTo string) *amqp.Publishing {
	message.CorrelationId = correlationID
	message.ReplyTo = replyTo
	message = expireIfNeeded(ctx, message)  // expire a message if a deadline specified for the context
	message = addMetaIfNeeded(ctx, message) // add metadata to the message if exist in the context
	return message
}

func expireIfNeeded(ctx context.Context, message *amqp.Publishing) *amqp.Publishing {
	timeout, exist := xcontext.Timeout(ctx)
	if !exist {
		return message
	}

	exp := int(timeout.Seconds() * 1000) // expiration in ms
	message.Expiration = fmt.Sprint(exp)
	return message
}

func addMetaIfNeeded(ctx context.Context, message *amqp.Publishing) *amqp.Publishing {
	meta, exist := Meta(ctx)
	if !exist {
		return message
	}

	message.Headers = meta
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
			return nil, errors.WithMessage(err, "failed to handle a reply for the remote client")
		}
	}
}
