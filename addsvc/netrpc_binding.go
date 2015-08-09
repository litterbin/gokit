package main

import (
	"fmt"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
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
		logger      = ctx.Value("logger").(log.Logger)
	)
	defer cancel()

	logger.Log("binding", "netrpc")

	go func() {
		resp, err := b.Endpoint(ctx, &request)
		if err != nil {
			logger.Log("endpoint", "add", "err", err)
			errs <- err
			return
		}

		addResp, ok := resp.(*AddResponse)
		if !ok {
			errs <- endpoint.ErrBadCast
			return
		}
		logger.Log("endpoint", "add", "resp", addResp.V)
		responses <- *addResp
	}()

	select {
	case <-ctx.Done():
		return context.DeadlineExceeded
	case err := <-errs:
		return err
	case resp := <-responses:
		(*response) = resp
		logger.Log("response", fmt.Sprintf("%v", response))
		return nil
	}
}
