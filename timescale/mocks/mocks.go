package mocks

import (
	"context"

	"github.com/mainflux/callhome"
	"github.com/stretchr/testify/mock"
)

var _ callhome.TelemetryRepo = (*mockRepo)(nil)

type mockRepo struct {
	mock.Mock
}

func (mr *mockRepo) RetrieveAll(ctx context.Context, pm callhome.PageMetadata, filter callhome.TelemetryFilters) (callhome.TelemetryPage, error) {
	ret := mr.Called(ctx, pm)
	return ret.Get(0).(callhome.TelemetryPage), ret.Error(1)
}

func (mr *mockRepo) Save(ctx context.Context, t callhome.Telemetry) error {
	ret := mr.Called(ctx, t)
	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, callhome.Telemetry) error); ok {
		r0 = rf(ctx, t)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// RetrieveDistinctIPsCountries retrieve distinct
func (*mockRepo) RetrieveDistinctIPsCountries(ctx context.Context, filter callhome.TelemetryFilters) (callhome.TelemetrySummary, error) {
	return callhome.TelemetrySummary{}, nil
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
