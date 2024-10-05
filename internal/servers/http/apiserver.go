package httpserver

import (
	api "github.com/Onnywrite/ssonny/api/oapi"
	"github.com/Onnywrite/ssonny/internal/servers/http/handlers"

	// authh "github.com/Onnywrite/ssonny/internal/servers/http/handlers/auth"
	// "github.com/Onnywrite/ssonny/internal/servers/http/middlewares"
	fiber "github.com/gofiber/fiber/v2"
)

// type AuthService interface {
// 	authh.Registrator
// 	authh.Loginer
// 	authh.Logouter
// 	authh.Refresher
// 	authh.Verifier
// }

// type TokenParser interface {
// 	middlewares.AccessTokenParser
// 	authh.RefreshTokenParser
// 	authh.EmailTokenParser
// }

func InitApi(r fiber.Router) { //authService AuthService, tokenParser TokenParser) {
	r.Get("/ping", handlers.Ping())
	// {
	// 	auth := r.Group("/auth")

	// 	auth.Post("/registerWithPassword", authh.RegisterWithPassword(authService))
	// 	auth.Post("/loginWithPassword", authh.LoginWithPassword(authService))
	// 	auth.Post("/refresh", authh.Refresh(authService, tokenParser))
	// 	auth.Post("/verify/email", authh.VerifyEmail(authService, tokenParser))
	// 	auth.Post("/logout", authh.Logout(authService, tokenParser))
	// }
	r.Static("/swagger", "cmd")

	sr := api.NewStrictHandler(nil, nil)
	api.RegisterHandlersWithOptions(r, sr, api.FiberServerOptions{
		BaseURL:     "/api",
		Middlewares: []api.MiddlewareFunc{},
	})
}
