package erix

import (
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

	return msg
}

func GrpcCode(err error) codes.Code {
	if e, ok := err.(*Error); ok { //nolint: errorlint
		return ToGrpc(e.code)
	}

	return codes.Unknown
}

func HttpCode(err error) int {
	if e, ok := err.(*Error); ok { //nolint: errorlint
		return ToHttp(e.code)
	}

	return CodeInternalServerError
}
