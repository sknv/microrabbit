package xcontext

import (
	"context"
	"time"
)

func Timeout(ctx context.Context) (time.Duration, bool) {
	deadline, isset := ctx.Deadline()
	if !isset {
		return 0, false
	}
	return deadline.Sub(time.Now()), true
}
