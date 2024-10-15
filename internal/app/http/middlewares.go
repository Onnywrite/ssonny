package httpapp

import (
	"net/http"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/rs/zerolog"
)

// logging is a middleware that logs the request and response.
// Probably, should move it to the pkg/*.
func logging(logger zerolog.Logger) fiber.Handler {
	return func(c fiber.Ctx) error {
		start := time.Now()
		err := c.Next()
		end := time.Now()

		logger.Info().
			Str("method", c.Route().Path).
			Str("status", http.StatusText(c.Response().StatusCode())).
			TimeDiff("elapsed", end, start).
			Bool("is_error", err != nil).
			Err(err).
			Msg("http request done")

		return err
	}
}
