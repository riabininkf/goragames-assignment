// Code generated by mockery v2.42.3. DO NOT EDIT.

package repository

import (
	context "context"

	entity "github.com/riabininkf/goragames-assignment/internal/repository/entity"
	mock "github.com/stretchr/testify/mock"
)

// CodesMock is an autogenerated mock type for the Codes type
type CodesMock struct {
	mock.Mock
}

// DecrementUsagesByName provides a mock function with given fields: ctx, name
func (_m *CodesMock) DecrementUsagesByName(ctx context.Context, name string) error {
	ret := _m.Called(ctx, name)

	if len(ret) == 0 {
		panic("no return value specified for DecrementUsagesByName")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string) error); ok {
		r0 = rf(ctx, name)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetByNameWithLock provides a mock function with given fields: ctx, name
func (_m *CodesMock) GetByNameWithLock(ctx context.Context, name string) (*entity.Code, error) {
	ret := _m.Called(ctx, name)

	if len(ret) == 0 {
		panic("no return value specified for GetByNameWithLock")
	}

	var r0 *entity.Code
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (*entity.Code, error)); ok {
		return rf(ctx, name)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) *entity.Code); ok {
		r0 = rf(ctx, name)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*entity.Code)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, name)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewCodesMock creates a new instance of CodesMock. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewCodesMock(t interface {
	mock.TestingT
	Cleanup(func())
}) *CodesMock {
	mock := &CodesMock{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
