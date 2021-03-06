// Code generated by mockery v0.0.0-dev. DO NOT EDIT.

package mocks

import (
	order "coretrix/internal/order"

	mock "github.com/stretchr/testify/mock"
)

// CartStorage is an autogenerated mock type for the CartStorage type
type CartStorage struct {
	mock.Mock
}

// DeleteCart provides a mock function with given fields: userID
func (_m *CartStorage) DeleteCart(userID int) error {
	ret := _m.Called(userID)

	var r0 error
	if rf, ok := ret.Get(0).(func(int) error); ok {
		r0 = rf(userID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetCart provides a mock function with given fields: userID
func (_m *CartStorage) GetCart(userID int) (order.Cart, error) {
	ret := _m.Called(userID)

	var r0 order.Cart
	if rf, ok := ret.Get(0).(func(int) order.Cart); ok {
		r0 = rf(userID)
	} else {
		r0 = ret.Get(0).(order.Cart)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(int) error); ok {
		r1 = rf(userID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpdateCart provides a mock function with given fields: userID, cart
func (_m *CartStorage) UpdateCart(userID int, cart order.Cart) error {
	ret := _m.Called(userID, cart)

	var r0 error
	if rf, ok := ret.Get(0).(func(int, order.Cart) error); ok {
		r0 = rf(userID, cart)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
