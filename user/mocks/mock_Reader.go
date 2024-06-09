// Code generated by mockery v2.43.2. DO NOT EDIT.

package user_mocks

import (
	user "github.com/rubemlrm/go-api-bootstrap/user"
	mock "github.com/stretchr/testify/mock"
)

// MockReader is an autogenerated mock type for the Reader type
type MockReader struct {
	mock.Mock
}

type MockReader_Expecter struct {
	mock *mock.Mock
}

func (_m *MockReader) EXPECT() *MockReader_Expecter {
	return &MockReader_Expecter{mock: &_m.Mock}
}

// All provides a mock function with given fields:
func (_m *MockReader) All() (*[]user.User, error) {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for All")
	}

	var r0 *[]user.User
	var r1 error
	if rf, ok := ret.Get(0).(func() (*[]user.User, error)); ok {
		return rf()
	}
	if rf, ok := ret.Get(0).(func() *[]user.User); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*[]user.User)
		}
	}

	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockReader_All_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'All'
type MockReader_All_Call struct {
	*mock.Call
}

// All is a helper method to define mock.On call
func (_e *MockReader_Expecter) All() *MockReader_All_Call {
	return &MockReader_All_Call{Call: _e.mock.On("All")}
}

func (_c *MockReader_All_Call) Run(run func()) *MockReader_All_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockReader_All_Call) Return(_a0 *[]user.User, _a1 error) *MockReader_All_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockReader_All_Call) RunAndReturn(run func() (*[]user.User, error)) *MockReader_All_Call {
	_c.Call.Return(run)
	return _c
}

// Get provides a mock function with given fields: id
func (_m *MockReader) Get(id user.ID) (*user.User, error) {
	ret := _m.Called(id)

	if len(ret) == 0 {
		panic("no return value specified for Get")
	}

	var r0 *user.User
	var r1 error
	if rf, ok := ret.Get(0).(func(user.ID) (*user.User, error)); ok {
		return rf(id)
	}
	if rf, ok := ret.Get(0).(func(user.ID) *user.User); ok {
		r0 = rf(id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*user.User)
		}
	}

	if rf, ok := ret.Get(1).(func(user.ID) error); ok {
		r1 = rf(id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockReader_Get_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Get'
type MockReader_Get_Call struct {
	*mock.Call
}

// Get is a helper method to define mock.On call
//   - id user.ID
func (_e *MockReader_Expecter) Get(id interface{}) *MockReader_Get_Call {
	return &MockReader_Get_Call{Call: _e.mock.On("Get", id)}
}

func (_c *MockReader_Get_Call) Run(run func(id user.ID)) *MockReader_Get_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(user.ID))
	})
	return _c
}

func (_c *MockReader_Get_Call) Return(_a0 *user.User, _a1 error) *MockReader_Get_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockReader_Get_Call) RunAndReturn(run func(user.ID) (*user.User, error)) *MockReader_Get_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockReader creates a new instance of MockReader. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockReader(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockReader {
	mock := &MockReader{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
