package callhome_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/ip2location/ip2location-go/v9"
	"github.com/mainflux/callhome/callhome"
	"github.com/mainflux/callhome/callhome/mocks"
	"github.com/mainflux/callhome/callhome/repository"
	repoMocks "github.com/mainflux/callhome/callhome/repository/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestRetrieve(t *testing.T) {
	ctx := context.TODO()

	t.Run("failed to identify", func(t *testing.T) {
		sheetRepo := repoMocks.NewTelemetryRepo(t)
		timescaleRepo := repoMocks.NewTelemetryRepo(t)
		svc := callhome.New(timescaleRepo, sheetRepo, nil)
		experr := fmt.Errorf("failed to identify")
		_, err := svc.Retrieve(ctx, callhome.SheetsRepo, callhome.PageMetadata{})
		assert.NotNil(t, err)
		assert.Equal(t, experr, err)
	})
	t.Run("failed authentication", func(t *testing.T) {
		sheetRepo := repoMocks.NewTelemetryRepo(t)
		timescaleRepo := repoMocks.NewTelemetryRepo(t)
		svc := callhome.New(timescaleRepo, sheetRepo, nil)
		_, err := svc.Retrieve(ctx, callhome.SheetsRepo, callhome.PageMetadata{})
		assert.NotNil(t, err)

	})
	t.Run("failed repo save", func(t *testing.T) {
		sheetRepo := repoMocks.NewTelemetryRepo(t)
		timescaleRepo := repoMocks.NewTelemetryRepo(t)
		svc := callhome.New(timescaleRepo, sheetRepo, nil)
		sheetRepo.On("RetrieveAll", ctx, callhome.PageMetadata{}).Return(callhome.TelemetryPage{}, repository.ErrSaveEvent)
		_, err := svc.Retrieve(ctx, callhome.SheetsRepo, callhome.PageMetadata{})
		assert.NotNil(t, err)
		assert.Equal(t, repository.ErrSaveEvent, err)
	})
	t.Run("success", func(t *testing.T) {
		sheetRepo := repoMocks.NewTelemetryRepo(t)
		timescaleRepo := repoMocks.NewTelemetryRepo(t)
		svc := callhome.New(timescaleRepo, sheetRepo, nil)
		sheetRepo.On("RetrieveAll", ctx, callhome.PageMetadata{}).Return(callhome.TelemetryPage{}, nil)
		_, err := svc.Retrieve(ctx, callhome.SheetsRepo, callhome.PageMetadata{})
		assert.Nil(t, err)
	})
}

func TestSave(t *testing.T) {
	ctx := context.TODO()
	t.Run("error obtaining location", func(t *testing.T) {
		sheetRepo := repoMocks.NewTelemetryRepo(t)
		timescaleRepo := repoMocks.NewTelemetryRepo(t)
		locMock := mocks.NewLocationService(t)
		locMock.On("GetLocation", "").Return(ip2location.IP2Locationrecord{}, fmt.Errorf("error getting loc"))
		svc := callhome.New(timescaleRepo, sheetRepo, locMock)
		err := svc.Save(ctx, callhome.Telemetry{})
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
		timescaleRepo.On("Save", ctx, mock.AnythingOfType("callhome.Telemetry")).Return(repository.ErrSaveEvent)
		svc := callhome.New(timescaleRepo, sheetRepo, locMock)
		err := svc.Save(ctx, callhome.Telemetry{})
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
		timescaleRepo.On("Save", ctx, mock.AnythingOfType("callhome.Telemetry")).Return(nil)
		sheetRepo.On("RetrieveByIP", ctx, "").Return(callhome.Telemetry{}, fmt.Errorf("error getting record"))
		svc := callhome.New(timescaleRepo, sheetRepo, locMock)
		err := svc.Save(ctx, callhome.Telemetry{})
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
		timescaleRepo.On("Save", ctx, mock.AnythingOfType("callhome.Telemetry")).Return(nil)
		sheetRepo.On("RetrieveByIP", ctx, "").Return(callhome.Telemetry{}, repository.ErrRecordNotFound)
		sheetRepo.On("Save", ctx, mock.AnythingOfType("callhome.Telemetry")).Return(nil)
		svc := callhome.New(timescaleRepo, sheetRepo, locMock)
		err := svc.Save(ctx, callhome.Telemetry{})
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
		timescaleRepo.On("Save", ctx, mock.AnythingOfType("callhome.Telemetry")).Return(nil)
		sheetRepo.On("RetrieveByIP", ctx, "").Return(callhome.Telemetry{}, nil)
		sheetRepo.On("Save", ctx, mock.AnythingOfType("callhome.Telemetry")).Return(nil)
		sheetRepo.On("Update", ctx, mock.AnythingOfType("callhome.Telemetry")).Return(nil)
		svc := callhome.New(timescaleRepo, sheetRepo, locMock)
		err := svc.Save(ctx, callhome.Telemetry{})
		assert.Nil(t, err)
	})
}
