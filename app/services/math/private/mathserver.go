package private

import (
	"context"
	"math"

	"github.com/sknv/microrabbit/app/lib/xamqp"
	"github.com/sknv/microrabbit/app/services/math/public"
)

type MathServer struct{}

func (*MathServer) Rect(_ context.Context, args *public.RectArgs) (*public.RectReply, error) {
	if args.Width <= 0 || args.Height <= 0 {
		return nil, xamqp.NewError(xamqp.InvalidArgument, "width and height must be positive numbers")
	}

	return &public.RectReply{
		Perimeter: 2*args.Width + 2*args.Height,
		Square:    args.Width * args.Height,
	}, nil
}

func (*MathServer) Circle(_ context.Context, args *public.CircleArgs) (*public.CircleReply, error) {
	if args.Radius <= 0 {
		return nil, xamqp.NewError(xamqp.InvalidArgument, "radius must be a positive number")
	}

	return &public.CircleReply{
		Length: 2 * math.Pi * args.Radius,
		Square: math.Pi * args.Radius * args.Radius,
	}, nil
}
