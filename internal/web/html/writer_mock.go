// Code generated by mockery v2.43.1. DO NOT EDIT.

package html

import (
	http "net/http"

	mock "github.com/stretchr/testify/mock"
)

// Mock is an autogenerated mock type for the Writer type
type Mock struct {
	mock.Mock
}

// WriteHTMLErrorPage provides a mock function with given fields: w, r, err
func (_m *Mock) WriteHTMLErrorPage(w http.ResponseWriter, r *http.Request, err error) {
	_m.Called(w, r, err)
}

// WriteHTMLTemplate provides a mock function with given fields: w, r, status, template
func (_m *Mock) WriteHTMLTemplate(w http.ResponseWriter, r *http.Request, status int, template Templater) {
	_m.Called(w, r, status, template)
}

// NewMock creates a new instance of Mock. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMock(t interface {
	mock.TestingT
	Cleanup(func())
}) *Mock {
	mock := &Mock{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
