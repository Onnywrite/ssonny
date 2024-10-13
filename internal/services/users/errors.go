package users

import "errors"

var (
	ErrUserNotFound       = errors.New("user not found")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrEmailUnverified    = errors.New("email is not verified")
	ErrNicknameInUse      = errors.New("nickname is already in use")
	ErrInternal           = errors.New("internal error")
)
