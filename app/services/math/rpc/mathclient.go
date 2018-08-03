package rpc

import (
	"context"

	"github.com/pkg/errors"
	"github.com/streadway/amqp"

	"github.com/sknv/microrabbit/app/lib/rmq/rpc"
)

const (
	CirclePattern = "/rpc/math/circle"
)

type MathClient struct {
	*rpc.ProtoClient
}

func NewClient(rconn *amqp.Connection) *MathClient {
	return &MathClient{ProtoClient: rpc.NewProtoClient(rconn)}
}

func (c *MathClient) Circle(ctx context.Context, args *CircleArgs) (*CircleReply, error) {
	reply := new(CircleReply)
	if err := c.Call(ctx, CirclePattern, args, reply); err != nil {
		return nil, errors.Wrap(err, "failed to call Math.Circle")
	}
	return reply, nil
}
