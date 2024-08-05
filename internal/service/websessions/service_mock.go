// Code generated by mockery v2.43.1. DO NOT EDIT.

package websessions

import (
	context "context"
	http "net/http"

	mock "github.com/stretchr/testify/mock"

	secret "github.com/Peltoche/zapette/internal/tools/secret"

	sqlstorage "github.com/Peltoche/zapette/internal/tools/sqlstorage"

	uuid "github.com/Peltoche/zapette/internal/tools/uuid"
)

// MockService is an autogenerated mock type for the Service type
type MockService struct {
	mock.Mock
}

// Create provides a mock function with given fields: ctx, cmd
func (_m *MockService) Create(ctx context.Context, cmd *CreateCmd) (*Session, error) {
	ret := _m.Called(ctx, cmd)

	if len(ret) == 0 {
		panic("no return value specified for Create")
	}

	var r0 *Session
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *CreateCmd) (*Session, error)); ok {
		return rf(ctx, cmd)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *CreateCmd) *Session); ok {
		r0 = rf(ctx, cmd)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*Session)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *CreateCmd) error); ok {
		r1 = rf(ctx, cmd)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Delete provides a mock function with given fields: ctx, cmd
func (_m *MockService) Delete(ctx context.Context, cmd *DeleteCmd) error {
	ret := _m.Called(ctx, cmd)

	if len(ret) == 0 {
		panic("no return value specified for Delete")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *DeleteCmd) error); ok {
		r0 = rf(ctx, cmd)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// DeleteAll provides a mock function with given fields: ctx, userID
func (_m *MockService) DeleteAll(ctx context.Context, userID uuid.UUID) error {
	ret := _m.Called(ctx, userID)

	if len(ret) == 0 {
		panic("no return value specified for DeleteAll")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID) error); ok {
		r0 = rf(ctx, userID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetAllForUser provides a mock function with given fields: ctx, userID, cmd
func (_m *MockService) GetAllForUser(ctx context.Context, userID uuid.UUID, cmd *sqlstorage.PaginateCmd) ([]Session, error) {
	ret := _m.Called(ctx, userID, cmd)

	if len(ret) == 0 {
		panic("no return value specified for GetAllForUser")
	}

	var r0 []Session
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID, *sqlstorage.PaginateCmd) ([]Session, error)); ok {
		return rf(ctx, userID, cmd)
	}
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID, *sqlstorage.PaginateCmd) []Session); ok {
		r0 = rf(ctx, userID, cmd)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]Session)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, uuid.UUID, *sqlstorage.PaginateCmd) error); ok {
		r1 = rf(ctx, userID, cmd)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetByToken provides a mock function with given fields: ctx, token
func (_m *MockService) GetByToken(ctx context.Context, token secret.Text) (*Session, error) {
	ret := _m.Called(ctx, token)

	if len(ret) == 0 {
		panic("no return value specified for GetByToken")
	}

	var r0 *Session
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, secret.Text) (*Session, error)); ok {
		return rf(ctx, token)
	}
	if rf, ok := ret.Get(0).(func(context.Context, secret.Text) *Session); ok {
		r0 = rf(ctx, token)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*Session)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, secret.Text) error); ok {
		r1 = rf(ctx, token)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetFromReq provides a mock function with given fields: r
func (_m *MockService) GetFromReq(r *http.Request) (*Session, error) {
	ret := _m.Called(r)

	if len(ret) == 0 {
		panic("no return value specified for GetFromReq")
	}

	var r0 *Session
	var r1 error
	if rf, ok := ret.Get(0).(func(*http.Request) (*Session, error)); ok {
		return rf(r)
	}
	if rf, ok := ret.Get(0).(func(*http.Request) *Session); ok {
		r0 = rf(r)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*Session)
		}
	}

	if rf, ok := ret.Get(1).(func(*http.Request) error); ok {
		r1 = rf(r)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Logout provides a mock function with given fields: r, w
func (_m *MockService) Logout(r *http.Request, w http.ResponseWriter) error {
	ret := _m.Called(r, w)

	if len(ret) == 0 {
		panic("no return value specified for Logout")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(*http.Request, http.ResponseWriter) error); ok {
		r0 = rf(r, w)
	} else {
		r0 = ret.Error(0)
	}

	return r0
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
