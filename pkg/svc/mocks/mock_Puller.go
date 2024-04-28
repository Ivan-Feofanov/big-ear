// Code generated by mockery v2.42.3. DO NOT EDIT.

package servicemock

import (
	mock "github.com/stretchr/testify/mock"

	stream "github.com/Ivan-Feofanov/big-ear/pkg/stream"
)

// MockPuller is an autogenerated mock type for the Puller type
type MockPuller struct {
	mock.Mock
}

type MockPuller_Expecter struct {
	mock *mock.Mock
}

func (_m *MockPuller) EXPECT() *MockPuller_Expecter {
	return &MockPuller_Expecter{mock: &_m.Mock}
}

// Pull provides a mock function with given fields: _a0
func (_m *MockPuller) Pull(_a0 *stream.Stream) error {
	ret := _m.Called(_a0)

	if len(ret) == 0 {
		panic("no return value specified for Pull")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(*stream.Stream) error); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockPuller_Pull_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Pull'
type MockPuller_Pull_Call struct {
	*mock.Call
}

// Pull is a helper method to define mock.On call
//   - _a0 *stream.Stream
func (_e *MockPuller_Expecter) Pull(_a0 interface{}) *MockPuller_Pull_Call {
	return &MockPuller_Pull_Call{Call: _e.mock.On("Pull", _a0)}
}

func (_c *MockPuller_Pull_Call) Run(run func(_a0 *stream.Stream)) *MockPuller_Pull_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(*stream.Stream))
	})
	return _c
}

func (_c *MockPuller_Pull_Call) Return(_a0 error) *MockPuller_Pull_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockPuller_Pull_Call) RunAndReturn(run func(*stream.Stream) error) *MockPuller_Pull_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockPuller creates a new instance of MockPuller. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockPuller(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockPuller {
	mock := &MockPuller{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
