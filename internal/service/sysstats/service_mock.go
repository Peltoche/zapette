// Code generated by mockery v2.43.1. DO NOT EDIT.

package sysstats

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
)

// MockService is an autogenerated mock type for the Service type
type MockService struct {
	mock.Mock
}

// GetLatest provides a mock function with given fields: ctx
func (_m *MockService) GetLatest(ctx context.Context) (*Stats, error) {
	ret := _m.Called(ctx)

	if len(ret) == 0 {
		panic("no return value specified for GetLatest")
	}

	var r0 *Stats
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context) (*Stats, error)); ok {
		return rf(ctx)
	}
	if rf, ok := ret.Get(0).(func(context.Context) *Stats); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*Stats)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// fetchAndRegister provides a mock function with given fields: ctx
func (_m *MockService) fetchAndRegister(ctx context.Context) (*Stats, error) {
	ret := _m.Called(ctx)

	if len(ret) == 0 {
		panic("no return value specified for fetchAndRegister")
	}

	var r0 *Stats
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context) (*Stats, error)); ok {
		return rf(ctx)
	}
	if rf, ok := ret.Get(0).(func(context.Context) *Stats); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*Stats)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewMockService creates a new instance of MockService. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockService(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockService {
	mock := &MockService{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
