// Code generated by mockery v2.9.4. DO NOT EDIT.

package mocks

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
)

// CouchbaseClusterDiscovery is an autogenerated mock type for the CouchbaseClusterDiscovery type
type CouchbaseClusterDiscovery struct {
	mock.Mock
}

// Discover provides a mock function with given fields: ctx
func (_m *CouchbaseClusterDiscovery) Discover(ctx context.Context) error {
	ret := _m.Called(ctx)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context) error); ok {
		r0 = rf(ctx)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
