package main

import (
	"github.com/go-kit/kit/endpoint"
	"golang.org/x/net/context"
	"log"
)

func Chain(middewares ...endpoint.Middleware) func(endpoint.Endpoint) endpoint.Endpoint {
	return func(e endpoint.Endpoint) endpoint.Endpoint {
		for i := len(middewares) - 1; i >= 0; i-- {
			e = middewares[i](e)

		}
		return e
	}
}

func main() {
	e := func(_ context.Context, request interface{}) (interface{}, error) {
		log.Println("e")
		return nil, nil
	}
	mw1 := func(next endpoint.Endpoint) endpoint.Endpoint {
		log.Println("B:mw1")
		return func(ctx context.Context, request interface{}) (interface{}, error) {
			log.Println("mw1")
			return next(ctx, request)
		}
	}
	mw2 := func(next endpoint.Endpoint) endpoint.Endpoint {
		log.Println("B:mw2")
		return func(ctx context.Context, request interface{}) (interface{}, error) {
			log.Println("mw2")
			return next(ctx, request)
		}
	}

	e = Chain(mw1, mw2)(e)
	ctx := context.Background()
	log.Println("aaa")
	log.Println(e(ctx, nil))

}
