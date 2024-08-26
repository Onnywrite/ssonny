package handlersapiauth

import (
	"context"

	"github.com/Onnywrite/ssonny/internal/lib/fiberutil"
	"github.com/Onnywrite/ssonny/internal/lib/tokens"
	"github.com/Onnywrite/ssonny/internal/services/auth"

	"github.com/gofiber/fiber/v3"
)

type Refresher interface {
	Refresh(ctx context.Context, token tokens.Refresh) (*auth.Tokens, error)
}

func Refresh(service Refresher, parser RefreshTokenParser) func(c fiber.Ctx) error {
	type request struct {
		RefreshToken string
	}

	return func(c fiber.Ctx) error {
		var data request
		if err := c.Bind().JSON(&data); err != nil {
			return c.SendStatus(fiber.StatusUnprocessableEntity)
		}

		parserRefresh, err := parser.ParseRefresh(data.RefreshToken)
		if err != nil {
			return c.SendStatus(fiber.StatusUnauthorized)
		}

		tokens, err := service.Refresh(c.Context(), *parserRefresh)
		if err != nil {
			return fiberutil.Error(c, err)
		}

		return c.Status(fiber.StatusOK).JSON(tokens)
	}
}
