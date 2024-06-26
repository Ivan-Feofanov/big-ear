// Code generated by mockery v2.42.3. DO NOT EDIT.

package ethereummock

import (
	context "context"
	big "math/big"

	mock "github.com/stretchr/testify/mock"

	types "github.com/ethereum/go-ethereum/core/types"
)

// MockEthClient is an autogenerated mock type for the EthClient type
type MockEthClient struct {
	mock.Mock
}

type MockEthClient_Expecter struct {
	mock *mock.Mock
}

func (_m *MockEthClient) EXPECT() *MockEthClient_Expecter {
	return &MockEthClient_Expecter{mock: &_m.Mock}
}

// BlockByNumber provides a mock function with given fields: ctx, number
func (_m *MockEthClient) BlockByNumber(ctx context.Context, number *big.Int) (*types.Block, error) {
	ret := _m.Called(ctx, number)

	if len(ret) == 0 {
		panic("no return value specified for BlockByNumber")
	}

	var r0 *types.Block
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *big.Int) (*types.Block, error)); ok {
		return rf(ctx, number)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *big.Int) *types.Block); ok {
		r0 = rf(ctx, number)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*types.Block)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *big.Int) error); ok {
		r1 = rf(ctx, number)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockEthClient_BlockByNumber_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'BlockByNumber'
type MockEthClient_BlockByNumber_Call struct {
	*mock.Call
}

// BlockByNumber is a helper method to define mock.On call
//   - ctx context.Context
//   - number *big.Int
func (_e *MockEthClient_Expecter) BlockByNumber(ctx interface{}, number interface{}) *MockEthClient_BlockByNumber_Call {
	return &MockEthClient_BlockByNumber_Call{Call: _e.mock.On("BlockByNumber", ctx, number)}
}

func (_c *MockEthClient_BlockByNumber_Call) Run(run func(ctx context.Context, number *big.Int)) *MockEthClient_BlockByNumber_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*big.Int))
	})
	return _c
}

func (_c *MockEthClient_BlockByNumber_Call) Return(_a0 *types.Block, _a1 error) *MockEthClient_BlockByNumber_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockEthClient_BlockByNumber_Call) RunAndReturn(run func(context.Context, *big.Int) (*types.Block, error)) *MockEthClient_BlockByNumber_Call {
	_c.Call.Return(run)
	return _c
}

// BlockNumber provides a mock function with given fields: ctx
func (_m *MockEthClient) BlockNumber(ctx context.Context) (uint64, error) {
	ret := _m.Called(ctx)

	if len(ret) == 0 {
		panic("no return value specified for BlockNumber")
	}

	var r0 uint64
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context) (uint64, error)); ok {
		return rf(ctx)
	}
	if rf, ok := ret.Get(0).(func(context.Context) uint64); ok {
		r0 = rf(ctx)
	} else {
		r0 = ret.Get(0).(uint64)
	}

	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockEthClient_BlockNumber_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'BlockNumber'
type MockEthClient_BlockNumber_Call struct {
	*mock.Call
}

// BlockNumber is a helper method to define mock.On call
//   - ctx context.Context
func (_e *MockEthClient_Expecter) BlockNumber(ctx interface{}) *MockEthClient_BlockNumber_Call {
	return &MockEthClient_BlockNumber_Call{Call: _e.mock.On("BlockNumber", ctx)}
}

func (_c *MockEthClient_BlockNumber_Call) Run(run func(ctx context.Context)) *MockEthClient_BlockNumber_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context))
	})
	return _c
}

func (_c *MockEthClient_BlockNumber_Call) Return(_a0 uint64, _a1 error) *MockEthClient_BlockNumber_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockEthClient_BlockNumber_Call) RunAndReturn(run func(context.Context) (uint64, error)) *MockEthClient_BlockNumber_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockEthClient creates a new instance of MockEthClient. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockEthClient(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockEthClient {
	mock := &MockEthClient{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
