package main

import (
	"github.com/go-kit/kit/endpoint"
	"golang.org/x/net/context"
)

type UppercaseRequest struct {
	S string `json:"s"`
}

type UppercaseResponse struct {
	V   string `json:"v"`
	Err error  `json:"err"`
}

type CountRequest struct {
	S string `json:"s"`
}

type CountResponse struct {
	V int `json:"v"`
}

func makeUppercaseEndpoint(svc StringService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(UppercaseRequest)
		v, err := svc.Uppercase(req.S)
		return UppercaseResponse{v, err}, nil
	}
}

func makeCountEndpoint(svc StringService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(CountRequest)
		v := svc.Count(req.S)
		return CountResponse{v}, nil
	}
}
