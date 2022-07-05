package main

import (
	"context"
	"flag"
	"fmt"
	
	"net"
	"net/url"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"strconv"

	"github.com/crossnokaye/facilityconfig/envloader"
	calc "github.com/crossnokaye/carbon/services/calc/gen/calc"
	calcpb "github.com/crossnokaye/carbon/services/calc/gen/grpc/calc/pb"
	calcsvr "github.com/crossnokaye/carbon/services/calc/gen/grpc/calc/server"

	calcapi "github.com/crossnokaye/carbon/services/calc"


	"github.com/crossnokaye/carbon/services/calc/clients/storage"
	"github.com/crossnokaye/carbon/services/calc/clients/power"

	"github.com/go-pg/pg"
	"github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/prometheus/client_golang/prometheus"
	"goa.design/clue/health"
	"goa.design/clue/log"
	"goa.design/clue/metrics"
	"goa.design/clue/trace"
	"goa.design/goa/v3/grpc/middleware"
	"golang.org/x/sync/errgroup"

	"google.golang.org/grpc/credentials/insecure"
	grpcmiddleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpcmdlwr "goa.design/goa/v3/grpc/middleware"
	"goa.design/goa/v3/middleware"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func resolveStringEnv(key, def, desc string) string {
	res := def
	if env := os.Getenv(key); env != "" {
		res = env
	}
	return res
}

func resolveBooleanEnv(key string, def bool, desc string) bool {
	res := def
	if env := os.Getenv(key); env != "" {
		if b, err := strconv.ParseBool(env); err == nil {
			res = b
		}
	}
	return res
}

func initLogging(debug bool) context.Context {
	ctx := log.Context(context.Background(),
		log.WithFormat(log.FormatJSON),
		log.WithFunc(trace.Log))
	ctx = log.With(ctx, log.KV{K: "svc", V: calc.ServiceName})
	if debug {
		ctx = log.Context(ctx, log.WithDebug())
		log.Debugf(ctx, "debug logging enabled")
	}
	return ctx
}

func initTrace(ctx context.Context, grafanaAgentAddr string) context.Context {
	log.Debugf(ctx, "connecting to Grafana agent %s", grafanaAgentAddr)
	conn, err := grpc.DialContext(ctx, grafanaAgentAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock())
	if err != nil {
		log.Errorf(ctx, err, "failed to connect to Grafana agent")
		os.Exit(1)
	}
	log.Debugf(ctx, "connected to Grafana agent %s", grafanaAgentAddr)
	ctx, err = trace.Context(ctx, calc.ServiceName, trace.WithGRPCExporter(conn))
	if err != nil {
		log.Errorf(ctx, err, "failed to initialize tracing")
		os.Exit(1)
	}
	return ctx
}

func initTransport(ctx context.Context, endpoints *calc.Endpoints) *grpc.Server {
	server := calcsvr.New(endpoints, nil)
	grpcServer := grpc.NewServer(
		grpc_middleware.WithUnaryServerChain(
			log.UnaryServerInterceptor(ctx),
			trace.UnaryServerInterceptor(ctx),
			middleware.UnaryRequestID(),
			middleware.UnaryServerLogContext(log.AsGoaMiddlewareLogger),
			metrics.UnaryServerInterceptor(ctx),
		))
	calcpb.RegisterPastValuesServer(grpcServer, server)
	reflection.Register(grpcServer)
	for svc, info := range grpcServer.GetServiceInfo() {
		for _, m := range info.Methods {
			log.Print(ctx, log.KV{K: "method", V: svc + "/" + m.Name})
		}
	}
	return grpcServer
}
/**
func initHealthCheck(ctx context.Context, httpAddr string) *http.Server {
	check := log.HTTP(ctx)(health.Handler(health.NewChecker()))
	http.Handle("/healthz", check)
	http.Handle("/livez", check)
	http.Handle("/metrics", metrics.Handler(ctx))
	return &http.Server{Addr: httpAddr}
}
*/
func main() {
	// Define command line flags, add any other flag required to configure the
	// service.
	var (
		hostF     = flag.String("host", "localhost", "Server host (valid values: localhost)")
		domainF   = flag.String("domain", "", "Host domain name (overrides host domain specified in service design)")
		grpcPortF = flag.String("grpc-port", "", "gRPC port (overrides host gRPC port specified in service design)")
		secureF   = flag.Bool("secure", false, "Use secure scheme (https or grpcs)")
		dbgF      = flag.Bool("debug", false, "Log request and response bodies")
		env = resolveStringEnv("ENV", "dev", "facility environment")

		chaddr = resolveStringEnv("CH_ADDRESS", "127.0.0.1:8087", "click house server url")
		chuser = resolveStringEnv("CH_USER", "atlas", "click house user name")
		chpwd  = resolveStringEnv("CH_PASSWORD", "atlas", "click house password")
		chsec  = resolveBooleanEnv("CH_SECURE", false, "click house secure connection")

		debug            = resolveBooleanEnv("DEBUG", false, "enable debug logging")
		grafanaAgentAddr = resolveStringEnv("GRAFANA_AGENT_ADDR", ":4317", "grafana agent address")

	


	)
	flag.Parse()

	// Setup logger. Replace logger with your own log package of choice.
	ctx := initLogging(debug)

	ctx = initTrace(ctx, grafanaAgentAddr)
	ctx = metrics.Context(ctx, calc.ServiceName)
	
	log.Printf(ctx, "starting %s service", calc.ServiceName)
	
	log.Print(ctx, log.KV{K: "creating", V: "report repositories"})

	//initialize the clients
	
	// Initialize the services.
	var (
		calcSvc calc.Service
	)
	{
		calcSvc = calcapi.NewCalc(logger)
	}

	// Wrap the services in endpoints that can be invoked from other services
	// potentially running in different processes.
	var (
		calcEndpoints *calc.Endpoints
	)
	{
		calcEndpoints = calc.NewEndpoints(calcSvc)
	}

	// Create channel used by both the signal handler and server goroutines
	// to notify the main goroutine when to stop the server.
	errc := make(chan error)

	// Setup interrupt handler. This optional step configures the process so
	// that SIGINT and SIGTERM signals cause the services to stop gracefully.
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errc <- fmt.Errorf("%s", <-c)
	}()

	var wg sync.WaitGroup
	ctx, cancel := context.WithCancel(context.Background())

	// Start the servers and send errors (if any) to the error channel.
	switch *hostF {
	case "localhost":

		{
			addr := "grpc://localhost:8080"
			u, err := url.Parse(addr)
			if err != nil {
				fmt.Fprintf(os.Stderr, "invalid URL %#v: %s\n", addr, err)
				os.Exit(1)
			}
			if *secureF {
				u.Scheme = "grpcs"
			}
			if *domainF != "" {
				u.Host = *domainF
			}
			if *grpcPortF != "" {
				h, _, err := net.SplitHostPort(u.Host)
				if err != nil {
					fmt.Fprintf(os.Stderr, "invalid URL %#v: %s\n", u.Host, err)
					os.Exit(1)
				}
				u.Host = net.JoinHostPort(h, *grpcPortF)
			} else if u.Port() == "" {
				u.Host = net.JoinHostPort(u.Host, "8080")
			}
			handleGRPCServer(ctx, u, calcEndpoints, &wg, errc, logger, *dbgF)
		}

	default:
		fmt.Fprintf(os.Stderr, "invalid host argument: %q (valid hosts: localhost)\n", *hostF)
	}

	// Wait for signal.
	logger.Printf("exiting (%v)", <-errc)

	// Send cancellation signal to the goroutines.
	cancel()

	wg.Wait()
	logger.Println("exited")
}
