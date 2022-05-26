package telemetry

import (
	// stdlib
	"errors"
	"io"
	"strings"

	// external
	"go.opencensus.io/stats/view"
	"go.opencensus.io/trace"
)

type multiExporter struct {
	t []trace.Exporter
	v []view.Exporter
	c []io.Closer
}

// ExportSpan implements trace.Exporter
func (m *multiExporter) ExportSpan(s *trace.SpanData) {
	for _, t := range m.t {
		t.ExportSpan(s)
	}
}

// ExportView implements view.Exporter
func (m *multiExporter) ExportView(d *view.Data) {
	for _, v := range m.v {
		v.ExportView(d)
	}
}

// Close implements io.Closer
func (m *multiExporter) Close() error {
	var e []string

	for _, c := range m.c {
		if err := c.Close(); err != nil {
			e = append(e, err.Error())
		}
	}

	if len(e) > 0 {
		return errors.New("ERRORS: " + strings.Join(e, " ; "))
	}
	return nil
}
