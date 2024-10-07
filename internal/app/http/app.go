package httpapp

import (
	"errors"
	"fmt"

	httpapi "github.com/Onnywrite/ssonny/api/oapi"
	httpserver "github.com/Onnywrite/ssonny/internal/servers/http"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/recover"
	"github.com/rs/zerolog"
)

type App struct {
	log    *zerolog.Logger
	server *fiber.App
	port   string

	useTls   bool
	certPath string
	keyPath  string
}

type Options struct {
	Port     int
	UseTLS   bool
	CertPath string
	KeyPath  string
}

type Dependecies struct {
	AuthService httpserver.AuthService
	TokenParser httpserver.TokenParser
}

func New(logger *zerolog.Logger, opts Options, deps Dependecies) *App {
	httpLogger := logger.With().Logger()
	//nolint: exhaustruct
	app := fiber.New(fiber.Config{
		ErrorHandler: FiberErrorHandler,
	})

	app.Use(logging(&httpLogger))
	//nolint: exhaustruct
	app.Use(recover.New(recover.Config{EnableStackTrace: true}))

	httpserver.InitApi(app.Group("/api"), deps.AuthService, deps.TokenParser)

	return &App{
		log:      logger,
		server:   app,
		port:     fmt.Sprintf(":%d", opts.Port),
		useTls:   opts.UseTLS,
		certPath: opts.CertPath,
		keyPath:  opts.KeyPath,
	}
}

func (a *App) MustStart() {
	if err := a.Start(); err != nil {
		panic(err)
	}
}

func (a *App) Start() error {
	go func() {
		//nolint: exhaustruct
		config := fiber.ListenConfig{
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

		err := a.server.Listen(a.port, config)
		if err != nil {
			a.log.Error().Err(err).Msg("error while starting http")

			return
		}
	}()

	a.log.Info().Str("port", a.port).Msg("http started")

	return nil
}

func (a *App) Stop() error {
	if err := a.server.Shutdown(); err != nil {
		a.log.Error().Err(err).Msg("error while stopping http")

		return err
	}

	a.log.Info().Str("port", a.port).Msg("stopped http")

	return nil
}

func FiberErrorHandler(c fiber.Ctx, err error) error {
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
