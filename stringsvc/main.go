package main

import (
	"golang.org/x/net/context"

	"log"
	"net/http"
	"net/rpc"
)

func main() {
	ctx := context.Background()
	svc := stringService{}

	//net/rpc

	netRpcB := NetRpcBinding{ctx, makeUppercaseEndpoint(svc), makeCountEndpoint(svc)}

	s := rpc.NewServer()
	s.RegisterName("stringsvc", netRpcB)
	s.HandleHTTP(rpc.DefaultRPCPath, rpc.DefaultDebugPath)
	err := http.ListenAndServe(":8080", s)
	if err != nil {
		log.Fatal(err)
	}

}
