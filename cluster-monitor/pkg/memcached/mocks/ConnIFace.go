// Code generated by mockery v2.9.4. DO NOT EDIT.

package mocks

import (
	memcached "github.com/couchbaselabs/workbench-prototype/cluster-monitor/pkg/memcached"
	mock "github.com/stretchr/testify/mock"

	values "github.com/couchbaselabs/workbench-prototype/cluster-monitor/pkg/values"
)

// ConnIFace is an autogenerated mock type for the ConnIFace type
type ConnIFace struct {
	mock.Mock
}

// CheckpointStats provides a mock function with given fields: host, bucket
func (_m *ConnIFace) CheckpointStats(host string, bucket string) (memcached.BucketCheckpointStats, error) {
	ret := _m.Called(host, bucket)

	var r0 memcached.BucketCheckpointStats
	if rf, ok := ret.Get(0).(func(string, string) memcached.BucketCheckpointStats); ok {
		r0 = rf(host, bucket)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(memcached.BucketCheckpointStats)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, string) error); ok {
		r1 = rf(host, bucket)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Close provides a mock function with given fields:
func (_m *ConnIFace) Close() error {
	ret := _m.Called()

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// DCPStats provides a mock function with given fields: bucket
func (_m *ConnIFace) DCPStats(bucket string) ([]*memcached.DCPMemStats, error) {
	ret := _m.Called(bucket)

	var r0 []*memcached.DCPMemStats
	if rf, ok := ret.Get(0).(func(string) []*memcached.DCPMemStats); ok {
		r0 = rf(bucket)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*memcached.DCPMemStats)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(bucket)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DefaultStats provides a mock function with given fields: bucket
func (_m *ConnIFace) DefaultStats(bucket string) ([]*memcached.DefStats, error) {
	ret := _m.Called(bucket)

	var r0 []*memcached.DefStats
	if rf, ok := ret.Get(0).(func(string) []*memcached.DefStats); ok {
		r0 = rf(bucket)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*memcached.DefStats)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(bucket)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Hosts provides a mock function with given fields:
func (_m *ConnIFace) Hosts() []string {
	ret := _m.Called()

	var r0 []string
	if rf, ok := ret.Get(0).(func() []string); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]string)
		}
	}

	return r0
}

// MemStats provides a mock function with given fields: bucket
func (_m *ConnIFace) MemStats(bucket string) ([]*memcached.MemoryStats, error) {
	ret := _m.Called(bucket)

	var r0 []*memcached.MemoryStats
	if rf, ok := ret.Get(0).(func(string) []*memcached.MemoryStats); ok {
		r0 = rf(bucket)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*memcached.MemoryStats)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(bucket)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
