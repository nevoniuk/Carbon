package main

import (
	//import all packages
	"context"
	//"crypto/tls"
	"flag"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	//import grpc stuff
	//ch "github.com/ClickHouse/clickhouse-go/v2"
	grpcmiddleware "github.com/grpc-ecosystem/go-grpc-middleware"
	//"goa.design/clue/health"
	"goa.design/clue/log"
	//"goa.design/clue/metrics"
	"goa.design/clue/trace"
	goagrpcmiddleware "goa.design/goa/v3/grpc/middleware"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
	
	//"github.com/crossnokaye/rates/clients/clickhouse"
	"github.com/crossnokaye/carbon/poller"
	"github.com/crossnokaye/carbon/poller/clients/storage"
	"github.com/crossnokaye/carbon/poller/clients/carbonara"
	genpb "github.com/crossnokaye/carbon/poller/gen/grpc/data/pb"
	gengrpc "github.com/crossnokaye/carbon/poller/gen/grpc/data/server"
	genpoller "github.com/crossnokaye/carbon/poller/gen/data"
)

func main() {
	var (
		grpcaddr = flag.String("grpc-addr", ":12500", "gRPC listen address")
		httpaddr = flag.String("http-addr", ":12501", "HTTP listen address (health checks and metrics)")
		//agentaddr = flag.String("agent-addr", ":4317", "Grafana agent listen address")

		//chaddr = flag.String("ch-addr", os.Getenv("CLICKHOUSE_ADDR"), "ClickHouse host address")
		//chuser = flag.String("ch-user", os.Getenv("CLICKHOUSE_USER"), "ClickHouse user")
		//chpwd  = flag.String("ch-pwd", os.Getenv("CLICKHOUSE_PASSWORD"), "ClickHouse password")
		//chssl  = flag.Bool("ch-ssl", os.Getenv("CLICKHOUSE_SSL") != "", "ClickHouse connection SSL")

		debug = flag.Bool("debug", false, "Enable debug logs")
		//test  = flag.Bool("test", os.Getenv("TEST_ENV") != "", "Enable test mode")
	)
	flag.Parse()

	//sets up JSON logger
	format := log.FormatJSON
	if log.IsTerminal() {
		format = log.FormatTerminal
	}
	ctx := log.Context(context.Background(), log.WithFormat(format), log.WithFunc(trace.Log))
	ctx = log.With(ctx, log.KV{K: "svc", V: gencarbon.ServiceName})
	if *debug {
		ctx = log.Context(ctx, log.WithDebug())
		log.Debugf(ctx, "debug logs enabled")
	}
	log.Printf(ctx, "starting")
	//set up tracing
	//set up metrics
	//ctx = metrics.Context(ctx, gencarbon.ServiceName)
	//set up clients
	c := &http.Client{Transport: trace.Client(ctx, http.DefaultTransport)}
	csc := carbonara.New(c)
	//create db client
	//create service and endpoints
	svc := smartservice.New(csc)
	endpoints := gencarbon.NewEndpoints(svc)
	//create transport

	server := gengrpc.New(endpoints, nil)
	grpcsvr := grpc.NewServer(
		grpcmiddleware.WithUnaryServerChain(
			log.UnaryServerInterceptor(ctx),
			trace.UnaryServerInterceptor(ctx),
			goagrpcmiddleware.UnaryRequestID(),
			goagrpcmiddleware.UnaryServerLogContext(log.AsGoaMiddlewareLogger),
			//metrics.UnaryServerInterceptor(ctx),
		))
	genpb.RegisterCarbonServer(grpcsvr, server)
	reflection.Register(grpcsvr)
	for svc, info := range grpcsvr.GetServiceInfo() {
		for _, m := range info.Methods {
			log.Print(ctx, log.KV{K: "method", V: svc + "/" + m.Name})
		}
	}
	//start health check
	//make gRPC and HTTP servers
	httpsvr := &http.Server{Addr: *httpaddr}
	errc := make(chan error)
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errc <- fmt.Errorf("%s", <-c)
	}()
	ctx, cancel := context.WithCancel(ctx)

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()

		go func() {
			log.Printf(ctx, "HTTP server listening on %s", httpsvr.Addr)
			errc <- httpsvr.ListenAndServe()
		}()

		var l net.Listener
		go func() {
			var err error
			l, err = net.Listen("tcp", *grpcaddr)
			if err != nil {
				errc <- err
				return
			}
			log.Printf(ctx, "gRPC server listening on %s", l.Addr())
			errc <- grpcsvr.Serve(l)
		}()

		<-ctx.Done()
		log.Printf(ctx, "shutting down HTTP and gRPC servers")

		// Shutdown gracefully with a 30s timeout.
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		grpcsvr.GracefulStop()
		httpsvr.Shutdown(ctx)
	}()

	// Cleanup
	err := <-errc
	if err != nil && err.Error() != "interrupt" {
		log.Errorf(ctx, err, "exiting")
	} else {
		log.Printf(ctx, "exiting")
	}
	cancel()
	wg.Wait()
	log.Printf(ctx, "exited")
}
