package httpserver

import (
	api "github.com/Onnywrite/ssonny/api/oapi"
	"github.com/Onnywrite/ssonny/internal/servers/http/handlers"
	authh "github.com/Onnywrite/ssonny/internal/servers/http/handlers/auth"
	usersh "github.com/Onnywrite/ssonny/internal/servers/http/handlers/users"
	"github.com/Onnywrite/ssonny/internal/servers/http/middlewares"

	"github.com/gofiber/fiber/v3"
)

type AuthService = authh.AuthService
type UsersService = usersh.UsersService

type TokenParser interface {
	middlewares.AccessTokenParser
	authh.RefreshTokenParser
	authh.EmailTokenParser
}

type handler struct {
	*authh.AuthHandler
	*handlers.InternalHandler
	*usersh.UsersHandler
}

func InitApi(r fiber.Router, authService AuthService, tokenParser TokenParser, usersService UsersService) {
	sr := api.NewStrictHandler(handler{
		AuthHandler: &authh.AuthHandler{
			Service:          authService,
			RefreshParser:    tokenParser,
			EmailTokenParser: tokenParser,
		},
		InternalHandler: &handlers.InternalHandler{},
		UsersHandler: &usersh.UsersHandler{
			Service: usersService,
		},
	}, nil)

	api.RegisterHandlersWithOptions(r, sr, api.FiberServerOptions{
		Middlewares: nil,
		EndpointMiddlewares: map[string][]fiber.Handler{
			api.EP_GetAuthCheck: {middlewares.Authorization(tokenParser)},
			api.EP_GetProfile:   {middlewares.Authorization(tokenParser, "get:profile", "profile")},
		},
	})
}
