// Code generated by mockery v1.0.0. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"
import redigoredis "github.com/garyburd/redigo/redis"

// Pool is an autogenerated mock type for the Pool type
type Pool struct {
	mock.Mock
}

// Close provides a mock function with given fields:
func (_m *Pool) Close() error {
	ret := _m.Called()

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Get provides a mock function with given fields:
func (_m *Pool) Get() redigoredis.Conn {
	ret := _m.Called()

	var r0 redigoredis.Conn
	if rf, ok := ret.Get(0).(func() redigoredis.Conn); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(redigoredis.Conn)
		}
	}

	return r0
}
