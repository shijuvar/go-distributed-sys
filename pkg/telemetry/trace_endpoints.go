package telemetry

import (
	"time"

	// external
	"github.com/go-kit/kit/endpoint"
	kitoc "github.com/go-kit/kit/tracing/opencensus"
	"go.opencensus.io/trace"
)

// BalancerType used in Retry logic
type BalancerType string

// BalancerTypes
const (
	Random     BalancerType = "random"
	RoundRobin BalancerType = "round robin"
)

// ClientTracingEndpoint adds our Endpoint Tracing middleware to the existing client
// side endpoint.
func ClientTracingEndpoint(operationName string, attrs ...trace.Attribute) endpoint.Middleware {
	attrs = append(
		attrs, trace.StringAttribute("gokit.endpoint.type", "client"),
	)
	return kitoc.TraceEndpoint(
		"gokit/endpoint "+operationName,
		kitoc.WithEndpointAttributes(attrs...),
	)
}

// ServerTracingEndpoint adds our Endpoint Tracing middleware to the existing server
// side endpoint.
func ServerTracingEndpoint(operationName string, attrs ...trace.Attribute) endpoint.Middleware {
	attrs = append(
		attrs, trace.StringAttribute("gokit.endpoint.type", "server"),
	)
	return kitoc.TraceEndpoint(
		"gokit/endpoint "+operationName,
		kitoc.WithEndpointAttributes(attrs...),
	)
}

// RetryEndpoint wraps a Go kit lb.Retry endpoint with an annotated span.
func RetryEndpoint(
	operationName string, balancer BalancerType, max int, timeout time.Duration,
) endpoint.Middleware {
	return kitoc.TraceEndpoint("gokit/retry "+operationName,
		kitoc.WithEndpointAttributes(
			trace.StringAttribute("gokit.balancer.type", string(balancer)),
			trace.StringAttribute("gokit.retry.timeout", timeout.String()),
			trace.Int64Attribute("gokit.retry.max_count", int64(max)),
		),
	)
}
