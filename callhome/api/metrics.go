package api

import (
	"context"
	"time"

	"github.com/go-kit/kit/metrics"
	"github.com/mainflux/callhome/callhome"
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

// GetAll add metrics middleware to get all service.
func (mm *metricsMiddleware) Retrieve(ctx context.Context, repo string, pm callhome.PageMetadata) (callhome.TelemetryPage, error) {
	defer func(begin time.Time) {
		mm.counter.With("method", "get all").Add(1)
		mm.latency.With("method", "get all").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return mm.svc.Retrieve(ctx, repo, pm)
}

// Save adds metrics middleware to save service.
func (mm *metricsMiddleware) Save(ctx context.Context, t callhome.Telemetry) error {
	defer func(begin time.Time) {
		mm.counter.With("method", "save").Add(1)
		mm.latency.With("method", "save").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return mm.svc.Save(ctx, t)
}
