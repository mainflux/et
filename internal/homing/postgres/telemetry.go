package postgres

import (
	"context"

	"github.com/mainflux/et/internal/homing"
)

var _ homing.TelemetryRepo = (*repo)(nil)

type repo struct {
}

// RetrieveAll implements homing.TelemetryRepo
func (*repo) RetrieveAll(ctx context.Context, pm homing.PageMetadata) ([]homing.Telemetry, error) {
	panic("unimplemented")
}

// RetrieveByIP implements homing.TelemetryRepo
func (*repo) RetrieveByIP(ctx context.Context, email string) (*homing.Telemetry, error) {
	panic("unimplemented")
}

// Save implements homing.TelemetryRepo
func (*repo) Save(ctx context.Context, t homing.Telemetry) error {
	panic("unimplemented")
}

// UpdateTelemetry implements homing.TelemetryRepo
func (*repo) UpdateTelemetry(ctx context.Context, u homing.Telemetry) error {
	panic("unimplemented")
}
