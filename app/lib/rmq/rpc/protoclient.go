package rpc

import (
	"context"

	"github.com/golang/protobuf/proto"
	"github.com/pkg/errors"
	"github.com/streadway/amqp"

	"github.com/sknv/microrabbit/app/lib/rmq"
	"github.com/sknv/microrabbit/app/lib/rmq/status"
)

type ProtoClient struct {
	*RemoteClient
}

func NewProtoClient(rmqConn *rmq.Connection) *ProtoClient {
	return &ProtoClient{RemoteClient: NewRemoteClient(rmqConn)}
}

func (c *ProtoClient) Call(ctx context.Context, method string, args proto.Message, reply proto.Message) error {
	data, err := proto.Marshal(args)
	if err != nil {
		return errors.WithMessage(err, "failed to marshal args to protobuf")
	}

	msg, err := c.RemoteClient.Call(ctx, method, &amqp.Publishing{Body: data})
	if err != nil {
		return errors.WithMessage(err, "failed to call the remote method")
	}

	// handle an error transfered over the network
	if rmq.HasError(msg) {
		status := new(status.Status)
		if err = proto.Unmarshal(msg.Body, status); err != nil {
			return errors.WithMessage(err, "failed to unmarshal the error from protobuf")
		}
		return status
	}

	// handle a reply
	if err = proto.Unmarshal(msg.Body, reply); err != nil {
		return errors.WithMessage(err, "failed to unmarshal the reply from protobuf")
	}
	return nil
}
