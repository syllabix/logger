package mocks

import mock "github.com/stretchr/testify/mock"

// Conn is a mocked implementation of a redigo Conn
type Conn struct {
	mock.Mock
}

// Close provides a mock function with given fields:
func (c *Conn) Close() error {
	ret := c.Called()

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Err provides a mock function with given fields:
func (c *Conn) Err() error {
	ret := c.Called()

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Do provides a mock function with given fields:
func (c *Conn) Do(commandName string, args ...interface{}) (reply interface{}, err error) {
	ret := c.Called(commandName, args)
	return ret.Get(0), ret.Error(1)
}

// Send provides a mock function with given fields:
func (c *Conn) Send(commandName string, args ...interface{}) error {
	ret := c.Called(commandName, args)

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Flush provides a mock function with given fields:
func (c *Conn) Flush() error {
	ret := c.Called()

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Receive provides a mock function with given fields:
func (c *Conn) Receive() (reply interface{}, err error) {
	ret := c.Called()
	return ret.Get(0), ret.Error(1)
}
