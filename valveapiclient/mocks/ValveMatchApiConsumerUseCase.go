// Code generated by mockery v2.9.4. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"

// ValveMatchApiConsumerUseCase is an autogenerated mock type for the ValveMatchApiConsumerUseCase type
type ValveMatchApiConsumerUseCase struct {
	mock.Mock
}

// RequestNextShareCode provides a mock function with given fields: steamApiKey, steamId, historyAuthenticationCode, lastShareCode
func (_m *ValveMatchApiConsumerUseCase) RequestNextShareCode(steamApiKey string, steamId uint64, historyAuthenticationCode string, lastShareCode string) (string, error) {
	ret := _m.Called(steamApiKey, steamId, historyAuthenticationCode, lastShareCode)

	var r0 string
	if rf, ok := ret.Get(0).(func(string, uint64, string, string) string); ok {
		r0 = rf(steamApiKey, steamId, historyAuthenticationCode, lastShareCode)
	} else {
		r0 = ret.Get(0).(string)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, uint64, string, string) error); ok {
		r1 = rf(steamApiKey, steamId, historyAuthenticationCode, lastShareCode)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}