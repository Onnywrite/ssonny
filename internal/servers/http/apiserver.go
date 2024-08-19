package httpserver

import (
	"github.com/Onnywrite/ssonny/internal/servers/http/handlers"
	handlers_api_auth "github.com/Onnywrite/ssonny/internal/servers/http/handlers/api/auth"
	"github.com/gofiber/fiber/v3"
)

type AuthService interface {
	handlers_api_auth.Registrator
}

func InitApi(r fiber.Router, authService AuthService) {
	r.Get("/ping", handlers.Ping())
	{
		auth := r.Group("/auth")

		auth.Post("/registerWithPassword", handlers_api_auth.RegisterWithPassword(authService))
	}
}
