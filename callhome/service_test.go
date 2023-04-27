package callhome_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/ip2location/ip2location-go/v9"
	"github.com/mainflux/callhome/callhome"
	"github.com/mainflux/callhome/callhome/mocks"
	"github.com/mainflux/callhome/callhome/repository"
	repoMocks "github.com/mainflux/callhome/callhome/repository/mocks"
	"github.com/mainflux/mainflux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestRetrieve(t *testing.T) {

	t.Run("failed to identify", func(t *testing.T) {
		sheetRepo := repoMocks.NewTelemetryRepo(t)
		timescaleRepo := repoMocks.NewTelemetryRepo(t)
		authMock := mocks.NewAuthMockRepo(t)
		svc := callhome.New(timescaleRepo, sheetRepo, nil)
		experr := fmt.Errorf("failed to identify")
		authMock.On("Identify", context.Background(), &mainflux.Token{}, mock.Anything).Return(&mainflux.UserIdentity{}, experr)

		_, err := svc.Retrieve(context.Background(), callhome.SheetsRepo, callhome.PageMetadata{})
		assert.NotNil(t, err)
		assert.Equal(t, experr, err)
	})
	t.Run("failed authentication", func(t *testing.T) {
		sheetRepo := repoMocks.NewTelemetryRepo(t)
		timescaleRepo := repoMocks.NewTelemetryRepo(t)
		authMock := mocks.NewAuthMockRepo(t)
		svc := callhome.New(timescaleRepo, sheetRepo, nil)
		experr := fmt.Errorf("failed authentication")
		authMock.On("Identify", context.Background(), &mainflux.Token{}, mock.Anything).Return(&mainflux.UserIdentity{}, nil)
		authMock.On("Authorize", context.Background(), &mainflux.AuthorizeReq{Obj: "users", Act: "member"}, mock.Anything).Return(&mainflux.AuthorizeRes{Authorized: false}, experr)
		_, err := svc.Retrieve(context.Background(), callhome.SheetsRepo, callhome.PageMetadata{})
		assert.NotNil(t, err)

	})
	t.Run("failed repo save", func(t *testing.T) {
		sheetRepo := repoMocks.NewTelemetryRepo(t)
		timescaleRepo := repoMocks.NewTelemetryRepo(t)
		authMock := mocks.NewAuthMockRepo(t)
		svc := callhome.New(timescaleRepo, sheetRepo, nil)
		sheetRepo.On("RetrieveAll", context.Background(), callhome.PageMetadata{}).Return(callhome.TelemetryPage{}, repository.ErrSaveEvent)
		authMock.On("Identify", context.Background(), &mainflux.Token{}, mock.Anything).Return(&mainflux.UserIdentity{}, nil)
		authMock.On("Authorize", context.Background(), &mainflux.AuthorizeReq{Obj: "users", Act: "member"}, mock.Anything).Return(&mainflux.AuthorizeRes{Authorized: true}, nil)
		_, err := svc.Retrieve(context.Background(), callhome.SheetsRepo, callhome.PageMetadata{})
		assert.NotNil(t, err)
		assert.Equal(t, repository.ErrSaveEvent, err)
	})
	t.Run("success", func(t *testing.T) {
		sheetRepo := repoMocks.NewTelemetryRepo(t)
		timescaleRepo := repoMocks.NewTelemetryRepo(t)
		authMock := mocks.NewAuthMockRepo(t)
		svc := callhome.New(timescaleRepo, sheetRepo, nil)
		sheetRepo.On("RetrieveAll", context.Background(), callhome.PageMetadata{}).Return(callhome.TelemetryPage{}, nil)
		authMock.On("Identify", context.Background(), &mainflux.Token{}, mock.Anything).Return(&mainflux.UserIdentity{}, nil)
		authMock.On("Authorize", context.Background(), &mainflux.AuthorizeReq{Obj: "users", Act: "member"}, mock.Anything).Return(&mainflux.AuthorizeRes{Authorized: true}, nil)
		_, err := svc.Retrieve(context.Background(), callhome.SheetsRepo, callhome.PageMetadata{})
		assert.Nil(t, err)
	})
}

func TestSave(t *testing.T) {
	t.Run("error obtaining location", func(t *testing.T) {
		sheetRepo := repoMocks.NewTelemetryRepo(t)
		timescaleRepo := repoMocks.NewTelemetryRepo(t)
		locMock := mocks.NewLocationService(t)
		locMock.On("GetLocation", "").Return(ip2location.IP2Locationrecord{}, fmt.Errorf("error getting loc"))
		svc := callhome.New(timescaleRepo, sheetRepo, locMock)
		err := svc.Save(context.Background(), callhome.Telemetry{})
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
		timescaleRepo.On("Save", context.Background(), mock.AnythingOfType("callhome.Telemetry")).Return(repository.ErrSaveEvent)
		svc := callhome.New(timescaleRepo, sheetRepo, locMock)
		err := svc.Save(context.Background(), callhome.Telemetry{})
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
		timescaleRepo.On("Save", context.Background(), mock.AnythingOfType("callhome.Telemetry")).Return(nil)
		sheetRepo.On("RetrieveByIP", context.Background(), "").Return(callhome.Telemetry{}, fmt.Errorf("error getting record"))
		svc := callhome.New(timescaleRepo, sheetRepo, locMock)
		err := svc.Save(context.Background(), callhome.Telemetry{})
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
		timescaleRepo.On("Save", context.Background(), mock.AnythingOfType("callhome.Telemetry")).Return(nil)
		sheetRepo.On("RetrieveByIP", context.Background(), "").Return(callhome.Telemetry{}, repository.ErrRecordNotFound)
		sheetRepo.On("Save", context.Background(), mock.AnythingOfType("callhome.Telemetry")).Return(nil)
		svc := callhome.New(timescaleRepo, sheetRepo, locMock)
		err := svc.Save(context.Background(), callhome.Telemetry{})
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
		timescaleRepo.On("Save", context.Background(), mock.AnythingOfType("callhome.Telemetry")).Return(nil)
		sheetRepo.On("RetrieveByIP", context.Background(), "").Return(callhome.Telemetry{ID: uuid.NewString()}, nil)
		sheetRepo.On("Save", context.Background(), mock.AnythingOfType("callhome.Telemetry")).Return(nil)
		sheetRepo.On("Update", context.Background(), mock.AnythingOfType("callhome.Telemetry")).Return(nil)
		svc := callhome.New(timescaleRepo, sheetRepo, locMock)
		err := svc.Save(context.Background(), callhome.Telemetry{})
		assert.Nil(t, err)
	})
}
