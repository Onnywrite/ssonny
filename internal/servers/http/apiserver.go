package httpserver

import (
	"github.com/Onnywrite/ssonny/internal/servers/http/handlers"
	authh "github.com/Onnywrite/ssonny/internal/servers/http/handlers/auth"
	"github.com/Onnywrite/ssonny/internal/servers/http/middlewares"
	"github.com/gofiber/fiber/v3"
)

type AuthService interface {
	authh.Registrator
	authh.Loginer
	authh.Logouter
	authh.Refresher
	authh.Verifier
}

type TokenParser interface {
	middlewares.AccessTokenParser
	authh.RefreshTokenParser
	authh.EmailTokenParser
}

func InitApi(r fiber.Router, authService AuthService, tokenParser TokenParser) {
	r.Get("/ping", handlers.Ping())
	{
		auth := r.Group("/auth")

		auth.Post("/registerWithPassword", authh.RegisterWithPassword(authService))
		auth.Post("/loginWithPassword", authh.LoginWithPassword(authService))
		auth.Post("/refresh", authh.Refresh(authService, tokenParser))
		auth.Post("/verify/email", authh.VerifyEmail(authService, tokenParser))
		auth.Post("/logout", authh.Logout(authService, tokenParser))
	}
}
