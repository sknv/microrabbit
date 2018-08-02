package rpc

import (
	"context"

	"github.com/pkg/errors"
	"github.com/streadway/amqp"

	"github.com/sknv/microrabbit/app/lib/rmq/rpc"
)

const (
	PatternCircle = "/rpc/math/circle"
)

type MathClient struct {
	*rpc.ProtoClient
}

func NewClient(rconn *amqp.Connection) *MathClient {
	return &MathClient{ProtoClient: rpc.NewProtoClient(rconn)}
}

func (c *MathClient) Circle(_ context.Context, args *CircleArgs) (*CircleReply, error) {
	var reply CircleReply
	if err := c.Call(PatternCircle, args, &reply); err != nil {
		return nil, errors.Wrap(err, "failed to call Math.Circle")
	}
	return &reply, nil
}
