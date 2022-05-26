package telemetry

import (
	// stdlib
	"io"
	"time"

	// external
	"go.opencensus.io/stats/view"
	"go.opencensus.io/trace"
)

// Setup OpenCensus
func Setup(serviceName string) io.Closer {
	// Always trace for this demo.
	trace.ApplyConfig(trace.Config{DefaultSampler: trace.AlwaysSample()})

	// Report stats at every second.
	view.SetReportingPeriod(1 * time.Second)

	zipkinExporter, zipkinCloser := setupZipkin(serviceName)

	exporter := &multiExporter{
		t: []trace.Exporter{
			zipkinExporter,
		},
		v: []view.Exporter{},
		c: []io.Closer{
			zipkinCloser,
		},
	}

	trace.RegisterExporter(exporter)
	view.RegisterExporter(exporter)

	return exporter
}
