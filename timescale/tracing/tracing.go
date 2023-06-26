package tracing

import (
	"context"

	"github.com/mainflux/callhome"
	"go.opentelemetry.io/otel/trace"
)

const (
	retrieveAllOp                  = "retrieve_all_op"
	retrieveDistinctIPsCountriesOp = "retrieve_distinct_IP_countries_op"
	saveOp                         = "save_op"
)

var _ callhome.TelemetryRepo = (*repoTracer)(nil)

type repoTracer struct {
	repo   callhome.TelemetryRepo
	tracer trace.Tracer
}

// New adds tracing middleware to callhome.TelemetryRepo.
func New(tracer trace.Tracer, repo callhome.TelemetryRepo) callhome.TelemetryRepo {
	return &repoTracer{
		tracer: tracer,
		repo:   repo,
	}
}

// RetrieveAll adds tracing middleware to retrieve all method.
func (rt *repoTracer) RetrieveAll(ctx context.Context, pm callhome.PageMetadata) (callhome.TelemetryPage, error) {
	ctx, span := rt.tracer.Start(ctx, retrieveAllOp)
	defer span.End()
	return rt.repo.RetrieveAll(ctx, pm)
}

// RetrieveDistinctIPsCountries adds tracing middleware to retrieve distinct ips countries method.
func (rt *repoTracer) RetrieveDistinctIPsCountries(ctx context.Context) (callhome.TelemetrySummary, error) {
	ctx, span := rt.tracer.Start(ctx, retrieveDistinctIPsCountriesOp)
	defer span.End()
	return rt.repo.RetrieveDistinctIPsCountries(ctx)
}

// Save adds tracing middleware to save method.
func (rt *repoTracer) Save(ctx context.Context, t callhome.Telemetry) error {
	ctx, span := rt.tracer.Start(ctx, saveOp)
	defer span.End()
	return rt.repo.Save(ctx, t)
}
