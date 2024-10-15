package grpcapp

import (
	"context"
	"time"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	"github.com/rs/zerolog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// loggingInterceptor is a middleware that logs the request and response.
func loggingInterceptor(logger zerolog.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (any, error) {
		start := time.Now()
		resp, err := handler(ctx, req)
		end := time.Now()

		code := codes.Unknown
		if s, ok := status.FromError(err); ok {
			code = s.Code()
		}

		logger.Info().
			Str("method", info.FullMethod).
			Stringer("status", code).
			TimeDiff("elapsed", end, start).
			Bool("is_error", err != nil).
			Err(err).
			Msg("grpc request done")

		return resp, nil
	}
}

// recoverInterceptor is a middleware that recovers from panics.
func recoverInterceptor(logger zerolog.Logger) grpc.UnaryServerInterceptor {
	return recovery.UnaryServerInterceptor(
		recovery.WithRecoveryHandler(func(p any) error {
			logger.Error().Any("error", p).Msg("panic was recovered")

			return status.Error(codes.Internal, "TODO")
		}))
}
