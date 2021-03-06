package main

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"

	"net"
	"net/http"
	"net/rpc"

	"golang.org/x/net/context"
	"google.golang.org/grpc"

	"github.com/go-kit/kit/endpoint"
	//	"github.com/go-kit/kit/metrics"
	//	"github.com/go-kit/kit/metrics/expvar"
	//httptransport "github.com/go-kit/kit/transport/http"

	kitlog "github.com/go-kit/kit/log"
	stdlog "log"

	"github.com/litterbin/gokit/addsvc/pb"
)

func main() {
	fs := flag.NewFlagSet("", flag.ExitOnError)

	var (
		debugAddr = fs.String("debug.addr", ":8000", "Address for HTTP debug/instrumentation server")
		//httpAddr  = fs.String("http.addr", ":8001", "Address for HTTP (JSON) server")
		grpcAddr   = fs.String("grpc.addr", ":8002", "Address for gRPC server")
		netrpcAddr = fs.String("netrpc.addr", ":8003", "Address for net/rpc server")

		//zipkinServiceName            = fs.String("zipkin.service.name", "addsvc", "Zipkin service name")
		//zipkinCollectorAddr          = fs.String("zipkin.collector.addr", "", "Zipkin Scribe collector address (empty will log spans)")
		//zipkinCollectorTimeout       = fs.Duration("zipkin.collector.timeout", time.Second, "Zipkin collector timeout")
		//zipkinCollectorBatchSize     = fs.Int("zipkin.collector.batch.size", 100, "Zipkin collector batch size")
		//zipkinCollectorBatchInterval = fs.Duration("zipkin.collector.batch.interval", time.Second, "Zipkin collector batch interval")
	)

	flag.Usage = fs.Usage // only show our flags
	fs.Parse(os.Args[1:])

	// `package log` domain
	var logger kitlog.Logger
	logger = kitlog.NewLogfmtLogger(os.Stderr)
	logger = kitlog.NewContext(logger).With("ts", kitlog.DefaultTimestampUTC)
	stdlog.SetOutput(kitlog.NewStdlibAdapter(logger)) // redirect stdlib logging to us
	stdlog.SetFlags(0)                                // flags are handled in our logger

	var a Add = pureAdd
	a = logging(logger)(a)

	var e endpoint.Endpoint
	e = makeEndpoint(a)

	// Mechanical stuff
	rand.Seed(time.Now().UnixNano())

	bg := context.Background()
	root := context.WithValue(bg, "logger", logger)
	errc := make(chan error)

	go func() {
		errc <- interrupt()
	}()

	//Transport: HTTP (debug/intrumentation)
	go func() {
		logger.Log("addr", *debugAddr, "transport", "debug")
		errc <- http.ListenAndServe(*debugAddr, nil)
	}()

	//Transport: HTTP (JSON)
	/*
		go func() {
			ctx, cancel := context.WithCancel(root)
			defer cancel()

				after := httptransport.BindingAfter(httptransport.SetContentType("application/json"))

				makeRequest := func() interface{} { return &addRequest{} }

				var handler http.Handler
				handler = httptransport.NewBinding(ctx, makeRequest, jsoncodec.New(), e, nil, after)

				mux := http.NewServeMux()
				mux.Handle("/add", handler)
				logger.Log("addr", *httpAddr, "transport", "HTTP")
				errc <- http.ListenAdnServe(*httpAddr, mux)
		}()
	*/

	//Transport: gRPC
	go func() {
		ln, err := net.Listen("tcp", *grpcAddr)
		if err != nil {
			errc <- err
			return
		}

		s := grpc.NewServer()

		var addServer pb.AddServer
		addServer = grpcBinding{e}

		pb.RegisterAddServer(s, addServer)
		logger.Log("addr", *grpcAddr, "transport", "gRPC")
		errc <- s.Serve(ln)

	}()

	//Transport: net/rpc
	go func() {
		ctx, cancel := context.WithCancel(root)
		defer cancel()
		s := rpc.NewServer()
		s.RegisterName("addsvc", NetRpcBinding{ctx, e})
		s.HandleHTTP(rpc.DefaultRPCPath, rpc.DefaultDebugPath)
		logger.Log("addr", *netrpcAddr, "transport", "net/rpc")
		errc <- http.ListenAndServe(*netrpcAddr, s)
	}()

	logger.Log("fatal", <-errc)
}

func interrupt() error {
	c := make(chan os.Signal)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	return fmt.Errorf("%s", <-c)
}
