package homing

import (
	"context"

	"golang.org/x/exp/slices"
)

type Service interface {
	Save(ctx context.Context, t Telemetry, serviceName string) error
	GetAll(ctx context.Context, token string, pm PageMetadata) (TelemetryPage, error)
}

var _ Service = (*telemetryService)(nil)

type telemetryService struct {
	repo   TelemetryRepo
	locSvc LocationService
}

func New(repo TelemetryRepo, locSvc LocationService) Service {
	return &telemetryService{
		repo:   repo,
		locSvc: locSvc,
	}
}

// GetAll implements Service
func (ts *telemetryService) GetAll(ctx context.Context, token string, pm PageMetadata) (TelemetryPage, error) {
	telemetry, err := ts.repo.RetrieveAll(ctx, pm)

	return TelemetryPage{
		Telemetry:    telemetry,
		PageMetadata: pm,
	}, err
}

// Save implements Service
func (ts *telemetryService) Save(ctx context.Context, t Telemetry, serviceName string) error {
	telemetry, row, err := ts.repo.RetrieveByIP(ctx, t.IpAddress)
	if err != nil {
		return err
	}
	long, lat, err := ts.locSvc.GetLocation(t.IpAddress)
	if err != nil {
		return err
	}

	t.Latitutde = float64(lat)
	t.Longitude = float64(long)

	if telemetry == nil {
		t.Services = append(t.Services, serviceName)
		err = ts.repo.Save(ctx, t)
		return err
	}
	t.ID = telemetry.ID
	t.Services = telemetry.Services

	if !slices.Contains(t.Services, serviceName) {
		t.Services = append(t.Services, serviceName)
	}

	return ts.repo.UpdateTelemetry(ctx, t, row)

}
