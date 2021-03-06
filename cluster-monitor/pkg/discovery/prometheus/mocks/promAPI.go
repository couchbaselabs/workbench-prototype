// Code generated by mockery v2.9.4. DO NOT EDIT.

package mocks

import (
	context "context"

	mock "github.com/stretchr/testify/mock"

	v1 "github.com/prometheus/client_golang/api/prometheus/v1"
)

// PromAPI is an autogenerated mock type for the promAPI type
type PromAPI struct {
	mock.Mock
}

// Targets provides a mock function with given fields: ctx
func (_m *PromAPI) Targets(ctx context.Context) (v1.TargetsResult, error) {
	ret := _m.Called(ctx)

	var r0 v1.TargetsResult
	if rf, ok := ret.Get(0).(func(context.Context) v1.TargetsResult); ok {
		r0 = rf(ctx)
	} else {
		r0 = ret.Get(0).(v1.TargetsResult)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
