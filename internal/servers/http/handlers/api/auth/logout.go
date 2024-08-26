package handlersapiauth

import (
	"context"
	"fmt"

	"github.com/Onnywrite/ssonny/internal/lib/fiberutil"
	"github.com/Onnywrite/ssonny/internal/lib/tokens"

	"github.com/gofiber/fiber/v3"
)

type Logouter interface {
	Logout(ctx context.Context, jwtId uint64) error
}

type RefreshTokenParser interface {
	ParseRefresh(token string) (*tokens.Refresh, error)
}

func Logout(service Logouter, parser RefreshTokenParser) func(c fiber.Ctx) error {
	type request struct {
		RefreshToken string
	}

	return func(c fiber.Ctx) error {
		var data request
		if err := c.Bind().JSON(&data); err != nil {
			return c.SendStatus(fiber.StatusUnprocessableEntity)
		}

		parsedRefresh, err := parser.ParseRefresh(data.RefreshToken)
		if err != nil {
			return c.SendStatus(fiber.StatusUnauthorized)
		}

		access, ok := c.Locals("parsedAccessToken").(*tokens.Access)
		if !ok {
			c.SendStatus(fiber.StatusInternalServerError)
			return fmt.Errorf("leak of access token in Logout handler")
		}

		if access.Subject != parsedRefresh.Subject ||
			access.Audience != parsedRefresh.Audience ||
			access.AuthorizedParty != parsedRefresh.AuthorizedParty {
			return c.SendStatus(fiber.StatusUnauthorized)
		}

		if err := service.Logout(c.Context(), parsedRefresh.Id); err != nil {
			return fiberutil.Error(c, err)
		}

		return c.SendStatus(fiber.StatusOK)
	}
}
