package main

import (
	"flag"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"os"
)

func main() {
	fs := flag.NewFlagSet("", flag.ExitOnError)
	var (
		grpcAddr = fs.String("grpc.addr", "localhost:8002", "gRPC address")
		a        = fs.Int64("a", 1, "a value")
		b        = fs.Int64("b", 2, "b value")
	)

	fs.Parse(os.Args[1:])

	logger := log.NewLogfmtLogger(os.Stdout)
	logCtx := log.NewContext(logger).With("ts", log.DefaultTimestampUTC)

	var e endpoint.Endpoint

	cc, err := grpc.Dial(*grpcAddr)
	if err != nil {
		logCtx.Log("grpc", "connect", "err", err)
		os.Exit(1)
	}
	e = NewGRPCClient(cc)
	response, err := e(context.Background(), AddRequest{A: *a, B: *b})
	if err != nil {
		logCtx.Log("response", response, "err", err)
		os.Exit(1)
	}

	addResponse, ok := response.(AddResponse)
	if !ok {
		logCtx.Log("response", response, "ok", ok)
		os.Exit(1)
	}

	logCtx.Log("response", addResponse.V)
}
