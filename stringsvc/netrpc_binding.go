package main

import (
	"github.com/go-kit/kit/endpoint"
	"golang.org/x/net/context"
)

type NetRpcBinding struct {
	Context           context.Context
	uppercaseEndpoint endpoint.Endpoint
	countEndpoint     endpoint.Endpoint
}

func (n *NetRpcBinding) Uppercase(req uppercaseRequest, res *uppercaseResponse) error {
	ctx, cancel := context.WithCancel(n.Context)
	defer cancel()
	responses := make(chan uppercaseResponse, 1)
	errs := make(chan error, 1)

	go func() {
		resp, err := n.uppercaseEndpoint(ctx, req)
		if err != nil {
			errs <- err
			return
		}
		responses <- resp.(uppercaseResponse)
		return
	}()

	select {
	case <-ctx.Done():
		return context.DeadlineExceeded
	case err := <-errs:
		return err
	case resp := <-responses:
		(*res) = resp
	}

	return nil
}

func (n *NetRpcBinding) Count(req countRequest, res *countResponse) error {
	ctx, cancel := context.WithCancel(n.Context)
	defer cancel()
	responses := make(chan countResponse, 1)
	errs := make(chan error, 1)

	go func() {
		resp, err := n.countEndpoint(ctx, req)
		if err != nil {
			errs <- err
			return
		}
		responses <- resp.(countResponse)
		return
	}()

	select {
	case <-ctx.Done():
		return context.DeadlineExceeded
	case err := <-errs:
		return err
	case resp := <-responses:
		(*res) = resp
	}

	return nil
}
