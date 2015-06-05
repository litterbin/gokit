package main

import (
	"time"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/metrics"
	"github.com/litterbin/gokit/addsvc-grpc/pb"
	"golang.org/x/net/context"
)

type grpcBinding struct{ endpoint.Endpoint }

func (b grpcBinding) Add(ctx context.Context, req *pb.AddRequest) (*pb.AddReply, error) {
	addReq := addRequest{req.A, req.B}
	//    r,err :=
}
