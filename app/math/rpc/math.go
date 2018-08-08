package rpc

import (
	"context"
)

const (
	CirclePattern = "/rpc/math/circle"
	RectPattern   = "/rpc/math/rect"
)

type Math interface {
	Circle(context.Context, *CircleArgs) (*CircleReply, error)
	Rect(context.Context, *RectArgs) (*RectReply, error)
}
