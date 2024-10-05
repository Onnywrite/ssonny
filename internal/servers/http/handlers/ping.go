package handlers

import "github.com/gofiber/fiber/v2"

func Ping() func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).SendString("pong")
	}
}
