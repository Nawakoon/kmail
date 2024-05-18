// Code generated by mockery v2.39.1. DO NOT EDIT.

package mocks

import (
	mail "passwordless-mail-server/pkg/mail"

	mock "github.com/stretchr/testify/mock"

	model "passwordless-mail-server/pkg/model"
)

// MailStore is an autogenerated mock type for the MailStore type
type MailStore struct {
	mock.Mock
}

// GetInbox provides a mock function with given fields: query
func (_m *MailStore) GetInbox(query mail.StoreGetInboxQuery) ([]model.MailEntity, error) {
	ret := _m.Called(query)

	if len(ret) == 0 {
		panic("no return value specified for GetInbox")
	}

	var r0 []model.MailEntity
	var r1 error
	if rf, ok := ret.Get(0).(func(mail.StoreGetInboxQuery) ([]model.MailEntity, error)); ok {
		return rf(query)
	}
	if rf, ok := ret.Get(0).(func(mail.StoreGetInboxQuery) []model.MailEntity); ok {
		r0 = rf(query)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]model.MailEntity)
		}
	}

	if rf, ok := ret.Get(1).(func(mail.StoreGetInboxQuery) error); ok {
		r1 = rf(query)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetMail provides a mock function with given fields:
func (_m *MailStore) GetMail() {
	_m.Called()
}

// InsertMail provides a mock function with given fields: _a0
func (_m *MailStore) InsertMail(_a0 model.Mail) (*model.MailEntity, error) {
	ret := _m.Called(_a0)

	if len(ret) == 0 {
		panic("no return value specified for InsertMail")
	}

	var r0 *model.MailEntity
	var r1 error
	if rf, ok := ret.Get(0).(func(model.Mail) (*model.MailEntity, error)); ok {
		return rf(_a0)
	}
	if rf, ok := ret.Get(0).(func(model.Mail) *model.MailEntity); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.MailEntity)
		}
	}

	if rf, ok := ret.Get(1).(func(model.Mail) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewMailStore creates a new instance of MailStore. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMailStore(t interface {
	mock.TestingT
	Cleanup(func())
}) *MailStore {
	mock := &MailStore{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}