// Code generated by mockery v2.10.4. DO NOT EDIT.

package mocks

import (
	context "context"

	faceitapiclient "github.com/Cludch/csgo-microservices/faceitapiclient/proto"
	mock "github.com/stretchr/testify/mock"
)

// UserHandlerUseCase is an autogenerated mock type for the UserHandlerUseCase type
type UserHandlerUseCase struct {
	mock.Mock
}

// CreateUser provides a mock function with given fields: ctx, req
func (_m *UserHandlerUseCase) CreateUser(ctx context.Context, req *faceitapiclient.CreateUserRequest) (*faceitapiclient.User, error) {
	ret := _m.Called(ctx, req)

	var r0 *faceitapiclient.User
	if rf, ok := ret.Get(0).(func(context.Context, *faceitapiclient.CreateUserRequest) *faceitapiclient.User); ok {
		r0 = rf(ctx, req)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*faceitapiclient.User)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *faceitapiclient.CreateUserRequest) error); ok {
		r1 = rf(ctx, req)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetUser provides a mock function with given fields: ctx, req
func (_m *UserHandlerUseCase) GetUser(ctx context.Context, req *faceitapiclient.GetUserRequest) (*faceitapiclient.User, error) {
	ret := _m.Called(ctx, req)

	var r0 *faceitapiclient.User
	if rf, ok := ret.Get(0).(func(context.Context, *faceitapiclient.GetUserRequest) *faceitapiclient.User); ok {
		r0 = rf(ctx, req)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*faceitapiclient.User)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *faceitapiclient.GetUserRequest) error); ok {
		r1 = rf(ctx, req)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}