package main

import (
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"golang.org/x/net/context"
)

// Add is the abstract definition of what this service does.
type Add func(context.Context, int64, int64) int64

// pureAdd implements Add with no dependencies.
func pureAdd(_ context.Context, a, b int64) int64 { return a + b }

// proxyAdd returns an implementation of Add that invokes a remote Add
// service.
func proxyAdd(e endpoint.Endpoint, logger log.Logger) Add {
	return func(ctx context.Context, a, b int64) int64 {
		resp, err := e(ctx, &addRequest{a, b})
		if err != nil {
			logger.Log("err", err)
			return 0
		}
		addResp, ok := resp.(*addResponse)
		if !ok {
			logger.Log("err", endpoint.ErrBadCast)
			return 0
		}
		return addResp.V
	}
}
