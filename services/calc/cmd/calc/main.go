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
	"github.com/crossnokaye/carbon/model"
    gencalc "github.com/crossnokaye/carbon/services/calc/gen/calc"
    calcpb "github.com/crossnokaye/carbon/services/calc/gen/grpc/calc/pb"
    calcsvr "github.com/crossnokaye/carbon/services/calc/gen/grpc/calc/server"
    calcapi "github.com/crossnokaye/carbon/services/calc"
    pastvalsvc "github.com/crossnokaye/carbon/services/calc/clients/power"
    "github.com/crossnokaye/carbon/clients/clickhouse"
    "github.com/crossnokaye/carbon/services/calc/clients/storage"
	"github.com/crossnokaye/carbon/services/calc/clients/facilityconfig"
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

	var (
		grpcaddr  = flag.String("grpc-addr", "0.0.0.0:12200", "gRPC listen address")
		httpaddr  = flag.String("http-addr", "0.0.0.0:12201", "HTTP listen address")
		agentaddr = flag.String("agent-addr", ":4317", "Grafana agent listen address")
		chaddr = flag.String("ch-addr", os.Getenv("CLICKHOUSE_ADDR"), "ClickHouse host address")
		chuser = flag.String("ch-user", os.Getenv("CLICKHOUSE_USER"), "ClickHouse user")
		chpwd  = flag.String("ch-pwd", os.Getenv("CLICKHOUSE_PASSWORD"), "ClickHouse password")
		chssl  = flag.Bool("ch-ssl", os.Getenv("CLICKHOUSE_SSL") != "", "ClickHouse connection SSL")
		monitoringEnabled = flag.Bool("monitoring-enabled", true, "monitoring")
		debug = flag.Bool("debug", false, "debug mode")
        //past val service host address
        pastValaddr = flag.String("pastval-add", "localhost:10110", "Past-Value host address")
        env = flag.String("dev", os.Getenv("ENV"), "facility environment") //what is this?
	)
	flag.Parse()

	// Enable logging
	format := log.FormatJSON
	if log.IsTerminal() {
		format = log.FormatTerminal
	}
	ctx := log.With(log.Context(context.Background(), log.WithFormat(format)), log.KV{K: "svc", V: gencalc.ServiceName})
	//log clickhouse credentials and monitoring status
	log.Info(ctx,
		log.KV{K: "pastval-add", V: *pastValaddr},
        log.KV{K: "ch-addr", V: *chaddr},
        log.KV{K: "ch-user", V: *chuser})
		log.Info(ctx,
			log.KV{K: "monitoringEnabled", V: *monitoringEnabled})
	
	if *debug {
		ctx = log.Context(ctx, log.WithDebug())
		log.Debugf(ctx, "debug logs enabled")
	}
	if *monitoringEnabled {
		conn, err := grpc.DialContext(ctx, *agentaddr,
			grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			log.Errorf(ctx, err, "failed to connect to Grafana agent")
			os.Exit(1)
		}
		log.Debugf(ctx, "connected to Grafana agent %s", agentaddr)
		if ctx, err = trace.Context(ctx, gencalc.ServiceName, trace.WithGRPCExporter(conn)); err != nil {
			log.Errorf(ctx, err, "failed to initialize tracing")
			os.Exit(1)
		}
	}
	//initialize the metrics
	ctx = metrics.Context(ctx, gencalc.ServiceName)

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
    facilityRepo := facilityconfig.New(*env)

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
	calcSvc = calcapi.NewCalc(ctx, pwc, dbc, facilityRepo)
    endpoints := gencalc.NewEndpoints(calcSvc)
	//remove later

	orgID := "52858b15-16ce-4998-b317-a1ce68c348c3"
	facID := "a5746ffa-2073-455e-b811-322ad3c3c4b7"
	locID := "cf153258-c08f-4ff0-9b01-d51d452e40e5"
	Interval := model.Hourly
	timeFormat := "2006-01-02T15:04:05-07:00"
	var startTime = time.Date(2021, time.June, 1, 0, 0, 0, 0, time.UTC)
	var endTime = time.Date(2021, time.June, 1, 1, 10, 0, 0, time.UTC)
	duration := &gencalc.Period{StartTime: startTime.Format(timeFormat), EndTime: endTime.Format(timeFormat)}
	payload := &gencalc.RequestPayload{OrgID: gencalc.UUID(orgID),
		 Duration: duration, FacilityID: gencalc.UUID(facID), Interval: Interval, LocationID: gencalc.UUID(locID)}
	calcSvc.HistoricalCarbonEmissions(ctx, payload)

    //intialize transport
	grpcserver := calcsvr.New(endpoints, nil)
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
