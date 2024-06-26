// Code generated by mockery v2.42.1. DO NOT EDIT.

package mocks

import (
	domain "github.com/peppydays/grpc-in-go-service/microservices/order/internal/domain"
	mock "github.com/stretchr/testify/mock"
)

// PaymentGateway is an autogenerated mock type for the PaymentGateway type
type PaymentGateway struct {
	mock.Mock
}

// Charge provides a mock function with given fields: _a0
func (_m *PaymentGateway) Charge(_a0 *domain.Order) error {
	ret := _m.Called(_a0)

	if len(ret) == 0 {
		panic("no return value specified for Charge")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(*domain.Order) error); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewPaymentGateway creates a new instance of PaymentGateway. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewPaymentGateway(t interface {
	mock.TestingT
	Cleanup(func())
}) *PaymentGateway {
	mock := &PaymentGateway{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
