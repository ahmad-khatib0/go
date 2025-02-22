// Code generated by mockery v2.52.1. DO NOT EDIT.

package mocks

import (
	db "github.com/ahmad-khatib0/go/test-driven-development/ch11_generics/db"
	mock "github.com/stretchr/testify/mock"
)

// PostingService is an autogenerated mock type for the PostingService type
type PostingService struct {
	mock.Mock
}

// NewBookOrder provides a mock function with given fields: b
func (_m *PostingService) NewBookOrder(b db.Book) error {
	ret := _m.Called(b)

	if len(ret) == 0 {
		panic("no return value specified for NewBookOrder")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(db.Book) error); ok {
		r0 = rf(b)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewMagazineOrder provides a mock function with given fields: m
func (_m *PostingService) NewMagazineOrder(m db.Magazine) error {
	ret := _m.Called(m)

	if len(ret) == 0 {
		panic("no return value specified for NewMagazineOrder")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(db.Magazine) error); ok {
		r0 = rf(m)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewPostingService creates a new instance of PostingService. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewPostingService(t interface {
	mock.TestingT
	Cleanup(func())
}) *PostingService {
	mock := &PostingService{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
