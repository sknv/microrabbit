package xcontext

import (
	"context"
	"time"
)

func Timeout(ctx context.Context) (time.Duration, bool) {
	deadline, ok := ctx.Deadline()
	if !ok {
		return 0, false
	}
	return deadline.Sub(time.Now()), true
}
