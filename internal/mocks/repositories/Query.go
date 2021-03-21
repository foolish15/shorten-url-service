// Code generated by mockery v2.4.0-beta. DO NOT EDIT.

package mocks

import (
	repositories "github.com/foolish15/shorten-url-service/internal/repositories"
	mock "github.com/stretchr/testify/mock"
	gorm "gorm.io/gorm"
)

// Query is an autogenerated mock type for the Query type
type Query struct {
	mock.Mock
}

// DB provides a mock function with given fields: _a0
func (_m *Query) DB(_a0 repositories.DB) *gorm.DB {
	ret := _m.Called(_a0)

	var r0 *gorm.DB
	if rf, ok := ret.Get(0).(func(repositories.DB) *gorm.DB); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*gorm.DB)
		}
	}

	return r0
}