package homing_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/ip2location/ip2location-go/v9"
	"github.com/mainflux/callhome/internal/homing"
	"github.com/mainflux/callhome/internal/homing/mocks"
	"github.com/mainflux/callhome/internal/homing/repository"
	repoMocks "github.com/mainflux/callhome/internal/homing/repository/mocks"
	"github.com/mainflux/mainflux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetAll(t *testing.T) {

	t.Run("failed to identify", func(t *testing.T) {
		sheetRepo := repoMocks.NewTelemetryRepo(t)
		timescaleRepo := repoMocks.NewTelemetryRepo(t)
		authMock := mocks.NewAuthMockRepo(t)
		svc := homing.New(timescaleRepo, sheetRepo, nil)
		experr := fmt.Errorf("failed to identify")
		authMock.On("Identify", context.Background(), &mainflux.Token{}, mock.Anything).Return(&mainflux.UserIdentity{}, experr)

		_, err := svc.Retrieve(context.Background(), homing.SheetsRepo, homing.PageMetadata{})
		assert.NotNil(t, err)
		assert.Equal(t, experr, err)
	})
	t.Run("failed authentication", func(t *testing.T) {
		sheetRepo := repoMocks.NewTelemetryRepo(t)
		timescaleRepo := repoMocks.NewTelemetryRepo(t)
		authMock := mocks.NewAuthMockRepo(t)
		svc := homing.New(timescaleRepo, sheetRepo, nil)
		experr := fmt.Errorf("failed authentication")
		authMock.On("Identify", context.Background(), &mainflux.Token{}, mock.Anything).Return(&mainflux.UserIdentity{}, nil)
		authMock.On("Authorize", context.Background(), &mainflux.AuthorizeReq{Obj: "users", Act: "member"}, mock.Anything).Return(&mainflux.AuthorizeRes{Authorized: false}, experr)
		_, err := svc.Retrieve(context.Background(), homing.SheetsRepo, homing.PageMetadata{})
		assert.NotNil(t, err)

	})
	t.Run("failed repo save", func(t *testing.T) {
		sheetRepo := repoMocks.NewTelemetryRepo(t)
		timescaleRepo := repoMocks.NewTelemetryRepo(t)
		authMock := mocks.NewAuthMockRepo(t)
		svc := homing.New(timescaleRepo, sheetRepo, nil)
		sheetRepo.On("RetrieveAll", context.Background(), homing.PageMetadata{}).Return(homing.TelemetryPage{}, repository.ErrSaveEvent)
		authMock.On("Identify", context.Background(), &mainflux.Token{}, mock.Anything).Return(&mainflux.UserIdentity{}, nil)
		authMock.On("Authorize", context.Background(), &mainflux.AuthorizeReq{Obj: "users", Act: "member"}, mock.Anything).Return(&mainflux.AuthorizeRes{Authorized: true}, nil)
		_, err := svc.Retrieve(context.Background(), homing.SheetsRepo, homing.PageMetadata{})
		assert.NotNil(t, err)
		assert.Equal(t, repository.ErrSaveEvent, err)
	})
	t.Run("success", func(t *testing.T) {
		sheetRepo := repoMocks.NewTelemetryRepo(t)
		timescaleRepo := repoMocks.NewTelemetryRepo(t)
		authMock := mocks.NewAuthMockRepo(t)
		svc := homing.New(timescaleRepo, sheetRepo, nil)
		sheetRepo.On("RetrieveAll", context.Background(), homing.PageMetadata{}).Return(homing.TelemetryPage{}, nil)
		authMock.On("Identify", context.Background(), &mainflux.Token{}, mock.Anything).Return(&mainflux.UserIdentity{}, nil)
		authMock.On("Authorize", context.Background(), &mainflux.AuthorizeReq{Obj: "users", Act: "member"}, mock.Anything).Return(&mainflux.AuthorizeRes{Authorized: true}, nil)
		_, err := svc.Retrieve(context.Background(), homing.SheetsRepo, homing.PageMetadata{})
		assert.Nil(t, err)
	})
}

func TestSave(t *testing.T) {
	t.Run("error obtaining location", func(t *testing.T) {
		sheetRepo := repoMocks.NewTelemetryRepo(t)
		timescaleRepo := repoMocks.NewTelemetryRepo(t)
		locMock := mocks.NewLocationService(t)
		locMock.On("GetLocation", "").Return(ip2location.IP2Locationrecord{}, fmt.Errorf("error getting loc"))
		svc := homing.New(timescaleRepo, sheetRepo, locMock)
		err := svc.Save(context.Background(), homing.Telemetry{})
		assert.NotNil(t, err)
	})
	t.Run("error saving to timescale", func(t *testing.T) {
		sheetRepo := repoMocks.NewTelemetryRepo(t)
		timescaleRepo := repoMocks.NewTelemetryRepo(t)
		locMock := mocks.NewLocationService(t)
		locMock.On("GetLocation", "").Return(ip2location.IP2Locationrecord{
			Latitude:     1.2,
			Longitude:    30,
			Country_long: "SomeCountry",
			City:         "someCity",
		}, nil)
		timescaleRepo.On("Save", context.Background(), mock.AnythingOfType("homing.Telemetry")).Return(repository.ErrSaveEvent)
		svc := homing.New(timescaleRepo, sheetRepo, locMock)
		err := svc.Save(context.Background(), homing.Telemetry{})
		assert.NotNil(t, err)
		assert.Equal(t, repository.ErrSaveEvent, err)
	})
	t.Run("error retrieve by ip", func(t *testing.T) {
		sheetRepo := repoMocks.NewTelemetryRepo(t)
		timescaleRepo := repoMocks.NewTelemetryRepo(t)
		locMock := mocks.NewLocationService(t)
		locMock.On("GetLocation", "").Return(ip2location.IP2Locationrecord{
			Latitude:     1.2,
			Longitude:    30,
			Country_long: "SomeCountry",
			City:         "someCity",
		}, nil)
		timescaleRepo.On("Save", context.Background(), mock.AnythingOfType("homing.Telemetry")).Return(nil)
		sheetRepo.On("RetrieveByIP", context.Background(), "").Return(homing.Telemetry{}, fmt.Errorf("error getting record"))
		svc := homing.New(timescaleRepo, sheetRepo, locMock)
		err := svc.Save(context.Background(), homing.Telemetry{})
		assert.NotNil(t, err)
	})
	t.Run("successful save", func(t *testing.T) {
		sheetRepo := repoMocks.NewTelemetryRepo(t)
		timescaleRepo := repoMocks.NewTelemetryRepo(t)
		locMock := mocks.NewLocationService(t)
		locMock.On("GetLocation", "").Return(ip2location.IP2Locationrecord{
			Latitude:     1.2,
			Longitude:    30,
			Country_long: "SomeCountry",
			City:         "someCity",
		}, nil)
		timescaleRepo.On("Save", context.Background(), mock.AnythingOfType("homing.Telemetry")).Return(nil)
		sheetRepo.On("RetrieveByIP", context.Background(), "").Return(homing.Telemetry{}, repository.ErrRecordNotFound)
		sheetRepo.On("Save", context.Background(), mock.AnythingOfType("homing.Telemetry")).Return(nil)
		svc := homing.New(timescaleRepo, sheetRepo, locMock)
		err := svc.Save(context.Background(), homing.Telemetry{})
		assert.Nil(t, err)
	})
	t.Run("successful update", func(t *testing.T) {
		sheetRepo := repoMocks.NewTelemetryRepo(t)
		timescaleRepo := repoMocks.NewTelemetryRepo(t)
		locMock := mocks.NewLocationService(t)
		locMock.On("GetLocation", "").Return(ip2location.IP2Locationrecord{
			Latitude:     1.2,
			Longitude:    30,
			Country_long: "SomeCountry",
			City:         "someCity",
		}, nil)
		timescaleRepo.On("Save", context.Background(), mock.AnythingOfType("homing.Telemetry")).Return(nil)
		sheetRepo.On("RetrieveByIP", context.Background(), "").Return(homing.Telemetry{ID: uuid.NewString()}, nil)
		sheetRepo.On("Save", context.Background(), mock.AnythingOfType("homing.Telemetry")).Return(nil)
		sheetRepo.On("UpdateTelemetry", context.Background(), mock.AnythingOfType("homing.Telemetry")).Return(nil)
		svc := homing.New(timescaleRepo, sheetRepo, locMock)
		err := svc.Save(context.Background(), homing.Telemetry{})
		assert.Nil(t, err)
	})
}
