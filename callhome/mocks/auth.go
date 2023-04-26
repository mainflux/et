package mocks

import (
	"context"

	"github.com/mainflux/mainflux"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

var _ mainflux.AuthServiceClient = (*authServiceMock)(nil)

type authServiceMock struct {
	mock.Mock
}

type mockConstructorTestingTNewTelemetryRepo interface {
	mock.TestingT
	Cleanup(func())
}

// NewTelemetryRepo creates a new instance of TelemetryRepo. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewAuthMockRepo(t mockConstructorTestingTNewTelemetryRepo) *authServiceMock {
	mock := &authServiceMock{}
	mock.Mock.Test(t)
	return mock
}

// AddPolicy implements mainflux.AuthServiceClient
func (*authServiceMock) AddPolicy(ctx context.Context, in *mainflux.AddPolicyReq, opts ...grpc.CallOption) (*mainflux.AddPolicyRes, error) {
	panic("unimplemented")
}

// Assign implements mainflux.AuthServiceClient
func (*authServiceMock) Assign(ctx context.Context, in *mainflux.Assignment, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	panic("unimplemented")
}

// Authorize implements mainflux.AuthServiceClient
func (a *authServiceMock) Authorize(ctx context.Context, in *mainflux.AuthorizeReq, opts ...grpc.CallOption) (*mainflux.AuthorizeRes, error) {
	ret := a.Called(ctx, in, opts)
	return ret.Get(0).(*mainflux.AuthorizeRes), ret.Error(1)
}

// DeletePolicy implements mainflux.AuthServiceClient
func (*authServiceMock) DeletePolicy(ctx context.Context, in *mainflux.DeletePolicyReq, opts ...grpc.CallOption) (*mainflux.DeletePolicyRes, error) {
	panic("unimplemented")
}

// Identify implements mainflux.AuthServiceClient
func (a *authServiceMock) Identify(ctx context.Context, in *mainflux.Token, opts ...grpc.CallOption) (*mainflux.UserIdentity, error) {
	ret := a.Called(ctx, in, opts)
	return ret.Get(0).(*mainflux.UserIdentity), ret.Error(1)
}

// Issue implements mainflux.AuthServiceClient
func (*authServiceMock) Issue(ctx context.Context, in *mainflux.IssueReq, opts ...grpc.CallOption) (*mainflux.Token, error) {
	panic("unimplemented")
}

// ListPolicies implements mainflux.AuthServiceClient
func (*authServiceMock) ListPolicies(ctx context.Context, in *mainflux.ListPoliciesReq, opts ...grpc.CallOption) (*mainflux.ListPoliciesRes, error) {
	panic("unimplemented")
}

// Members implements mainflux.AuthServiceClient
func (*authServiceMock) Members(ctx context.Context, in *mainflux.MembersReq, opts ...grpc.CallOption) (*mainflux.MembersRes, error) {
	panic("unimplemented")
}
