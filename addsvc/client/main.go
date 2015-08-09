package main

import (
	"flag"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"net/rpc"
	"os"
)

func main() {
	fs := flag.NewFlagSet("", flag.ExitOnError)
	var (
		grpcAddr   = fs.String("grpc.addr", "localhost:8002", "gRPC address")
		netRpcAddr = fs.String("netrpc.addr", "localhost:8003", "netrpc address")
		transport  = fs.String("transport", "grpc", "grpc netrpc")
		a          = fs.Int64("a", 1, "a value")
		b          = fs.Int64("b", 2, "b value")
	)

	fs.Parse(os.Args[1:])

	logger := log.NewLogfmtLogger(os.Stdout)
	logger = log.NewContext(logger).With("ts", log.DefaultTimestampUTC)

	var e endpoint.Endpoint

	switch *transport {
	case "grpc":
		cc, err := grpc.Dial(*grpcAddr)
		if err != nil {
			logger.Log("grpc", "connect", "err", err)
			os.Exit(1)
		}
		e = NewGRPCClient(cc)
	case "netrpc":
		client, err := rpc.DialHTTP("tcp", *netRpcAddr)
		if err != nil {
			logger.Log("netrpc", "connect", "err", err)
			os.Exit(1)
		}
		e = NewNetRpcClient(client)
	}

	response, err := e(context.Background(), AddRequest{A: *a, B: *b})
	if err != nil {
		logger.Log("response", response, "err", err)
		os.Exit(1)
	}

	addResponse, ok := response.(AddResponse)
	if !ok {
		logger.Log("response", response, "ok", ok)
		os.Exit(1)
	}

	logger.Log("response", addResponse.V)
}
