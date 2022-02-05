// Code generated by mockery v2.9.4. DO NOT EDIT.

package mocks

import (
	context "context"

	mock "github.com/stretchr/testify/mock"

	valveapiclient "github.com/Cludch/csgo-microservices/valveapiclient/proto"
)

// UserHandlerUseCase is an autogenerated mock type for the UserHandlerUseCase type
type UserHandlerUseCase struct {
	mock.Mock
}

// CreateUser provides a mock function with given fields: ctx, req
func (_m *UserHandlerUseCase) CreateUser(ctx context.Context, req *valveapiclient.CreateUserRequest) (*valveapiclient.User, error) {
	ret := _m.Called(ctx, req)

	var r0 *valveapiclient.User
	if rf, ok := ret.Get(0).(func(context.Context, *valveapiclient.CreateUserRequest) *valveapiclient.User); ok {
		r0 = rf(ctx, req)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*valveapiclient.User)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *valveapiclient.CreateUserRequest) error); ok {
		r1 = rf(ctx, req)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetUser provides a mock function with given fields: ctx, req
func (_m *UserHandlerUseCase) GetUser(ctx context.Context, req *valveapiclient.GetUserRequest) (*valveapiclient.User, error) {
	ret := _m.Called(ctx, req)

	var r0 *valveapiclient.User
	if rf, ok := ret.Get(0).(func(context.Context, *valveapiclient.GetUserRequest) *valveapiclient.User); ok {
		r0 = rf(ctx, req)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*valveapiclient.User)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *valveapiclient.GetUserRequest) error); ok {
		r1 = rf(ctx, req)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpdateUserApiCredentials provides a mock function with given fields: ctx, req
func (_m *UserHandlerUseCase) UpdateUserApiCredentials(ctx context.Context, req *valveapiclient.UpdateUserApiCredentialsRequest) (*valveapiclient.StatusResponse, error) {
	ret := _m.Called(ctx, req)

	var r0 *valveapiclient.StatusResponse
	if rf, ok := ret.Get(0).(func(context.Context, *valveapiclient.UpdateUserApiCredentialsRequest) *valveapiclient.StatusResponse); ok {
		r0 = rf(ctx, req)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*valveapiclient.StatusResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *valveapiclient.UpdateUserApiCredentialsRequest) error); ok {
		r1 = rf(ctx, req)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}