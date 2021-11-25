// Code generated by mockery v2.9.0. DO NOT EDIT.

package mocks

import (
	common "github.com/klaytn/klaytn/common"
	mock "github.com/stretchr/testify/mock"
)

// ETHKeyStore is an autogenerated mock type for the ETHKeyStore type
type ETHKeyStore struct {
	mock.Mock
}

// GetRoundRobinAddress provides a mock function with given fields: addrs
func (_m *ETHKeyStore) GetRoundRobinAddress(addrs ...common.Address) (common.Address, error) {
	_va := make([]interface{}, len(addrs))
	for _i := range addrs {
		_va[_i] = addrs[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 common.Address
	if rf, ok := ret.Get(0).(func(...common.Address) common.Address); ok {
		r0 = rf(addrs...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(common.Address)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(...common.Address) error); ok {
		r1 = rf(addrs...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
