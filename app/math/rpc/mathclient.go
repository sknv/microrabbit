package rpc

import (
	"context"

	"github.com/pkg/errors"
	"github.com/streadway/amqp"

	"github.com/sknv/microrabbit/app/lib/rmq/rpc"
)

type MathClient struct {
	*rpc.ProtoClient
}

func NewClient(rconn *amqp.Connection) Math {
	return &MathClient{ProtoClient: rpc.NewProtoClient(rconn)}
}

func (c *MathClient) Circle(ctx context.Context, args *CircleArgs) (*CircleReply, error) {
	reply := new(CircleReply)
	if err := c.Call(ctx, CirclePattern, args, reply); err != nil {
		return nil, errors.WithMessage(err, "failed to call Math.Circle")
	}
	return reply, nil
}

func (c *MathClient) Rect(ctx context.Context, args *RectArgs) (*RectReply, error) {
	reply := new(RectReply)
	if err := c.Call(ctx, RectPattern, args, reply); err != nil {
		return nil, errors.WithMessage(err, "failed to call Math.Rect")
	}
	return reply, nil
}
