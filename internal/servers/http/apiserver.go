package httpserver

import (
	"github.com/Onnywrite/ssonny/internal/servers/http/handlers"
	handlersapiauth "github.com/Onnywrite/ssonny/internal/servers/http/handlers/api/auth"
	"github.com/Onnywrite/ssonny/internal/servers/http/middlewares"

	"github.com/gofiber/fiber/v3"
)

type AuthService interface {
	handlersapiauth.Registrator
	handlersapiauth.Loginer
	handlersapiauth.Logouter
	handlersapiauth.Refresher
	handlersapiauth.Verifier
}

type TokenParser interface {
	middlewares.AccessTokenParser
	handlersapiauth.RefreshTokenParser
	handlersapiauth.EmailTokenParser
}

func InitApi(r fiber.Router, authService AuthService, tokenParser TokenParser) {
	r.Get("/ping", handlers.Ping())
	{
		auth := r.Group("/auth")

		auth.Post("/registerWithPassword", handlersapiauth.RegisterWithPassword(authService))
		auth.Post("/loginWithPassword", handlersapiauth.LoginWithPassword(authService))
		auth.Post("/refresh", handlersapiauth.Refresh(authService, tokenParser))
		auth.Post("/verify/email", handlersapiauth.VerifyEmail(authService, tokenParser))
		auth.Post("/logout",
			handlersapiauth.Logout(authService, tokenParser),
			middlewares.Authorization(tokenParser, "logout"),
		)
	}
}
