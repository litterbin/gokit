package main

import (
	"github.com/go-kit/kit/endpoint"
	"golang.org/x/net/context"
)

type NetRpcBinding struct {
	ctx context.Context
	endpoint.Endpoint
}

func (b NetRpcBinding) Add(request AddRequest, response *AddResponse) error {
	var (
		ctx, cancel = context.WithCancel(b.ctx)
		errs        = make(chan error, 1)
		responses   = make(chan AddResponse, 1)
	)
	defer cancel()

	go func() {
		resp, err := b.Endpoint(ctx, request)
		if err != nil {
			errs <- err
			return
		}

		addResp, ok := resp.(AddResponse)
		if !ok {
			errs <- endpoint.ErrBadCast
			return
		}
		responses <- addResp
	}()

	select {
	case <-ctx.Done():
		return context.DeadlineExceeded
	case err := <-errs:
		return err
	case resp := <-responses:
		(*response) = resp
		return nil
	}
}
