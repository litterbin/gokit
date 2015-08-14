package endpoints

import (
	"github.com/go-kit/kit/endpoint"
	"golang.org/x/net/context"
	"testing"
)

func AddEndpoint(_ context.Context, _ interface{}) (interface{}, error) {
	return nil, nil
}

func Middleware(next endpoint.Endpoint) endpoint.Endpoint {
	return func(ctx context.Context, req interface{}) (interface{}, error) {

		_, err := next(ctx, req)

		return 1, err
	}
}

func TestEndpoints(t *testing.T) {
	Add("add", AddEndpoint) // endpoints.Add
	Use(Middleware)

	resp, err := Endpoints["add"](context.Background(), nil)
	if err != nil {
		t.Error(err)
		return
	}
	ret, ok := resp.(int)
	if !ok {
		t.Error(err)
		return
	}
	if ret != 1 {
		t.Error(ret)
	}
}
