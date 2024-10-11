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

const allScopes = "*"

func Authorization(parser AccessTokenParser, allowedScopes ...string) func(fiber.Ctx) error {
	return func(c fiber.Ctx) error {
		token := strings.TrimPrefix(c.Get("Authorization"), "Bearer ")
		if token == "" {
			return fiber.NewError(fiber.StatusUnauthorized, "missing 'Authorization: Bearer' header")
		}

		parsedToken, err := parser.ParseAccess(token)
		if err != nil {
			return fiber.NewError(fiber.StatusUnauthorized, "token is expired or invalid")
		}

		if !enoughPermissions(parsedToken.Scopes, allowedScopes) {
			return fiber.NewError(fiber.StatusForbidden, "not enough permissions")
		}

		c.Locals("parsedAccessToken", parsedToken)
		c.Locals("currentUserId", parsedToken.Subject)
		c.Locals("currentAppId", parsedToken.Audience)

		return c.Next()
	}
}

func enoughPermissions(scopes, allowedScopes []string) bool {
	return slices.Contains(scopes, allScopes) ||
		len(allowedScopes) == 0 ||
		slices.ContainsFunc(scopes, func(s string) bool {
			return slices.Contains(allowedScopes, s)
		})
}
