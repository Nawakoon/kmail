// Code generated by mockery v2.39.1. DO NOT EDIT.

package mocks

import (
	model "passwordless-mail-server/pkg/model"

	mock "github.com/stretchr/testify/mock"

	uuid "github.com/google/uuid"
)

// UuidStore is an autogenerated mock type for the UuidStore type
type UuidStore struct {
	mock.Mock
}

// GetUsedUUID provides a mock function with given fields: _a0
func (_m *UuidStore) GetUsedUUID(_a0 uuid.UUID) (*model.UsedUUIDEntity, error) {
	ret := _m.Called(_a0)

	if len(ret) == 0 {
		panic("no return value specified for GetUsedUUID")
	}

	var r0 *model.UsedUUIDEntity
	var r1 error
	if rf, ok := ret.Get(0).(func(uuid.UUID) (*model.UsedUUIDEntity, error)); ok {
		return rf(_a0)
	}
	if rf, ok := ret.Get(0).(func(uuid.UUID) *model.UsedUUIDEntity); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.UsedUUIDEntity)
		}
	}

	if rf, ok := ret.Get(1).(func(uuid.UUID) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// InsertUsedUUID provides a mock function with given fields: _a0
func (_m *UuidStore) InsertUsedUUID(_a0 uuid.UUID) error {
	ret := _m.Called(_a0)

	if len(ret) == 0 {
		panic("no return value specified for InsertUsedUUID")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(uuid.UUID) error); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewUuidStore creates a new instance of UuidStore. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewUuidStore(t interface {
	mock.TestingT
	Cleanup(func())
}) *UuidStore {
	mock := &UuidStore{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}