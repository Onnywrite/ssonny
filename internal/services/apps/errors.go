package apps

import (
	"errors"
)

var (
	ErrUserUnverified = errors.New("user's account is not verified")

	ErrAppNotFound      = errors.New("app not found")
	ErrAppAlreadyExists = errors.New("app already exists")

	ErrDomainAlreadyExists = errors.New("domain already exists")
	ErrDomainNotFound      = errors.New("domain not found")

	ErrInvalidData        = errors.New("app has invalid data")
	ErrDependencyNotFound = errors.New("data, app depends on, not found")
	ErrInternal           = errors.New("internal error")
)
