package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	kitoc "github.com/go-kit/kit/tracing/opencensus"
	kithttp "github.com/go-kit/kit/transport/http"

	ordersvc "github.com/shijuvar/go-distributed-sys/eventstream/gokit-ordersvc"
	"github.com/shijuvar/go-distributed-sys/eventstream/gokit-ordersvc/implementation"
	"github.com/shijuvar/go-distributed-sys/eventstream/gokit-ordersvc/middleware"

	"github.com/shijuvar/go-distributed-sys/eventstream/gokit-ordersvc/transport"
	httptransport "github.com/shijuvar/go-distributed-sys/eventstream/gokit-ordersvc/transport/http"
	"github.com/shijuvar/go-distributed-sys/pkg/telemetry"
)

func main() {
	var (
		httpAddr = flag.String("http.addr", ":8080", "HTTP listen address")
	)
	flag.Parse()
	// initialize our OpenCensus configuration and defer a clean-up
	defer telemetry.Setup("ordersvc").Close()

	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stderr)
		logger = log.NewSyncLogger(logger)
		logger = level.NewFilter(logger, level.AllowDebug())
		logger = log.With(logger,
			"svc", "ordersvc",
			"ts", log.DefaultTimestampUTC,
			"caller", log.DefaultCaller,
		)
	}

	level.Info(logger).Log("msg", "service started")
	defer level.Info(logger).Log("msg", "service ended")

	// Create Order Service
	var s ordersvc.Service
	{
		s = implementation.NewService()
		// Add service middlewares
		s = middleware.LoggingMiddleware(logger)(s)

	}
	// Create Go kit endpoints for the Order Service
	// Then decorates with endpoint middlewares
	var endpoints transport.Endpoints
	{
		endpoints = transport.MakeServerEndpoints(s)
		// Add endpoint middlewares
		// Add endpoint level middlewares here
		// Trace server side endpoints with open census
		endpoints = transport.Endpoints{
			CreateOrder:  telemetry.ServerTracingEndpoint("CreateOrder")(endpoints.CreateOrder),
			GetOrderByID: telemetry.ServerTracingEndpoint("GetOrderByID")(endpoints.GetOrderByID),
		}
	}
	var h http.Handler
	{
		ocTracing := kitoc.HTTPServerTrace()
		serverOptions := []kithttp.ServerOption{ocTracing}
		h = httptransport.NewService(endpoints, serverOptions, logger)
	}
	errs := make(chan error)
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errs <- fmt.Errorf("%s", <-c)
	}()

	go func() {
		level.Info(logger).Log("transport", "HTTP", "addr", *httpAddr)
		errs <- http.ListenAndServe(*httpAddr, h)
	}()

	level.Error(logger).Log("exit", <-errs)
}
