// Code generated by mockery v2.31.0. DO NOT EDIT.

package user_mocks

import (
	context "context"

	user "github.com/rubemlrm/go-api-bootstrap/user"
	mock "github.com/stretchr/testify/mock"
)

// MockRepository is an autogenerated mock type for the Repository type
type MockRepository struct {
	mock.Mock
}

type MockRepository_Expecter struct {
	mock *mock.Mock
}

func (_m *MockRepository) EXPECT() *MockRepository_Expecter {
	return &MockRepository_Expecter{mock: &_m.Mock}
}

// All provides a mock function with given fields: ctx
func (_m *MockRepository) All(ctx context.Context) (*[]user.User, error) {
	ret := _m.Called(ctx)

	var r0 *[]user.User
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context) (*[]user.User, error)); ok {
		return rf(ctx)
	}
	if rf, ok := ret.Get(0).(func(context.Context) *[]user.User); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*[]user.User)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockRepository_All_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'All'
type MockRepository_All_Call struct {
	*mock.Call
}

// All is a helper method to define mock.On call
//   - ctx context.Context
func (_e *MockRepository_Expecter) All(ctx interface{}) *MockRepository_All_Call {
	return &MockRepository_All_Call{Call: _e.mock.On("All", ctx)}
}

func (_c *MockRepository_All_Call) Run(run func(ctx context.Context)) *MockRepository_All_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context))
	})
	return _c
}

func (_c *MockRepository_All_Call) Return(_a0 *[]user.User, _a1 error) *MockRepository_All_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockRepository_All_Call) RunAndReturn(run func(context.Context) (*[]user.User, error)) *MockRepository_All_Call {
	_c.Call.Return(run)
	return _c
}

// Create provides a mock function with given fields: ctx, u
func (_m *MockRepository) Create(ctx context.Context, u *user.UserCreate) (user.ID, error) {
	ret := _m.Called(ctx, u)

	var r0 user.ID
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *user.UserCreate) (user.ID, error)); ok {
		return rf(ctx, u)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *user.UserCreate) user.ID); ok {
		r0 = rf(ctx, u)
	} else {
		r0 = ret.Get(0).(user.ID)
	}

	if rf, ok := ret.Get(1).(func(context.Context, *user.UserCreate) error); ok {
		r1 = rf(ctx, u)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockRepository_Create_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Create'
type MockRepository_Create_Call struct {
	*mock.Call
}

// Create is a helper method to define mock.On call
//   - ctx context.Context
//   - u *user.UserCreate
func (_e *MockRepository_Expecter) Create(ctx interface{}, u interface{}) *MockRepository_Create_Call {
	return &MockRepository_Create_Call{Call: _e.mock.On("Create", ctx, u)}
}

func (_c *MockRepository_Create_Call) Run(run func(ctx context.Context, u *user.UserCreate)) *MockRepository_Create_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*user.UserCreate))
	})
	return _c
}

func (_c *MockRepository_Create_Call) Return(_a0 user.ID, _a1 error) *MockRepository_Create_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockRepository_Create_Call) RunAndReturn(run func(context.Context, *user.UserCreate) (user.ID, error)) *MockRepository_Create_Call {
	_c.Call.Return(run)
	return _c
}

// Get provides a mock function with given fields: ctx, id
func (_m *MockRepository) Get(ctx context.Context, id user.ID) (*user.User, error) {
	ret := _m.Called(ctx, id)

	var r0 *user.User
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, user.ID) (*user.User, error)); ok {
		return rf(ctx, id)
	}
	if rf, ok := ret.Get(0).(func(context.Context, user.ID) *user.User); ok {
		r0 = rf(ctx, id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*user.User)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, user.ID) error); ok {
		r1 = rf(ctx, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockRepository_Get_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Get'
type MockRepository_Get_Call struct {
	*mock.Call
}

// Get is a helper method to define mock.On call
//   - ctx context.Context
//   - id user.ID
func (_e *MockRepository_Expecter) Get(ctx interface{}, id interface{}) *MockRepository_Get_Call {
	return &MockRepository_Get_Call{Call: _e.mock.On("Get", ctx, id)}
}

func (_c *MockRepository_Get_Call) Run(run func(ctx context.Context, id user.ID)) *MockRepository_Get_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(user.ID))
	})
	return _c
}

func (_c *MockRepository_Get_Call) Return(_a0 *user.User, _a1 error) *MockRepository_Get_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockRepository_Get_Call) RunAndReturn(run func(context.Context, user.ID) (*user.User, error)) *MockRepository_Get_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockRepository creates a new instance of MockRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockRepository(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockRepository {
	mock := &MockRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
