package api

import (
	"context"
	"time"

	"github.com/absmach/callhome"
	"github.com/go-kit/kit/metrics"
)

var _ callhome.Service = (*metricsMiddleware)(nil)

type metricsMiddleware struct {
	counter metrics.Counter
	latency metrics.Histogram
	svc     callhome.Service
}

// MetricsMiddleware instruments core service by tracking request count and latency.
func MetricsMiddleware(svc callhome.Service, counter metrics.Counter, latency metrics.Histogram) callhome.Service {
	return &metricsMiddleware{
		counter: counter,
		latency: latency,
		svc:     svc,
	}
}

// Retrieve add metrics middleware to retrieve service.
func (mm *metricsMiddleware) Retrieve(ctx context.Context, pm callhome.PageMetadata, filters callhome.TelemetryFilters) (callhome.TelemetryPage, error) {
	defer func(begin time.Time) {
		mm.counter.With("method", "retrieve").Add(1)
		mm.latency.With("method", "retrieve").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return mm.svc.Retrieve(ctx, pm, filters)
}

// Save adds metrics middleware to save service.
func (mm *metricsMiddleware) Save(ctx context.Context, t callhome.Telemetry) error {
	defer func(begin time.Time) {
		mm.counter.With("method", "save").Add(1)
		mm.latency.With("method", "save").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return mm.svc.Save(ctx, t)
}

// RetrieveSummary adds metrics middleware to retrieve summary service.
func (mm *metricsMiddleware) RetrieveSummary(ctx context.Context, filters callhome.TelemetryFilters) (callhome.TelemetrySummary, error) {
	defer func(begin time.Time) {
		mm.counter.With("method", "retrieve-summary").Add(1)
		mm.latency.With("method", "retrieve-summary").Observe(time.Since(begin).Seconds())
	}(time.Now())
	return mm.svc.RetrieveSummary(ctx, filters)
}

// ServeUI implements callhome.Service
func (mm *metricsMiddleware) ServeUI(ctx context.Context, filters callhome.TelemetryFilters) ([]byte, error) {
	defer func(begin time.Time) {
		mm.counter.With("method", "serve-ui").Add(1)
		mm.latency.With("method", "serve-ui").Observe(time.Since(begin).Seconds())
	}(time.Now())
	return mm.svc.ServeUI(ctx, filters)
}
