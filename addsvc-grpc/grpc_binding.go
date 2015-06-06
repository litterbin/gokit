package main

import (
	"time"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/metrics"
	"github.com/litterbin/gokit/addsvc-grpc/pb"
	"golang.org/x/net/context"
)

// A binding wraps an Endpoint so that it's usable by a transport.
// grpcBinding makes an Endpoint usable over gRPC.
type grpcBinding struct{ endpoint.Endpoint }

func (b grpcBinding) Add(ctx context.Context, req *pb.AddRequest) (*pb.AddReply, error) {
	addReq := addRequest{req.A, req.B}
	r, err := b.Endpoint(ctx, addReq)
	if err != nil {
		return nil, err
	}

	resp, ok := r.(*addResponse)
	if !ok {
		return nil, endpoint.ErrBadCast
	}

	return &pb.AddReply{
		V: resp.V,
	}, nil
}
