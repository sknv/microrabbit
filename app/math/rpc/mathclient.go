package rpc

import (
	"context"

	"github.com/sknv/microrabbit/app/lib/rmq"
	"github.com/sknv/microrabbit/app/lib/rmq/rpc"
)

type MathClient struct {
	*rpc.ProtoClient
}

func NewClient(conn *rmq.Connection) Math {
	return &MathClient{ProtoClient: rpc.NewProtoClient(conn)}
}

func (c *MathClient) Circle(ctx context.Context, args *CircleArgs) (*CircleReply, error) {
	reply := new(CircleReply)
	if err := c.Call(ctx, CirclePattern, args, reply); err != nil {
		return nil, err
	}
	return reply, nil
}

func (c *MathClient) Rect(ctx context.Context, args *RectArgs) (*RectReply, error) {
	reply := new(RectReply)
	if err := c.Call(ctx, RectPattern, args, reply); err != nil {
		return nil, err
	}
	return reply, nil
}
