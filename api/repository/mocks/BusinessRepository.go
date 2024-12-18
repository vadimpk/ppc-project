// Code generated by mockery v2.46.3. DO NOT EDIT.

package mocks

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
	entity "github.com/vadimpk/ppc-project/entity"
)

// BusinessRepository is an autogenerated mock type for the BusinessRepository type
type BusinessRepository struct {
	mock.Mock
}

// Create provides a mock function with given fields: ctx, business
func (_m *BusinessRepository) Create(ctx context.Context, business *entity.Business) error {
	ret := _m.Called(ctx, business)

	if len(ret) == 0 {
		panic("no return value specified for Create")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *entity.Business) error); ok {
		r0 = rf(ctx, business)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Get provides a mock function with given fields: ctx, id
func (_m *BusinessRepository) Get(ctx context.Context, id int) (*entity.Business, error) {
	ret := _m.Called(ctx, id)

	if len(ret) == 0 {
		panic("no return value specified for Get")
	}

	var r0 *entity.Business
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, int) (*entity.Business, error)); ok {
		return rf(ctx, id)
	}
	if rf, ok := ret.Get(0).(func(context.Context, int) *entity.Business); ok {
		r0 = rf(ctx, id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*entity.Business)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, int) error); ok {
		r1 = rf(ctx, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Update provides a mock function with given fields: ctx, business
func (_m *BusinessRepository) Update(ctx context.Context, business *entity.Business) error {
	ret := _m.Called(ctx, business)

	if len(ret) == 0 {
		panic("no return value specified for Update")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *entity.Business) error); ok {
		r0 = rf(ctx, business)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// UpdateAppearance provides a mock function with given fields: ctx, id, logoURL, colorScheme
func (_m *BusinessRepository) UpdateAppearance(ctx context.Context, id int, logoURL string, colorScheme map[string]interface{}) error {
	ret := _m.Called(ctx, id, logoURL, colorScheme)

	if len(ret) == 0 {
		panic("no return value specified for UpdateAppearance")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, int, string, map[string]interface{}) error); ok {
		r0 = rf(ctx, id, logoURL, colorScheme)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewBusinessRepository creates a new instance of BusinessRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewBusinessRepository(t interface {
	mock.TestingT
	Cleanup(func())
}) *BusinessRepository {
	mock := &BusinessRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
