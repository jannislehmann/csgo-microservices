// Code generated by mockery v2.10.0. DO NOT EDIT.

package mocks

import (
	player "github.com/Cludch/csgo-microservices/demoparser/internal/domain/player"
	uuid "github.com/google/uuid"
	mock "github.com/stretchr/testify/mock"
)

// UseCase is an autogenerated mock type for the UseCase type
type UseCase struct {
	mock.Mock
}

// AddResult provides a mock function with given fields: _a0, _a1
func (_m *UseCase) AddResult(_a0 *player.Player, _a1 *player.PlayerResult) error {
	ret := _m.Called(_a0, _a1)

	var r0 error
	if rf, ok := ret.Get(0).(func(*player.Player, *player.PlayerResult) error); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// CreatePlayer provides a mock function with given fields: steamId
func (_m *UseCase) CreatePlayer(steamId uint64) (*player.Player, error) {
	ret := _m.Called(steamId)

	var r0 *player.Player
	if rf, ok := ret.Get(0).(func(uint64) *player.Player); ok {
		r0 = rf(steamId)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*player.Player)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(uint64) error); ok {
		r1 = rf(steamId)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DeleteResult provides a mock function with given fields: p, matchId
func (_m *UseCase) DeleteResult(p *player.Player, matchId uuid.UUID) error {
	ret := _m.Called(p, matchId)

	var r0 error
	if rf, ok := ret.Get(0).(func(*player.Player, uuid.UUID) error); ok {
		r0 = rf(p, matchId)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetAll provides a mock function with given fields:
func (_m *UseCase) GetAll() ([]*player.Player, error) {
	ret := _m.Called()

	var r0 []*player.Player
	if rf, ok := ret.Get(0).(func() []*player.Player); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*player.Player)
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

// GetPlayer provides a mock function with given fields: _a0
func (_m *UseCase) GetPlayer(_a0 uint64) (*player.Player, error) {
	ret := _m.Called(_a0)

	var r0 *player.Player
	if rf, ok := ret.Get(0).(func(uint64) *player.Player); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*player.Player)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(uint64) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetResult provides a mock function with given fields: p, matchId
func (_m *UseCase) GetResult(p *player.Player, matchId uuid.UUID) (*player.PlayerResult, error) {
	ret := _m.Called(p, matchId)

	var r0 *player.PlayerResult
	if rf, ok := ret.Get(0).(func(*player.Player, uuid.UUID) *player.PlayerResult); ok {
		r0 = rf(p, matchId)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*player.PlayerResult)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*player.Player, uuid.UUID) error); ok {
		r1 = rf(p, matchId)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}