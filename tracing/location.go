// Copyright (c) Abstract Machines
// SPDX-License-Identifier: Apache-2.0

package tracing

import (
	"context"

	"github.com/absmach/callhome"
	"github.com/ip2location/ip2location-go/v9"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

const locOpName = "get_location_op"

var _ callhome.LocationService = (*locationServiceTracer)(nil)

type locationServiceTracer struct {
	svc    callhome.LocationService
	tracer trace.Tracer
}

// NewLocationService adds tracing middlware to callhome.LocationService.
func NewLocationService(tracer trace.Tracer, svc callhome.LocationService) callhome.LocationService {
	return &locationServiceTracer{
		tracer: tracer,
		svc:    svc,
	}
}

// GetLocation adds tracing middleware to location service.
func (lst *locationServiceTracer) GetLocation(ctx context.Context, ip string) (ip2location.IP2Locationrecord, error) {
	ctx, span := lst.tracer.Start(ctx, locOpName, trace.WithAttributes([]attribute.KeyValue{attribute.String("ip_address", ip)}...))
	defer span.End()
	return lst.svc.GetLocation(ctx, ip)
}
