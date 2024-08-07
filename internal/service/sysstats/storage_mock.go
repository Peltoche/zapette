// Code generated by mockery v2.43.2. DO NOT EDIT.

package sysstats

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
)

// mockStorage is an autogenerated mock type for the storage type
type mockStorage struct {
	mock.Mock
}

// GetLatest provides a mock function with given fields: ctx
func (_m *mockStorage) GetLatest(ctx context.Context) (*Stats, error) {
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

// Save provides a mock function with given fields: ctx, stats
func (_m *mockStorage) Save(ctx context.Context, stats *Stats) error {
	ret := _m.Called(ctx, stats)

	if len(ret) == 0 {
		panic("no return value specified for Save")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *Stats) error); ok {
		r0 = rf(ctx, stats)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// newMockStorage creates a new instance of mockStorage. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func newMockStorage(t interface {
	mock.TestingT
	Cleanup(func())
}) *mockStorage {
	mock := &mockStorage{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
