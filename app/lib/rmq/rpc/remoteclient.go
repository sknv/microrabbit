package rpc

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
	"github.com/streadway/amqp"

	"github.com/sknv/microrabbit/app/lib/rmq"
	"github.com/sknv/microrabbit/app/lib/rmq/status"
)

type RemoteClient struct {
	Conn *rmq.Connection
}

func NewRemoteClient(conn *rmq.Connection) *RemoteClient {
	return &RemoteClient{Conn: conn}
}

func (c *RemoteClient) Call(ctx context.Context, method string, message *amqp.Publishing) (*amqp.Delivery, error) {
	reply, err := c.Conn.Request(ctx, "", method, message)
	if err != nil { // handle network errors
		cause := errors.Cause(err)
		if cause == context.DeadlineExceeded { // handle timeout error if such exist
			cause = status.Error(status.DeadlineExceeded, err.Error())
		}
		return nil, errors.WithMessage(cause, fmt.Sprintf("failed to call %s", method))
	}
	return reply, nil
}
