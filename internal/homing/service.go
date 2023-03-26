package homing

import (
	"context"

	"github.com/mainflux/mainflux"
	"github.com/mainflux/mainflux/pkg/errors"
	"golang.org/x/exp/slices"
)

const (
	usersObjectKey    = "users"
	memberRelationKey = "member"
)

type Service interface {
	Save(ctx context.Context, t Telemetry, serviceName string) error
	GetAll(ctx context.Context, token string, pm PageMetadata) (TelemetryPage, error)
}

var _ Service = (*telemetryService)(nil)

type telemetryService struct {
	repo   TelemetryRepo
	locSvc LocationService
	auth   mainflux.AuthServiceClient
}

func New(repo TelemetryRepo, locSvc LocationService, auth mainflux.AuthServiceClient) Service {
	return &telemetryService{
		repo:   repo,
		locSvc: locSvc,
		auth:   auth,
	}
}

// GetAll implements Service
func (ts *telemetryService) GetAll(ctx context.Context, token string, pm PageMetadata) (TelemetryPage, error) {
	res, err := ts.auth.Identify(ctx, &mainflux.Token{Value: token})
	if err != nil {
		return TelemetryPage{}, err
	}

	if err := ts.authorize(ctx, res.GetId(), usersObjectKey, memberRelationKey); err != nil {
		return TelemetryPage{}, err
	}
	telemetry, err := ts.repo.RetrieveAll(ctx, pm)

	return TelemetryPage{
		Telemetry:    telemetry,
		PageMetadata: pm,
	}, err
}

// Save implements Service
func (ts *telemetryService) Save(ctx context.Context, t Telemetry, serviceName string) error {
	telemetry, err := ts.repo.RetrieveByIP(ctx, t.IpAddress)
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
