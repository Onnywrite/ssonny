// Code generated by mockery v2.45.0. DO NOT EDIT.

package authmocks

import (
	context "context"

	email "github.com/Onnywrite/ssonny/internal/services/email"
	mock "github.com/stretchr/testify/mock"
)

// EmailService is an autogenerated mock type for the EmailService type
type EmailService struct {
	mock.Mock
}

type EmailService_Expecter struct {
	mock *mock.Mock
}

func (_m *EmailService) EXPECT() *EmailService_Expecter {
	return &EmailService_Expecter{mock: &_m.Mock}
}

// SendVerificationEmail provides a mock function with given fields: _a0, _a1
func (_m *EmailService) SendVerificationEmail(_a0 context.Context, _a1 email.VerificationEmail) error {
	ret := _m.Called(_a0, _a1)

	if len(ret) == 0 {
		panic("no return value specified for SendVerificationEmail")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, email.VerificationEmail) error); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// EmailService_SendVerificationEmail_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'SendVerificationEmail'
type EmailService_SendVerificationEmail_Call struct {
	*mock.Call
}

// SendVerificationEmail is a helper method to define mock.On call
//   - _a0 context.Context
//   - _a1 email.VerificationEmail
func (_e *EmailService_Expecter) SendVerificationEmail(_a0 interface{}, _a1 interface{}) *EmailService_SendVerificationEmail_Call {
	return &EmailService_SendVerificationEmail_Call{Call: _e.mock.On("SendVerificationEmail", _a0, _a1)}
}

func (_c *EmailService_SendVerificationEmail_Call) Run(run func(_a0 context.Context, _a1 email.VerificationEmail)) *EmailService_SendVerificationEmail_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(email.VerificationEmail))
	})
	return _c
}

func (_c *EmailService_SendVerificationEmail_Call) Return(_a0 error) *EmailService_SendVerificationEmail_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *EmailService_SendVerificationEmail_Call) RunAndReturn(run func(context.Context, email.VerificationEmail) error) *EmailService_SendVerificationEmail_Call {
	_c.Call.Return(run)
	return _c
}

// NewEmailService creates a new instance of EmailService. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewEmailService(t interface {
	mock.TestingT
	Cleanup(func())
}) *EmailService {
	mock := &EmailService{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
