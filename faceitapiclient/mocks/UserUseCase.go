// Code generated by mockery v2.10.4. DO NOT EDIT.

package mocks

import (
	user "github.com/Cludch/csgo-microservices/faceitapiclient/internal/domain/user"
	uuid "github.com/google/uuid"
	mock "github.com/stretchr/testify/mock"
)

// UserUseCase is an autogenerated mock type for the UserUseCase type
type UserUseCase struct {
	mock.Mock
}

// CreateUser provides a mock function with given fields: id
func (_m *UserUseCase) CreateUser(id uuid.UUID) (*user.User, error) {
	ret := _m.Called(id)

	var r0 *user.User
	if rf, ok := ret.Get(0).(func(uuid.UUID) *user.User); ok {
		r0 = rf(id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*user.User)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(uuid.UUID) error); ok {
		r1 = rf(id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetUser provides a mock function with given fields: _a0
func (_m *UserUseCase) GetUser(_a0 uuid.UUID) (*user.User, error) {
	ret := _m.Called(_a0)

	var r0 *user.User
	if rf, ok := ret.Get(0).(func(uuid.UUID) *user.User); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*user.User)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(uuid.UUID) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetUsersWithApiEnabled provides a mock function with given fields:
func (_m *UserUseCase) GetUsersWithApiEnabled() ([]*user.User, error) {
	ret := _m.Called()

	var r0 []*user.User
	if rf, ok := ret.Get(0).(func() []*user.User); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*user.User)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
