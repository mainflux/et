package callhome

import (
	"context"
	"errors"
	"reflect"
	"time"

	"github.com/mainflux/callhome/callhome/repository"
	"golang.org/x/exp/slices"
)

const (
	usersObjectKey    = "users"
	memberRelationKey = "member"
	SheetsRepo        = "sheets"
	TimescaleRepo     = "timescale"
)

var errInvalidRepo error = errors.New("undefined repository")

// Service Service to receive homing telemetry data, persist and retrieve it.
type Service interface {
	// Save saves the homing telemetry data and its location information.
	Save(ctx context.Context, t Telemetry) error
	// Retrieve retrieves homing telemetry data from the specified repository.
	Retrieve(ctx context.Context, repo string, pm PageMetadata) (TelemetryPage, error)
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

// Retrieve retrieves homing telemetry data from the specified repository.
func (ts *telemetryService) Retrieve(ctx context.Context, repo string, pm PageMetadata) (TelemetryPage, error) {
	switch repo {
	case SheetsRepo:
		return ts.repo.RetrieveAll(ctx, pm)
	case TimescaleRepo:
		return ts.timescaleRepo.RetrieveAll(ctx, pm)
	default:
		return TelemetryPage{}, errInvalidRepo
	}
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

	if err := ts.timescaleRepo.Save(ctx, t); err != nil {
		return err
	}

	telemetry, err := ts.repo.RetrieveByIP(ctx, t.IpAddress)
	if err != nil && !errors.Is(err, repository.ErrRecordNotFound) {
		return err
	}
	if reflect.ValueOf(telemetry).IsZero() {
		t.Services = append(t.Services, t.Service)
		return ts.repo.Save(ctx, t)
	}
	t.Services = telemetry.Services
	if !slices.Contains(t.Services, t.Service) {
		t.Services = append(t.Services, t.Service)
	}
	return ts.repo.Update(ctx, t)
}
