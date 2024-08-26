package handlersapiauth

import (
	"context"

	"github.com/Onnywrite/ssonny/internal/lib/fiberutil"
	"github.com/Onnywrite/ssonny/internal/services/auth"

	"github.com/gofiber/fiber/v3"
)

type Loginer interface {
	LoginWithPassword(ctx context.Context, data auth.LoginWithPasswordData) (*auth.AuthenticatedUser, error)
}

func LoginWithPassword(service Loginer) func(c fiber.Ctx) error {
	return func(c fiber.Ctx) error {
		var data auth.LoginWithPasswordData
		if err := c.Bind().JSON(&data); err != nil {
			return c.SendStatus(fiber.StatusUnprocessableEntity)
		}

		data.UserInfo = getUserInfo(c)

		authUser, err := service.LoginWithPassword(c.Context(), data)
		if err != nil {
			return fiberutil.Error(c, err)
		}

		return c.Status(fiber.StatusOK).JSON(authUser)
	}
}
