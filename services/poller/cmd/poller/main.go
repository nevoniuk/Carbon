package main

import (
	"context"
	"crypto/tls"
	"flag"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
	ch "github.com/ClickHouse/clickhouse-go/v2"
	grpcmiddleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"goa.design/clue/health"
	"goa.design/clue/log" //the log package provides a context based logging API
	"goa.design/clue/metrics" //metrics is for exposing a promethes compatible metrics HTTP endpoint
	"goa.design/clue/trace"
	goagrpcmiddleware "goa.design/goa/v3/grpc/middleware"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
	"github.com/crossnokaye/carbon/clients/clickhouse"
	pollerapi "github.com/crossnokaye/carbon/services/poller"
	"github.com/crossnokaye/carbon/services/poller/clients/carbonara"
	"github.com/crossnokaye/carbon/services/poller/clients/storage"
	genpoller "github.com/crossnokaye/carbon/services/poller/gen/poller"
	genpb "github.com/crossnokaye/carbon/services/poller/gen/grpc/poller/pb"
	gengrpc "github.com/crossnokaye/carbon/services/poller/gen/grpc/poller/server"
)

func main() {
	var (
		grpcaddr  = flag.String("grpc-addr", ":12500", "gRPC listen address")
		httpaddr  = flag.String("http-addr", ":12501", "HTTP listen address")
		agentaddr = flag.String("agent-addr", ":4317", "Grafana agent listen address")
		chaddr = flag.String("ch-addr", os.Getenv("CLICKHOUSE_ADDR"), "ClickHouse host address")
		chuser = flag.String("ch-user", os.Getenv("CLICKHOUSE_USER"), "ClickHouse user")
		chpwd  = flag.String("ch-pwd", os.Getenv("CLICKHOUSE_PASSWORD"), "ClickHouse password")
		chssl  = flag.Bool("ch-ssl", os.Getenv("CLICKHOUSE_SSL") != "", "ClickHouse connection SSL")
		monitoringEnabled = flag.Bool("monitoring-enabled", true, "monitoring")
		debug = flag.Bool("debug", false, "Enable debug logs")
		carbonKey = flag.String("singularity-key", os.Getenv("SINGULARITY_API_KEY"), "The API key for Singularity")
	)
	
	flag.Parse()
	format := log.FormatJSON
	if log.IsTerminal() {
		format = log.FormatTerminal
	}
	//log context
	ctx := log.With(log.Context(context.Background(), log.WithFormat(format)), log.KV{K: "svc", V: genpoller.ServiceName})
	
	//log clickhouse credentials and monitoring status
	log.Info(ctx,
		log.KV{K: "singularity-key", V: *carbonKey},
        log.KV{K: "ch-addr", V: *chaddr},
        log.KV{K: "ch-user", V: *chuser})
		
		log.Info(ctx,
			log.KV{K: "monitoringEnabled", V: *monitoringEnabled})

	if *debug {
		ctx = log.Context(ctx, log.WithDebug())
		log.Debugf(ctx, "debug logs enabled")
	}
	//monitoring enabled is true - initialize tracing - only in production
	//monitoring enabled is false - in janeway - dont initialize tracing. set in .env
	// Setup tracing
	if *monitoringEnabled {
		conn, err := grpc.DialContext(ctx, *agentaddr,
			grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			log.Errorf(ctx, err, "failed to connect to Grafana agent")
			os.Exit(1)
		}
		if ctx, err = trace.Context(ctx, genpoller.ServiceName, trace.WithGRPCExporter(conn)); err != nil {
			log.Errorf(ctx, err, "failed to initialize tracing")
			os.Exit(1)
		}
	}
	

	//initialize the metrics
	ctx = metrics.Context(ctx, genpoller.ServiceName)

	//initialize the clients
	c := &http.Client{}
	if *monitoringEnabled {
		c.Transport = trace.Client(ctx, http.DefaultTransport)
	}
	csc := carbonara.New(c, *carbonKey)
	chadd := *chaddr
	if chadd == "" {
		chadd = "localhost:8088" //dev default
	}
	var tlsConfig *tls.Config
	if *chssl {
		tlsConfig = &tls.Config{InsecureSkipVerify: true}
	}
	options := &ch.Options{
		TLS:  tlsConfig,
		Addr: []string{chadd},
		Auth: ch.Auth{
			Username: *chuser,
			Password: *chpwd},
	}
	chcon, err := ch.Open(options)

	retries := 0
	for err != nil && retries < 10 {
		time.Sleep(time.Second)
		chcon, err = ch.Open(options)
		retries++
	}
	con := clickhouse.New(chcon)
	if err != nil {
		log.Errorf(ctx, err, "could not connect to clickhouse")
	}
	defer con.Close()
	dbc := storage.New(con)
	if err := dbc.Init(ctx, true); err != nil {
		log.Errorf(ctx, err, "could not initialize clickhouse")
	}

	//setup the service
	pollerSvc := pollerapi.NewPoller(ctx, csc, dbc)
	endpoints := genpoller.NewEndpoints(pollerSvc)
	//create transport
	server := gengrpc.New(endpoints, nil)
	var grpcsvr *grpc.Server
	if *monitoringEnabled {
		grpcsvr = grpc.NewServer(
			grpcmiddleware.WithUnaryServerChain(
				log.UnaryServerInterceptor(ctx),
				trace.UnaryServerInterceptor(ctx),
				goagrpcmiddleware.UnaryRequestID(),
				goagrpcmiddleware.UnaryServerLogContext(log.AsGoaMiddlewareLogger),
			),
		)
	} else {
		grpcsvr = grpc.NewServer(
			grpcmiddleware.WithUnaryServerChain(
				log.UnaryServerInterceptor(ctx),
				goagrpcmiddleware.UnaryRequestID(),
				goagrpcmiddleware.UnaryServerLogContext(log.AsGoaMiddlewareLogger),
			),
		)
	}

	genpb.RegisterPollerServer(grpcsvr, server)
	reflection.Register(grpcsvr)
	for svc, info := range grpcsvr.GetServiceInfo() {
		for _, m := range info.Methods {
			log.Print(ctx, log.KV{K: "method", V: svc + "/" + m.Name})
		}
	}

	check := health.Handler(health.NewChecker(dbc))
	http.Handle("/healthz", check)
	http.Handle("/livez", check)
	http.Handle("/metrics", metrics.Handler(ctx))
	httpsvr := &http.Server{Addr: *httpaddr}


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
	err = <-errc
	if err != nil && err.Error() != "interrupt" {
		log.Errorf(ctx, err, "exiting")
	} else {
		log.Printf(ctx, "exiting")
	}
	cancel()
	wg.Wait()
	log.Printf(ctx, "exited")
}
