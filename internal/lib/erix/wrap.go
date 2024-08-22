package erix

import (
	"errors"

	"github.com/rotisserie/eris"
	"google.golang.org/grpc/codes"
)

type Error struct {
	erisError error
	code      int
}

func Wrap(err error, code int, stdErr error) error {
	return &Error{
		code:      code,
		erisError: eris.Wrap(err, stdErr.Error()),
	}
}

func (e *Error) Unwrap() error {
	return e.erisError
}

func (e *Error) Error() string {
	var msg string
	if up := eris.Unpack(e.erisError); up.ErrChain != nil {
		msg = up.ErrChain[0].Msg
	} else {
		msg = up.ErrRoot.Msg
	}

	return `{"Service":"ssonny","ErrorMessage":"` + msg + `"}`
}

func GrpcCode(err error) codes.Code {
	var thisError *Error

	if !errors.As(err, &thisError) {
		return codes.Unknown
	}

	return codes.Unknown
}

func HttpCode(err error) int {
	var thisError *Error

	if !errors.As(err, &thisError) {
		return CodeInternalServerError
	}

	return CodeInternalServerError
}
