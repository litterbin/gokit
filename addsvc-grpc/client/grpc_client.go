package main

import (
	"github.com/go-kit/kit/endpoint"
	"github.com/litterbin/gokit/addsvc-grpc/pb"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

func NewGRPCClient(cc *grpc.ClientConn) endpoint.Endpoint {
	client := pb.NewAddClient(cc)

	return func(ctx context.Context, request interface{}) (interface{}, error) {
		var (
			errs      = make(chan error, 1)
			responses = make(chan interface{}, 1)
		)
		go func() {
			addReq, ok := request.(AddRequest)
			if !ok {
				errs <- endpoint.ErrBadCast
				return
			}

			reply, err := client.Add(ctx, &pb.AddRequest{A: addReq.A, B: addReq.B})
			if err != nil {
				errs <- err
				return
			}
			responses <- AddResponse{V: reply.V}
		}()
		select {
		case <-ctx.Done():
			return nil, context.DeadlineExceeded
		case err := <-errs:
			return nil, err
		case response := <-responses:
			return response, nil
		}
	}
}
