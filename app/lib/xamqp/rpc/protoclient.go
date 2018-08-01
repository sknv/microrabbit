package rpc

import (
	"github.com/golang/protobuf/proto"
	"github.com/pkg/errors"
	"github.com/streadway/amqp"
)

type ProtoClient struct {
	*Client
}

func NewProtoClient(conn *amqp.Connection) *ProtoClient {
	return &ProtoClient{Client: NewClient(conn)}
}

func (c *ProtoClient) Call(method string, args proto.Message, reply proto.Message) error {
	data, err := proto.Marshal(args)
	if err != nil {
		return errors.Wrap(err, "failed to marshal args to protobuf")
	}

	msg, err := c.Client.Call(method, &amqp.Publishing{Body: data})
	if err != nil {
		return errors.Wrapf(err, "failed to call a remote method: %s", method)
	}

	if err = proto.Unmarshal(msg.Body, reply); err != nil {
		return errors.Wrap(err, "failed to unmarshal a reply from protobuf")
	}
	return nil
}
