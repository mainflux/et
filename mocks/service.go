package mocks

import (
	"context"

	"github.com/mainflux/callhome"
	"github.com/stretchr/testify/mock"
)

var _ callhome.Service = (*Service)(nil)

type Service struct {
	mock.Mock
}

// ServeUI implements callhome.Service
func (*Service) ServeUI(ctx context.Context, filters callhome.TelemetryFilters) ([]byte, error) {
	return nil, nil
}

func (s *Service) Retrieve(ctx context.Context, pm callhome.PageMetadata, filters callhome.TelemetryFilters) (callhome.TelemetryPage, error) {
	ret := s.Called(ctx, pm)
	var r0 callhome.TelemetryPage
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, callhome.PageMetadata) (callhome.TelemetryPage, error)); ok {
		return rf(ctx, pm)
	}
	if rf, ok := ret.Get(0).(func(context.Context, callhome.PageMetadata) callhome.TelemetryPage); ok {
		r0 = rf(ctx, pm)
	} else {
		r0 = ret.Get(0).(callhome.TelemetryPage)
	}

	if rf, ok := ret.Get(1).(func(context.Context, callhome.PageMetadata) error); ok {
		r1 = rf(ctx, pm)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

func (s *Service) Save(ctx context.Context, t callhome.Telemetry) error {
	ret := s.Called(ctx, t)
	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, callhome.Telemetry) error); ok {
		r0 = rf(ctx, t)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

func (*Service) RetrieveSummary(ctx context.Context, filters callhome.TelemetryFilters) (callhome.TelemetrySummary, error) {
	return callhome.TelemetrySummary{}, nil
}

type mockConstructorTestingTNewService interface {
	mock.TestingT
	Cleanup(func())
}

func NewService(t mockConstructorTestingTNewService) *Service {
	mock := &Service{}
	mock.Mock.Test(t)
	t.Cleanup(func() { mock.AssertExpectations(t) })
	return mock
}
