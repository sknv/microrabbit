package rpc

import (
	"context"
	"fmt"

	"github.com/golang/protobuf/proto"
	"github.com/pkg/errors"
	"github.com/streadway/amqp"

	"github.com/sknv/microrabbit/app/lib/rmq/status"
)

type ProtoClient struct {
	*RemoteClient
}

func NewProtoClient(conn *amqp.Connection) *ProtoClient {
	return &ProtoClient{RemoteClient: NewRemoteClient(conn)}
}

func (c *ProtoClient) Call(ctx context.Context, method string, args proto.Message, reply proto.Message) error {
	data, err := proto.Marshal(args)
	if err != nil {
		return errors.WithMessage(err, "failed to marshal args to protobuf")
	}

	msg, err := c.RemoteClient.Call(ctx, method, &amqp.Publishing{Body: data})
	if err != nil {
		return errors.WithMessage(err, fmt.Sprintf("failed to call a remote method: %s", method))
	}

	// handle an error transfered over the network
	if HasError(msg) {
		rerr := new(status.Status)
		if err = proto.Unmarshal(msg.Body, rerr); err != nil {
			return errors.WithMessage(err, "failed to unmarshal an error from protobuf")
		}
		return rerr
	}

	// handle a reply
	if err = proto.Unmarshal(msg.Body, reply); err != nil {
		return errors.WithMessage(err, "failed to unmarshal a reply from protobuf")
	}
	return nil
}
