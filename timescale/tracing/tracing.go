// Copyright (c) Abstract Machines
// SPDX-License-Identifier: Apache-2.0

package tracing

import (
	"context"

	"github.com/absmach/callhome"
	"go.opentelemetry.io/otel/trace"
)

const (
	retrieveAllOp     = "retrieve_all_op"
	retrieveSummaryOp = "retrieve_summary_op"
	saveOp            = "save_op"
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
func (rt *repoTracer) RetrieveAll(ctx context.Context, pm callhome.PageMetadata, filter callhome.TelemetryFilters) (callhome.TelemetryPage, error) {
	ctx, span := rt.tracer.Start(ctx, retrieveAllOp)
	defer span.End()
	return rt.repo.RetrieveAll(ctx, pm, filter)
}

// RetrieveSummary adds tracing middleware to retrieve summary method.
func (rt *repoTracer) RetrieveSummary(ctx context.Context, filter callhome.TelemetryFilters) (callhome.TelemetrySummary, error) {
	ctx, span := rt.tracer.Start(ctx, retrieveSummaryOp)
	defer span.End()
	return rt.repo.RetrieveSummary(ctx, filter)
}

// Save adds tracing middleware to save method.
func (rt *repoTracer) Save(ctx context.Context, t callhome.Telemetry) error {
	ctx, span := rt.tracer.Start(ctx, saveOp)
	defer span.End()
	return rt.repo.Save(ctx, t)
}
