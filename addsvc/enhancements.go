package main

import (
	"github.com/go-kit/kit/log"
	"golang.org/x/net/context"
	"time"
)

func logging(logger log.Logger) func(Add) Add {
	return func(next Add) Add {
		return func(ctx context.Context, a, b int64) (v int64) {
			defer func(begin time.Time) {
				logger.Log("a", a, "b", b, "result", v, "took", time.Since(begin))
			}(time.Now())
			v = next(ctx, a, b)
			return
		}
	}
}
