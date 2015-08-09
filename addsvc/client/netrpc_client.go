package main

import (
	"fmt"
	"github.com/go-kit/kit/endpoint"
	"golang.org/x/net/context"
	"net/rpc"
)

func NewNetRpcClient(c *rpc.Client) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		var (
			errs      = make(chan error, 1)
			responses = make(chan interface{}, 1)
		)

		go func() {
			var response AddResponse

			if err := c.Call("addsvc.Add", request, &response); err != nil {
				errs <- err
				return
			}
			fmt.Println("response", response)
			responses <- response

		}()
		select {
		case <-ctx.Done():
			return nil, context.DeadlineExceeded
		case err := <-errs:
			return nil, err
		case resp := <-responses:
			return resp, nil
		}

	}

}
