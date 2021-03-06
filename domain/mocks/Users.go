// Code generated by mockery v1.0.0. DO NOT EDIT.

package mocks

import (
	domain "github.com/concepts-system/go-paperless/domain"
	mock "github.com/stretchr/testify/mock"
)

// Users is an autogenerated mock type for the Users type
type Users struct {
	mock.Mock
}

// Add provides a mock function with given fields: user
func (_m *Users) Add(user *domain.User) (*domain.User, error) {
	ret := _m.Called(user)

	var r0 *domain.User
	if rf, ok := ret.Get(0).(func(*domain.User) *domain.User); ok {
		r0 = rf(user)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*domain.User)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*domain.User) error); ok {
		r1 = rf(user)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Delete provides a mock function with given fields: user
func (_m *Users) Delete(user *domain.User) error {
	ret := _m.Called(user)

	var r0 error
	if rf, ok := ret.Get(0).(func(*domain.User) error); ok {
		r0 = rf(user)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Find provides a mock function with given fields: page
func (_m *Users) Find(page domain.PageRequest) ([]domain.User, domain.Count, error) {
	ret := _m.Called(page)

	var r0 []domain.User
	if rf, ok := ret.Get(0).(func(domain.PageRequest) []domain.User); ok {
		r0 = rf(page)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]domain.User)
		}
	}

	var r1 domain.Count
	if rf, ok := ret.Get(1).(func(domain.PageRequest) domain.Count); ok {
		r1 = rf(page)
	} else {
		r1 = ret.Get(1).(domain.Count)
	}

	var r2 error
	if rf, ok := ret.Get(2).(func(domain.PageRequest) error); ok {
		r2 = rf(page)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// GetByUsername provides a mock function with given fields: username
func (_m *Users) GetByUsername(username domain.Name) (*domain.User, error) {
	ret := _m.Called(username)

	var r0 *domain.User
	if rf, ok := ret.Get(0).(func(domain.Name) *domain.User); ok {
		r0 = rf(username)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*domain.User)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(domain.Name) error); ok {
		r1 = rf(username)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Update provides a mock function with given fields: user
func (_m *Users) Update(user *domain.User) (*domain.User, error) {
	ret := _m.Called(user)

	var r0 *domain.User
	if rf, ok := ret.Get(0).(func(*domain.User) *domain.User); ok {
		r0 = rf(user)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*domain.User)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*domain.User) error); ok {
		r1 = rf(user)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
