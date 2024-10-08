// Code generated by mockery v2.45.0. DO NOT EDIT.

package authmocks

import (
	context "context"

	models "github.com/Onnywrite/ssonny/internal/domain/models"
	mock "github.com/stretchr/testify/mock"

	repo "github.com/Onnywrite/ssonny/internal/storage/repo"

	uuid "github.com/google/uuid"
)

// UserRepo is an autogenerated mock type for the UserRepo type
type UserRepo struct {
	mock.Mock
}

type UserRepo_Expecter struct {
	mock *mock.Mock
}

func (_m *UserRepo) EXPECT() *UserRepo_Expecter {
	return &UserRepo_Expecter{mock: &_m.Mock}
}

// SaveUser provides a mock function with given fields: _a0, _a1
func (_m *UserRepo) SaveUser(_a0 context.Context, _a1 models.User) (*models.User, repo.Transactor, error) {
	ret := _m.Called(_a0, _a1)

	if len(ret) == 0 {
		panic("no return value specified for SaveUser")
	}

	var r0 *models.User
	var r1 repo.Transactor
	var r2 error
	if rf, ok := ret.Get(0).(func(context.Context, models.User) (*models.User, repo.Transactor, error)); ok {
		return rf(_a0, _a1)
	}
	if rf, ok := ret.Get(0).(func(context.Context, models.User) *models.User); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.User)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, models.User) repo.Transactor); ok {
		r1 = rf(_a0, _a1)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(repo.Transactor)
		}
	}

	if rf, ok := ret.Get(2).(func(context.Context, models.User) error); ok {
		r2 = rf(_a0, _a1)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// UserRepo_SaveUser_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'SaveUser'
type UserRepo_SaveUser_Call struct {
	*mock.Call
}

// SaveUser is a helper method to define mock.On call
//   - _a0 context.Context
//   - _a1 models.User
func (_e *UserRepo_Expecter) SaveUser(_a0 interface{}, _a1 interface{}) *UserRepo_SaveUser_Call {
	return &UserRepo_SaveUser_Call{Call: _e.mock.On("SaveUser", _a0, _a1)}
}

func (_c *UserRepo_SaveUser_Call) Run(run func(_a0 context.Context, _a1 models.User)) *UserRepo_SaveUser_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(models.User))
	})
	return _c
}

func (_c *UserRepo_SaveUser_Call) Return(_a0 *models.User, _a1 repo.Transactor, _a2 error) *UserRepo_SaveUser_Call {
	_c.Call.Return(_a0, _a1, _a2)
	return _c
}

func (_c *UserRepo_SaveUser_Call) RunAndReturn(run func(context.Context, models.User) (*models.User, repo.Transactor, error)) *UserRepo_SaveUser_Call {
	_c.Call.Return(run)
	return _c
}

// UpdateUser provides a mock function with given fields: ctx, userId, newValues
func (_m *UserRepo) UpdateUser(ctx context.Context, userId uuid.UUID, newValues map[string]interface{}) error {
	ret := _m.Called(ctx, userId, newValues)

	if len(ret) == 0 {
		panic("no return value specified for UpdateUser")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID, map[string]interface{}) error); ok {
		r0 = rf(ctx, userId, newValues)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// UserRepo_UpdateUser_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'UpdateUser'
type UserRepo_UpdateUser_Call struct {
	*mock.Call
}

// UpdateUser is a helper method to define mock.On call
//   - ctx context.Context
//   - userId uuid.UUID
//   - newValues map[string]interface{}
func (_e *UserRepo_Expecter) UpdateUser(ctx interface{}, userId interface{}, newValues interface{}) *UserRepo_UpdateUser_Call {
	return &UserRepo_UpdateUser_Call{Call: _e.mock.On("UpdateUser", ctx, userId, newValues)}
}

func (_c *UserRepo_UpdateUser_Call) Run(run func(ctx context.Context, userId uuid.UUID, newValues map[string]interface{})) *UserRepo_UpdateUser_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(uuid.UUID), args[2].(map[string]interface{}))
	})
	return _c
}

func (_c *UserRepo_UpdateUser_Call) Return(_a0 error) *UserRepo_UpdateUser_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *UserRepo_UpdateUser_Call) RunAndReturn(run func(context.Context, uuid.UUID, map[string]interface{}) error) *UserRepo_UpdateUser_Call {
	_c.Call.Return(run)
	return _c
}

// UserByEmail provides a mock function with given fields: _a0, _a1
func (_m *UserRepo) UserByEmail(_a0 context.Context, _a1 string) (*models.User, error) {
	ret := _m.Called(_a0, _a1)

	if len(ret) == 0 {
		panic("no return value specified for UserByEmail")
	}

	var r0 *models.User
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (*models.User, error)); ok {
		return rf(_a0, _a1)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) *models.User); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.User)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UserRepo_UserByEmail_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'UserByEmail'
type UserRepo_UserByEmail_Call struct {
	*mock.Call
}

// UserByEmail is a helper method to define mock.On call
//   - _a0 context.Context
//   - _a1 string
func (_e *UserRepo_Expecter) UserByEmail(_a0 interface{}, _a1 interface{}) *UserRepo_UserByEmail_Call {
	return &UserRepo_UserByEmail_Call{Call: _e.mock.On("UserByEmail", _a0, _a1)}
}

func (_c *UserRepo_UserByEmail_Call) Run(run func(_a0 context.Context, _a1 string)) *UserRepo_UserByEmail_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string))
	})
	return _c
}

func (_c *UserRepo_UserByEmail_Call) Return(_a0 *models.User, _a1 error) *UserRepo_UserByEmail_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *UserRepo_UserByEmail_Call) RunAndReturn(run func(context.Context, string) (*models.User, error)) *UserRepo_UserByEmail_Call {
	_c.Call.Return(run)
	return _c
}

// UserById provides a mock function with given fields: _a0, _a1
func (_m *UserRepo) UserById(_a0 context.Context, _a1 uuid.UUID) (*models.User, error) {
	ret := _m.Called(_a0, _a1)

	if len(ret) == 0 {
		panic("no return value specified for UserById")
	}

	var r0 *models.User
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID) (*models.User, error)); ok {
		return rf(_a0, _a1)
	}
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID) *models.User); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.User)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, uuid.UUID) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UserRepo_UserById_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'UserById'
type UserRepo_UserById_Call struct {
	*mock.Call
}

// UserById is a helper method to define mock.On call
//   - _a0 context.Context
//   - _a1 uuid.UUID
func (_e *UserRepo_Expecter) UserById(_a0 interface{}, _a1 interface{}) *UserRepo_UserById_Call {
	return &UserRepo_UserById_Call{Call: _e.mock.On("UserById", _a0, _a1)}
}

func (_c *UserRepo_UserById_Call) Run(run func(_a0 context.Context, _a1 uuid.UUID)) *UserRepo_UserById_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(uuid.UUID))
	})
	return _c
}

func (_c *UserRepo_UserById_Call) Return(_a0 *models.User, _a1 error) *UserRepo_UserById_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *UserRepo_UserById_Call) RunAndReturn(run func(context.Context, uuid.UUID) (*models.User, error)) *UserRepo_UserById_Call {
	_c.Call.Return(run)
	return _c
}

// UserByNickname provides a mock function with given fields: _a0, _a1
func (_m *UserRepo) UserByNickname(_a0 context.Context, _a1 string) (*models.User, error) {
	ret := _m.Called(_a0, _a1)

	if len(ret) == 0 {
		panic("no return value specified for UserByNickname")
	}

	var r0 *models.User
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (*models.User, error)); ok {
		return rf(_a0, _a1)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) *models.User); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.User)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UserRepo_UserByNickname_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'UserByNickname'
type UserRepo_UserByNickname_Call struct {
	*mock.Call
}

// UserByNickname is a helper method to define mock.On call
//   - _a0 context.Context
//   - _a1 string
func (_e *UserRepo_Expecter) UserByNickname(_a0 interface{}, _a1 interface{}) *UserRepo_UserByNickname_Call {
	return &UserRepo_UserByNickname_Call{Call: _e.mock.On("UserByNickname", _a0, _a1)}
}

func (_c *UserRepo_UserByNickname_Call) Run(run func(_a0 context.Context, _a1 string)) *UserRepo_UserByNickname_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string))
	})
	return _c
}

func (_c *UserRepo_UserByNickname_Call) Return(_a0 *models.User, _a1 error) *UserRepo_UserByNickname_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *UserRepo_UserByNickname_Call) RunAndReturn(run func(context.Context, string) (*models.User, error)) *UserRepo_UserByNickname_Call {
	_c.Call.Return(run)
	return _c
}

// NewUserRepo creates a new instance of UserRepo. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewUserRepo(t interface {
	mock.TestingT
	Cleanup(func())
}) *UserRepo {
	mock := &UserRepo{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
