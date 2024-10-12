package users

import "errors"

var (
	ErrUserNotFound       = errors.New("user not found")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrEmailUnverified    = errors.New("email is not verified")
	ErrInternal           = errors.New("internal error")
)
