package httpapp

import (
	"fmt"

	httpserver "github.com/Onnywrite/ssonny/internal/servers/http"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/recover"
	"github.com/rs/zerolog"
)

type App struct {
	log    *zerolog.Logger
	server *fiber.App
	port   string
}

type Options struct {
	Port uint16
}

type Dependecies struct {
	AuthService httpserver.AuthService
}

func New(logger *zerolog.Logger, opts Options, deps Dependecies) *App {
	httpLogger := logger.With().Logger()
	s := fiber.New()
	s.Use(logging(&httpLogger))
	s.Use(recover.New(recover.Config{EnableStackTrace: true}))

	httpserver.InitApi(s.Group("/api"), deps.AuthService)

	return &App{
		log:    logger,
		server: s,
		port:   fmt.Sprintf(":%d", opts.Port),
	}
}

func (a *App) MustStart() {
	if err := a.Start(); err != nil {
		panic(err)
	}
}

func (a *App) Start() error {
	go func() {
		if err := a.server.Listen(a.port, fiber.ListenConfig{}); err != nil {
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
