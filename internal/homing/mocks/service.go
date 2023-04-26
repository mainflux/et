package mocks

import (
	"context"

	"github.com/mainflux/callhome/internal/homing"
	"github.com/stretchr/testify/mock"
)

var _ homing.Service = (*Service)(nil)

type Service struct {
	mock.Mock
}

func (s *Service) Retrieve(ctx context.Context, repo string, pm homing.PageMetadata) (homing.TelemetryPage, error) {
	ret := s.Called(ctx, repo, pm)
	var r0 homing.TelemetryPage
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string, homing.PageMetadata) (homing.TelemetryPage, error)); ok {
		return rf(ctx, repo, pm)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, homing.PageMetadata) homing.TelemetryPage); ok {
		r0 = rf(ctx, repo, pm)
	} else {
		r0 = ret.Get(0).(homing.TelemetryPage)
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, homing.PageMetadata) error); ok {
		r1 = rf(ctx, repo, pm)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

func (s *Service) Save(ctx context.Context, t homing.Telemetry) error {
	ret := s.Called(ctx, t)
	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, homing.Telemetry) error); ok {
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
