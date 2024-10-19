package httpserver

import (
	"time"

	api "github.com/Onnywrite/ssonny/api/oapi"
	"github.com/Onnywrite/ssonny/internal/config"
	"github.com/Onnywrite/ssonny/internal/servers/http/handlers"
	authh "github.com/Onnywrite/ssonny/internal/servers/http/handlers/auth"
	usersh "github.com/Onnywrite/ssonny/internal/servers/http/handlers/users"
	"github.com/Onnywrite/ssonny/internal/servers/http/middlewares"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/limiter"
	"github.com/google/uuid"
)

type (
	AuthService  = authh.AuthService
	UsersService = usersh.UsersService
)

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

type Dependecies struct {
	AuthService            AuthService
	TokenParser            TokenParser
	UsersService           UsersService
	PasswordLimiterStorage fiber.Storage
}

type Config struct {
	Dependecies
	UpdatePasswordTimeout time.Duration
}

func InitApi(r fiber.Router, deps Dependecies) {
	c := config.Get()

	InitApiWithConfig(r, Config{
		Dependecies:           deps,
		UpdatePasswordTimeout: c.Limits.Password.ChangeTimeout,
	})
}

func InitApiWithConfig(r fiber.Router, c Config) {
	sr := api.NewStrictHandler(handler{
		AuthHandler: &authh.AuthHandler{
			Service:          c.AuthService,
			RefreshParser:    c.TokenParser,
			EmailTokenParser: c.TokenParser,
		},
		InternalHandler: &handlers.InternalHandler{},
		UsersHandler: &usersh.UsersHandler{
			Service: c.UsersService,
		},
	}, nil)

	passwordLimiterConfig := limiter.Config{
		Next:               skipLimiter,
		Max:                1,
		KeyGenerator:       passwordKeyGen,
		Expiration:         c.UpdatePasswordTimeout,
		SkipFailedRequests: true,
		Storage:            c.PasswordLimiterStorage,
		LimitReached:       passwordLimitReached,
		LimiterMiddleware:  middlewares.FixedWindow{},
	}

	api.RegisterHandlersWithOptions(r, sr, api.FiberServerOptions{
		Middlewares: nil,
		EndpointMiddlewares: map[string][]fiber.Handler{
			api.EP_GetAuthCheck: {middlewares.Authorization(c.TokenParser)},
			api.EP_GetProfile:   {middlewares.Authorization(c.TokenParser, "get:profile", "profile")},
			api.EP_PutProfile:   {middlewares.Authorization(c.TokenParser, "put:profile", "profile")},
			api.EP_PutProfilePassword: {
				middlewares.Authorization(c.TokenParser, "put:profile/password"),
				limiter.New(passwordLimiterConfig),
			},
		},
	})
}

func skipLimiter(c fiber.Ctx) bool {
	headers := c.GetReqHeaders()
	_, skipHeaderSet := headers["X-Skip-Limiter"]

	skipQuerySet := c.Query("skip_limiter", "") != ""

	return c.IP() == "127.0.0.1" && (skipHeaderSet || skipQuerySet)
}

func passwordLimitReached(c fiber.Ctx) error {
	return c.Status(fiber.StatusTooManyRequests).JSON(api.Err{
		Service: api.ErrServiceSsonny,
		Message: "password changed too recently, wait and try again later",
	})
}

func passwordKeyGen(c fiber.Ctx) string {
	return c.Locals("currentUserId").(uuid.UUID).String()
}
