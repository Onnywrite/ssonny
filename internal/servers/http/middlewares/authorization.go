package middlewares

import (
	"slices"
	"strings"

	"github.com/Onnywrite/ssonny/internal/lib/tokens"

	"github.com/gofiber/fiber/v3"
)

type AccessTokenParser interface {
	ParseAccess(token string) (*tokens.Access, error)
}

func Authorization(parser AccessTokenParser, requiredScopes ...string) func(fiber.Ctx) error {
	return func(c fiber.Ctx) error {
		token := strings.TrimPrefix(c.Get("Authorization"), "Bearer ")
		if token == "" {
			return c.SendStatus(fiber.StatusUnauthorized)
		}

		parsedToken, err := parser.ParseAccess(token)
		if err != nil {
			return c.SendStatus(fiber.StatusUnauthorized)
		}

		if parsedToken.Scopes[0] != "*" {
			for _, requiredScope := range requiredScopes {
				if !slices.Contains(parsedToken.Scopes, requiredScope) {
					return c.SendStatus(fiber.StatusForbidden)
				}
			}
		}

		c.Locals("parsedAccessToken", parsedToken)
		c.Locals("currentUserId", parsedToken.Subject)
		c.Locals("currentAppId", parsedToken.Audience)

		return c.Next()
	}
}
