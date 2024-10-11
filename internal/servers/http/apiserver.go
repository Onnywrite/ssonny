package httpserver

import (
	api "github.com/Onnywrite/ssonny/api/oapi"
	"github.com/Onnywrite/ssonny/internal/servers/http/handlers"
	authh "github.com/Onnywrite/ssonny/internal/servers/http/handlers/auth"
	"github.com/Onnywrite/ssonny/internal/servers/http/middlewares"

	"github.com/gofiber/fiber/v3"
)

type AuthService = authh.AuthService

type TokenParser interface {
	middlewares.AccessTokenParser
	authh.RefreshTokenParser
	authh.EmailTokenParser
}

type handler struct {
	*authh.AuthHandler
	*handlers.InternalHandler
}

func InitApi(r fiber.Router, authService AuthService, tokenParser TokenParser) {
	sr := api.NewStrictHandler(handler{
		AuthHandler: &authh.AuthHandler{
			Service:          authService,
			RefreshParser:    tokenParser,
			EmailTokenParser: tokenParser,
		},
		InternalHandler: &handlers.InternalHandler{},
	}, nil)

	api.RegisterHandlersWithOptions(r, sr, api.FiberServerOptions{
		Middlewares: nil,
		EndpointMiddlewares: map[string][]fiber.Handler{
			api.EP_GetAuthCheck: {middlewares.Authorization(tokenParser)},
		},
	})
}
