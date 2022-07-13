package main

import (
    "context"
    "flag"
    "fmt"
    "crypto/tls"
    "os"
    "os/signal"
    "sync"
    "syscall"
    "time"
    "net"
	"net/http"


    ch "github.com/ClickHouse/clickhouse-go/v2"
    "github.com/crossnokaye/facilityconfig/envloader"
    gencalc "github.com/crossnokaye/carbon/services/calc/gen/calc"
    calcpb "github.com/crossnokaye/carbon/services/calc/gen/grpc/calc/pb"
    calcsvr "github.com/crossnokaye/carbon/services/calc/gen/grpc/calc/server"

    calcapi "github.com/crossnokaye/carbon/services/calc"
    pastvalsvc "github.com/crossnokaye/carbon/services/calc/clients/power"

    //clients
    "github.com/crossnokaye/carbon/clients/clickhouse"
    "github.com/crossnokaye/carbon/services/calc/clients/storage"
	"github.com/crossnokaye/carbon/services/calc/clients/power_server"

    goagrpcmiddleware "goa.design/goa/v3/grpc/middleware"
    grpcmiddleware "github.com/grpc-ecosystem/go-grpc-middleware"
    "google.golang.org/grpc/credentials/insecure"
    "google.golang.org/grpc"
    "google.golang.org/grpc/reflection"
    "goa.design/clue/health"
    "goa.design/clue/log"
    "goa.design/clue/metrics"
    "goa.design/clue/trace"
 
    
)



func main() {
	// Define command line flags, add any other flag required to configure the
	// service.

	var (
		
		grpcaddr  = flag.String("grpc-addr", "0.0.0.0:12200", "gRPC listen address")
		httpaddr  = flag.String("http-addr", "0.0.0.0:12201", "HTTP listen address")
		agentaddr = flag.String("agent-addr", ":4317", "Grafana agent listen address")

		chaddr = flag.String("ch-addr", os.Getenv("CLICKHOUSE_ADDR"), "ClickHouse host address")
		chuser = flag.String("ch-user", os.Getenv("CLICKHOUSE_USER"), "ClickHouse user")
		chpwd  = flag.String("ch-pwd", os.Getenv("CLICKHOUSE_PASSWORD"), "ClickHouse password")
		chssl  = flag.Bool("ch-ssl", os.Getenv("CLICKHOUSE_SSL") != "", "ClickHouse connection SSL")

		debug = flag.Bool("debug", false, "debug mode")
        //past val service host address
        pastValaddr = flag.String("pastval-add", "localhost:10110", "Past-Value hos address")
        env = flag.String("dev", os.Getenv("ENV"), "facility environment")
		
       
	)
	flag.Parse()

	// Enable logging
	format := log.FormatJSON
	if log.IsTerminal() {
		format = log.FormatTerminal
	}
	ctx := log.With(log.Context(context.Background(), log.WithFormat(format)), log.KV{K: "svc", V: gencalc.ServiceName})
	if *debug {
		ctx = log.Context(ctx, log.WithDebug())
		log.Debugf(ctx, "debug logs enabled")
	}

	// Setup tracing
	grafanaconn, err := grpc.DialContext(ctx, *agentaddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Errorf(ctx, err, "failed to connect to Grafana agent")
		os.Exit(1)
	}
    log.Debugf(ctx, "connected to Grafana agent %s", agentaddr)
	ctx, err = trace.Context(ctx, gencalc.ServiceName, trace.WithGRPCExporter(grafanaconn))
	if err != nil {
		log.Errorf(ctx, err, "failed to initialize tracing")
		os.Exit(1)
	}


    //1.initialize clickhouse client
    chadd := *chaddr
    if chadd == "" {
        chadd = "localhost:8088" // dev default
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
        // CH can take a few seconds before it accepts connection on start
        time.Sleep(time.Second)
        chcon, err = ch.Open(options)
        retries++
    }
    con := clickhouse.New(chcon)
    if err != nil {
        log.Errorf(ctx, err, "could not connect to clickhouse")
    }
    defer con.Close()

	//2.initialize storage client
	dbc := storage.New(con)
	if err := dbc.Init(ctx, true); err != nil {
		log.Errorf(ctx, err, "could not initialize clickhouse")
        return
	}

    //3.initialize power_server repo with the env loader
    powerRepo := power_server.New(envloader.MustNew(*env))
    
	//4.initialize power client with past val grpc connection
    pastValConn, err := grpc.Dial(*pastValaddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		fmt.Errorf("could not connect to Past Values service: %v", err)
        return
	}
    pwc := pastvalsvc.New(pastValConn)

	// Initialize the services and endpoints
	log.Print(ctx, log.KV{K: "creating", V: "service endpoints"})

	var calcSvc gencalc.Service
	calcSvc = calcapi.NewCalc(ctx, pwc, dbc, powerRepo)
    endpoints := gencalc.NewEndpoints(calcSvc)
	

    //intialize transport
	grpcserver := calcsvr.New(endpoints, nil)
	grpcsvr := grpc.NewServer(
		grpcmiddleware.WithUnaryServerChain(
			goagrpcmiddleware.UnaryRequestID(),
			log.UnaryServerInterceptor(ctx),
			goagrpcmiddleware.UnaryServerLogContext(log.AsGoaMiddlewareLogger),
			metrics.UnaryServerInterceptor(ctx),
			trace.UnaryServerInterceptor(ctx),
		),
		grpcmiddleware.WithStreamServerChain(
			goagrpcmiddleware.StreamRequestID(),
			log.StreamServerInterceptor(ctx),
			goagrpcmiddleware.StreamServerLogContext(log.AsGoaMiddlewareLogger),
			metrics.StreamServerInterceptor(ctx),
			trace.StreamServerInterceptor(ctx),
		),
	)
	calcpb.RegisterCalcServer(grpcsvr, grpcserver)
	reflection.Register(grpcsvr)
	for svc, info := range grpcsvr.GetServiceInfo() {
		for _, m := range info.Methods {
			log.Print(ctx, log.KV{K: "method", V: svc + "/" + m.Name})
		}
	}

    
	// Mount health check & metrics on separate handler to avoid logging etc.
	check := health.Handler(health.NewChecker(dbc))
	http.Handle("/healthz", check)
	http.Handle("/livez", check)
	http.Handle("/metrics", metrics.Handler(ctx).(http.HandlerFunc))
	httpsvr := &http.Server{Addr: *httpaddr}

	// Signal handler
	errc := make(chan error)
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errc <- fmt.Errorf("%s", <-c)
	}()

	// Listen loop
	ctx, cancel := context.WithCancel(ctx)
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()

		go func() {
			log.Printf(ctx, "HTTP server listening on %s", httpsvr.Addr)
			errc <- httpsvr.ListenAndServe()
		}()

		go func() {
			l, err := net.Listen("tcp", *grpcaddr)
			if err != nil {
				errc <- fmt.Errorf("failed to start gRPC listener: %s", err.Error())
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
	if err := <-errc; err != nil {
		log.Errorf(ctx, err, "exiting")
	}
	cancel()
	wg.Wait()
	log.Printf(ctx, "exited")
}
