// Code generated by mockery v2.52.1. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"

// ResponseItemType is an autogenerated mock type for the ResponseItemType type
type ResponseItemType struct {
	mock.Mock
}

// NewResponseItemType creates a new instance of ResponseItemType. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewResponseItemType(t interface {
	mock.TestingT
	Cleanup(func())
}) *ResponseItemType {
	mock := &ResponseItemType{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
