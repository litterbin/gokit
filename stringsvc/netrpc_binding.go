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

func (n NetRpcBinding) Uppercase(req UppercaseRequest, res *UppercaseResponse) error {
	ctx, cancel := context.WithCancel(n.Context)
	defer cancel()
	responses := make(chan UppercaseResponse, 1)
	errs := make(chan error, 1)

	go func() {
		resp, err := n.uppercaseEndpoint(ctx, req)
		if err != nil {
			errs <- err
			return
		}
		responses <- resp.(UppercaseResponse)
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

func (n *NetRpcBinding) Count(req CountRequest, res *CountResponse) error {
	ctx, cancel := context.WithCancel(n.Context)
	defer cancel()
	responses := make(chan CountResponse, 1)
	errs := make(chan error, 1)

	go func() {
		resp, err := n.countEndpoint(ctx, req)
		if err != nil {
			errs <- err
			return
		}
		responses <- resp.(CountResponse)
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
