package isitjwt

import "errors"

var (
	ErrSecretTooShort      = errors.New("secret is too short, must be at least 32 bytes")
	ErrInvalidToken        = errors.New("invalid token")
	ErrInvalidTokenVersion = errors.New("invalid version, expected 1")
	ErrTokenExpired        = errors.New("token expired")
	ErrInvalidSubject      = errors.New("invalid subject")
)
