package mocks

import (
	"context"

	"github.com/mainflux/et/internal/homing"
	"github.com/stretchr/testify/mock"
)

var _ homing.Service = (*Service)(nil)

type Service struct {
	mock.Mock
}

func (s *Service) GetAll(ctx context.Context, repo string, token string, pm homing.PageMetadata) (homing.TelemetryPage, error) {
	ret := s.Called(ctx, repo, token, pm)
	return ret.Get(0).(homing.TelemetryPage), ret.Error(1)
}

func (s *Service) Save(ctx context.Context, t homing.Telemetry) error {
	ret := s.Called(ctx, t)
	return ret.Error(1)
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
