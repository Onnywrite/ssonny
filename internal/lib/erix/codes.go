package erix

import (
	"google.golang.org/grpc/codes"
)

const (
	CodeBadRequest         = 400
	CodeUnauthorized       = 401
	CodeForbidden          = 403
	CodeNotFound           = 404
	CodeRequestTimeout     = 408
	CodeConflict           = 409
	CodePreconditionFailed = 412
	CodeTooManyRequests    = 429

	CodeInternalServerError = 500
	CodeNotImplemented      = 501
	CodeServiceUnavailable  = 503
	CodeGatewayTimeout      = 504
)

func ToHttp(code int) int {
	return code
}

//nolint: cyclop
func ToGrpc(code int) codes.Code {
	switch code {
	case CodeRequestTimeout:
		return codes.Canceled
	case CodeBadRequest:
		return codes.InvalidArgument
	case CodeGatewayTimeout:
		return codes.DeadlineExceeded
	case CodeNotFound:
		return codes.NotFound
	case CodeConflict:
		return codes.AlreadyExists
	case CodeForbidden:
		return codes.PermissionDenied
	case CodeTooManyRequests:
		return codes.ResourceExhausted
	case CodePreconditionFailed:
		return codes.FailedPrecondition
	case CodeNotImplemented:
		return codes.Unimplemented
	case CodeInternalServerError:
		return codes.Internal
	case CodeServiceUnavailable:
		return codes.Unavailable
	case CodeUnauthorized:
		return codes.Unauthenticated
	}

	return codes.Unknown
}
