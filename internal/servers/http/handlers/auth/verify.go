package handlersapiauth

import (
	"context"

	"github.com/Onnywrite/ssonny/internal/lib/fiberutil"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

type Verifier interface {
	VerifyEmail(ctx context.Context, userId uuid.UUID) error
}

type EmailTokenParser interface {
	ParseEmail(token string) (uuid.UUID, error)
}

func VerifyEmail(service Verifier, verifier EmailTokenParser) func(c fiber.Ctx) error {
	return func(c fiber.Ctx) error {
		verificationToken := c.Query("token")
		if verificationToken == "" {
			return c.SendStatus(fiber.StatusBadRequest)
		}

		userId, err := verifier.ParseEmail(verificationToken)
		if err != nil {
			return fiberutil.ErrorWithCode(c, err, fiber.StatusBadRequest)
		}

		if err := service.VerifyEmail(c.Context(), userId); err != nil {
			return fiberutil.Error(c, err)
		}

		return c.SendStatus(fiber.StatusOK)
	}
}
