package httpapp

import (
	"errors"
	"fmt"

	httpapi "github.com/Onnywrite/ssonny/api/oapi"
	httpserver "github.com/Onnywrite/ssonny/internal/servers/http"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/gofiber/fiber/v3/middleware/recover"
	"github.com/rs/zerolog"
)

// App is a http application.
type App struct {
	log    zerolog.Logger
	server *fiber.App
	port   string

	useTls   bool
	certPath string
	keyPath  string
}

type FiberSubstorager interface {
	FiberSubstorage(namespace string) fiber.Storage
}

// Config is a http application configuration.
type Config struct {
	Port         int
	UseTls       bool
	CertPath     string
	KeyPath      string
	AuthService  httpserver.AuthService
	TokenParser  httpserver.TokenParser
	UsersService httpserver.UsersService
	Substorager  FiberSubstorager
}

// New creates a new HTTP application.
func New(logger zerolog.Logger, conf Config) *App {
	app := fiber.New(fiber.Config{ //nolint: exhaustruct
		ErrorHandler: fiberErrorHandler,
	})

	applyMiddlewares(logger, app)

	httpserver.InitApi(app.Group("/api"), httpserver.Dependecies{
		AuthService:            conf.AuthService,
		TokenParser:            conf.TokenParser,
		UsersService:           conf.UsersService,
		PasswordLimiterStorage: conf.Substorager.FiberSubstorage("f-lim-pswd"),
	})

	return &App{
		log:      logger,
		server:   app,
		port:     fmt.Sprintf(":%d", conf.Port),
		useTls:   conf.UseTls,
		certPath: conf.CertPath,
		keyPath:  conf.KeyPath,
	}
}

// Start starts the http application.
func (a *App) Start() error {
	config := fiber.ListenConfig{ //nolint: exhaustruct
		DisableStartupMessage: true,
	}

	if a.useTls {
		a.log.Info().
			Str("cert_path", a.certPath).
			Str("key_path", a.keyPath).
			Msg("http uses TLS certificate")

		config.CertFile = a.certPath
		config.CertKeyFile = a.keyPath
	}

	go func() {
		err := a.server.Listen(a.port, config)
		if err != nil {
			a.log.Error().Err(err).Msg("error while starting http")
		}
	}()

	a.log.Info().Str("port", a.port).Msg("http started")

	return nil
}

// Stop stops the http application.
func (a *App) Stop() error {
	if err := a.server.Shutdown(); err != nil {
		a.log.Error().Err(err).Msg("error while stopping http")

		return err
	}

	a.log.Info().Str("port", a.port).Msg("stopped http")

	return nil
}

// fiberErrorHandler is a custom error handler for fiber.
func fiberErrorHandler(c fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError

	var e *fiber.Error
	if errors.As(err, &e) {
		code = e.Code
	}

	c.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)

	return c.Status(code).JSON(httpapi.Err{
		Service: httpapi.ErrServiceSsonny,
		Message: err.Error(),
	})
}

// applyMiddlewares applies middlewares to the fiber application.
func applyMiddlewares(logger zerolog.Logger, app *fiber.App) {
	app.Use(logging(logger))
	app.Use(recover.New(recover.Config{EnableStackTrace: true})) //nolint: exhaustruct

	//nolint: exhaustruct
	app.Use(cors.New(cors.Config{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{"Origin", "Content-Type", "Accept", "Authorization", "User-Agent"},
		AllowMethods: []string{"GET", "POST", "HEAD", "PUT", "DELETE", "PATCH"},
	}))
}
