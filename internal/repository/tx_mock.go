// Code generated by mockery v2.42.3. DO NOT EDIT.

package repository

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
	mongo "go.mongodb.org/mongo-driver/mongo"
)

// TxMock is an autogenerated mock type for the Tx type
type TxMock struct {
	mock.Mock
}

// Do provides a mock function with given fields: ctx, callback
func (_m *TxMock) Do(ctx context.Context, callback func(mongo.SessionContext) (interface{}, error)) error {
	ret := _m.Called(ctx, callback)

	if len(ret) == 0 {
		panic("no return value specified for Do")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, func(mongo.SessionContext) (interface{}, error)) error); ok {
		r0 = rf(ctx, callback)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewTxMock creates a new instance of TxMock. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewTxMock(t interface {
	mock.TestingT
	Cleanup(func())
}) *TxMock {
	mock := &TxMock{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
