package fiberutil

import (
	"github.com/Onnywrite/ssonny/internal/lib/erix"

	"github.com/gofiber/fiber/v3"
)

func Error(c fiber.Ctx, err error) error {
	c.Response().SetBodyString(err.Error())
	c.Response().Header.SetContentType(fiber.MIMEApplicationJSON)

	return c.SendStatus(erix.HttpCode(err))
}

func ErrorWithCode(c fiber.Ctx, err error, code int) error {
	c.Response().SetBodyString(err.Error())
	c.Response().Header.SetContentType(fiber.MIMEApplicationJSON)

	return c.SendStatus(code)
}
