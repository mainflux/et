package callhome

import (
	"context"
	"time"
)

const (
	usersObjectKey    = "users"
	memberRelationKey = "member"
)

// Service Service to receive homing telemetry data, persist and retrieve it.
type Service interface {
	// Save saves the homing telemetry data and its location information.
	Save(ctx context.Context, t Telemetry) error
	// Retrieve retrieves homing telemetry data from the specified repository.
	Retrieve(ctx context.Context, pm PageMetadata) (TelemetryPage, error)
}

var _ Service = (*telemetryService)(nil)

type telemetryService struct {
	repo   TelemetryRepo
	locSvc LocationService
}

// New creates a new instance of the telemetry service.
func New(repo TelemetryRepo, locSvc LocationService) Service {
	return &telemetryService{
		repo:   repo,
		locSvc: locSvc,
	}
}

// Retrieve retrieves homing telemetry data from the specified repository.
func (ts *telemetryService) Retrieve(ctx context.Context, pm PageMetadata) (TelemetryPage, error) {
	return ts.repo.RetrieveAll(ctx, pm)
}

// Save saves the homing telemetry data and its location information.
func (ts *telemetryService) Save(ctx context.Context, t Telemetry) error {
	locRec, err := ts.locSvc.GetLocation(t.IpAddress)
	if err != nil {
		return err
	}
	t.City = locRec.City
	t.Country = locRec.Country_long
	t.Latitude = float64(locRec.Latitude)
	t.Longitude = float64(locRec.Longitude)
	t.LastSeen = time.Now()
	return ts.repo.Save(ctx, t)
}
