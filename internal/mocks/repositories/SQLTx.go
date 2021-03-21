// Code generated by mockery v2.4.0-beta. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"

// SQLTx is an autogenerated mock type for the SQLTx type
type SQLTx struct {
	mock.Mock
}

// Commit provides a mock function with given fields:
func (_m *SQLTx) Commit() error {
	ret := _m.Called()

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Rollback provides a mock function with given fields:
func (_m *SQLTx) Rollback() error {
	ret := _m.Called()

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
