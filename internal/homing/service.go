package homing

import (
	"context"
	"fmt"
	"reflect"

	goerrors "errors"

	"github.com/google/uuid"
	"github.com/mainflux/callhome/internal/homing/repository"
	"golang.org/x/exp/slices"
)

const (
	usersObjectKey    = "users"
	memberRelationKey = "member"
	SheetsRepo        = "sheets"
	TimescaleRepo     = "timescale"
)

// Service Service to receive homing telemetry data, persist and retrieve it.
type Service interface {
	// Save saves the homing telemetry data and its location information.
	Save(ctx context.Context, t Telemetry) error
	// GetAll retrieves homing telemetry data from the specified repository.
	GetAll(ctx context.Context, repo string, pm PageMetadata) (TelemetryPage, error)
}

var _ Service = (*telemetryService)(nil)

type telemetryService struct {
	repo          TelemetryRepo
	timescaleRepo TelemetryRepo
	locSvc        LocationService
}

// New creates a new instance of the telemetry service.
func New(timescaleRepo, repo TelemetryRepo, locSvc LocationService) Service {
	return &telemetryService{
		repo:          repo,
		locSvc:        locSvc,
		timescaleRepo: timescaleRepo,
	}
}

// GetAll retrieves homing telemetry data from the specified repository.
func (ts *telemetryService) GetAll(ctx context.Context, repo string, pm PageMetadata) (TelemetryPage, error) {
	switch repo {
	case SheetsRepo:
		return ts.repo.RetrieveAll(ctx, pm)
	case TimescaleRepo:
		return ts.timescaleRepo.RetrieveAll(ctx, pm)
	default:
		return TelemetryPage{}, fmt.Errorf("undefined repository")
	}
}

// Save saves the homing telemetry data and its location information.
func (ts *telemetryService) Save(ctx context.Context, t Telemetry) error {
	locRec, err := ts.locSvc.GetLocation(t.IpAddress)
	if err != nil {
		return err
	}
	t.ID = uuid.New().String()
	t.City = locRec.City
	t.Country = locRec.Country_long
	t.Latitude = float64(locRec.Latitude)
	t.Longitude = float64(locRec.Longitude)

	if err := ts.timescaleRepo.Save(ctx, t); err != nil {
		return err
	}

	telemetry, err := ts.repo.RetrieveByIP(ctx, t.IpAddress)
	if err != nil && !goerrors.Is(err, repository.ErrRecordNotFound) {
		return err
	}
	if reflect.ValueOf(telemetry).IsZero() {
		t.Services = append(t.Services, t.Service)
		return ts.repo.Save(ctx, t)
	}
	t.ID = telemetry.ID
	t.Services = telemetry.Services
	if !slices.Contains(t.Services, t.Service) {
		t.Services = append(t.Services, t.Service)
	}
	return ts.repo.UpdateTelemetry(ctx, t)
}
