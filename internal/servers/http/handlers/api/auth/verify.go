package handlersapiauth

import (
	"context"

	"github.com/Onnywrite/ssonny/internal/lib/fiberutil"
	"github.com/gofiber/fiber/v3"
)

type Verifier interface {
	VerifyEmail(ctx context.Context, token string) error
}

func VerifyEmail(service Verifier) func(c fiber.Ctx) error {
	return func(c fiber.Ctx) error {
		verificationToken := c.Query("token")
		if verificationToken == "" {
			return c.SendStatus(fiber.StatusBadRequest)
		}

		if err := service.VerifyEmail(c.Context(), verificationToken); err != nil {
			return fiberutil.Error(c, err)
		}

		return c.SendStatus(fiber.StatusOK)
	}
}
