package endpoints

import (
	"github.com/go-kit/kit/endpoint"
)

var Endpoints = map[string]endpoint.Endpoint{}

func Add(name string, ep endpoint.Endpoint) {
	Endpoints[name] = ep
}

func Use(middleware endpoint.Middleware) {
	for name, endpoint := range Endpoints {
		Endpoints[name] = middleware(endpoint)
	}
}
