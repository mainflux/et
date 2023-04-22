package homing

import (
	"context"
	"fmt"
	"reflect"

	goerrors "errors"

	"github.com/google/uuid"
	"github.com/mainflux/callhome/internal/homing/repository"
	"github.com/mainflux/mainflux"
	"github.com/mainflux/mainflux/pkg/errors"
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
	GetAll(ctx context.Context, repo, token string, pm PageMetadata) (TelemetryPage, error)
}

var _ Service = (*telemetryService)(nil)

type telemetryService struct {
	repo          TelemetryRepo
	timescaleRepo TelemetryRepo
	locSvc        LocationService
	auth          mainflux.AuthServiceClient
}

// New creates a new instance of the telemetry service.
func New(timescaleRepo, repo TelemetryRepo, locSvc LocationService, auth mainflux.AuthServiceClient) Service {
	return &telemetryService{
		repo:          repo,
		locSvc:        locSvc,
		auth:          auth,
		timescaleRepo: timescaleRepo,
	}
}

// GetAll retrieves homing telemetry data from the specified repository.
func (ts *telemetryService) GetAll(ctx context.Context, repo, token string, pm PageMetadata) (TelemetryPage, error) {
	res, err := ts.auth.Identify(ctx, &mainflux.Token{Value: token})
	if err != nil {
		return TelemetryPage{}, err
	}
	if err := ts.authorize(ctx, res.GetId(), usersObjectKey, memberRelationKey); err != nil {
		return TelemetryPage{}, err
	}

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

func (ts *telemetryService) authorize(ctx context.Context, subject, object string, relation string) error {
	req := &mainflux.AuthorizeReq{
		Sub: subject,
		Obj: object,
		Act: relation,
	}
	res, err := ts.auth.Authorize(ctx, req)
	if err != nil {
		return errors.Wrap(errors.ErrAuthorization, err)
	}
	if !res.GetAuthorized() {
		return errors.ErrAuthorization
	}
	return nil
}
