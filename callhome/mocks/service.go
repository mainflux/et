package mocks

import (
	"context"

	"github.com/mainflux/callhome/callhome"
	"github.com/stretchr/testify/mock"
)

var _ callhome.Service = (*Service)(nil)

type Service struct {
	mock.Mock
}

func (s *Service) Retrieve(ctx context.Context, repo string, pm callhome.PageMetadata) (callhome.TelemetryPage, error) {
	ret := s.Called(ctx, repo, pm)
	var r0 callhome.TelemetryPage
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string, callhome.PageMetadata) (callhome.TelemetryPage, error)); ok {
		return rf(ctx, repo, pm)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, callhome.PageMetadata) callhome.TelemetryPage); ok {
		r0 = rf(ctx, repo, pm)
	} else {
		r0 = ret.Get(0).(callhome.TelemetryPage)
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, callhome.PageMetadata) error); ok {
		r1 = rf(ctx, repo, pm)
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
