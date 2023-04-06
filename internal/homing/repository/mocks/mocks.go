package mocks

import (
	"context"

	"github.com/mainflux/et/internal/homing"
	"github.com/stretchr/testify/mock"
)

var _ homing.TelemetryRepo = (*mockRepo)(nil)

type mockRepo struct {
	mock.Mock
}

func (mr *mockRepo) RetrieveAll(ctx context.Context, pm homing.PageMetadata) (homing.TelemetryPage, error) {
	ret := mr.Called(ctx, pm)
	return ret.Get(0).(homing.TelemetryPage), ret.Error(1)
}

func (mr *mockRepo) RetrieveByIP(ctx context.Context, email string) (homing.Telemetry, error) {
	ret := mr.Called(ctx, email)
	return ret.Get(0).(homing.Telemetry), ret.Error(1)
}

func (mr *mockRepo) Save(ctx context.Context, t homing.Telemetry) error {
	ret := mr.Called(ctx, t)
	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, homing.Telemetry) error); ok {
		r0 = rf(ctx, t)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

func (mr *mockRepo) UpdateTelemetry(ctx context.Context, u homing.Telemetry) error {
	ret := mr.Called(ctx, u)
	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, homing.Telemetry) error); ok {
		r0 = rf(ctx, u)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

type mockConstructorTestingTNewTelemetryRepo interface {
	mock.TestingT
	Cleanup(func())
}

func NewTelemetryRepo(t mockConstructorTestingTNewTelemetryRepo) *mockRepo {
	mock := &mockRepo{}
	mock.Mock.Test(t)

	return mock
}
