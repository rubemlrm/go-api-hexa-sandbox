// Code generated by mockery v2.31.0. DO NOT EDIT.

package validations_mocks

import (
	validations "github.com/rubemlrm/go-api-bootstrap/pkg/validations"
	mock "github.com/stretchr/testify/mock"
)

// MockOption is an autogenerated mock type for the Option type
type MockOption struct {
	mock.Mock
}

type MockOption_Expecter struct {
	mock *mock.Mock
}

func (_m *MockOption) EXPECT() *MockOption_Expecter {
	return &MockOption_Expecter{mock: &_m.Mock}
}

// Execute provides a mock function with given fields: v
func (_m *MockOption) Execute(v *validations.Validator) error {
	ret := _m.Called(v)

	var r0 error
	if rf, ok := ret.Get(0).(func(*validations.Validator) error); ok {
		r0 = rf(v)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockOption_Execute_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Execute'
type MockOption_Execute_Call struct {
	*mock.Call
}

// Execute is a helper method to define mock.On call
//   - v *validations.Validator
func (_e *MockOption_Expecter) Execute(v interface{}) *MockOption_Execute_Call {
	return &MockOption_Execute_Call{Call: _e.mock.On("Execute", v)}
}

func (_c *MockOption_Execute_Call) Run(run func(v *validations.Validator)) *MockOption_Execute_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(*validations.Validator))
	})
	return _c
}

func (_c *MockOption_Execute_Call) Return(_a0 error) *MockOption_Execute_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockOption_Execute_Call) RunAndReturn(run func(*validations.Validator) error) *MockOption_Execute_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockOption creates a new instance of MockOption. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockOption(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockOption {
	mock := &MockOption{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}