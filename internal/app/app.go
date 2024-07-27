package app

import (
	"os"

	grpcapp "github.com/Onnywrite/ssonny/internal/app/grpc"
	httpapp "github.com/Onnywrite/ssonny/internal/app/http"
	"github.com/Onnywrite/ssonny/internal/config"
	"github.com/rs/zerolog"
)

type Application struct {
	cfg  *config.Config
	log  *zerolog.Logger
	http *httpapp.App
	grpc *grpcapp.App
}

func New(cfg *config.Config) *Application {
	l := zerolog.New(os.Stdout).
		Hook(zerolog.HookFunc(
			func(e *zerolog.Event, level zerolog.Level, message string) {
				e.Timestamp().Caller(4)
			}))

	grpc := grpcapp.NewGRPC(&l, grpcapp.Options{
		Port:           cfg.Grpc.Port,
		Timeout:        cfg.Grpc.Timeout,
		CurrentService: "ssonny",
	}, grpcapp.Dependecies{})

	http := httpapp.New(&l, httpapp.Options{
		Port: cfg.Https.Port,
	}, httpapp.Dependecies{})

	return &Application{
		cfg:  cfg,
		log:  &l,
		http: http,
		grpc: grpc,
	}
}

func (a *Application) MustStart() {
	if err := a.Start(); err != nil {
		panic(err)
	}
}

func (a *Application) Start() error {
	a.log.Info().Msg("starting application")
	if err := a.http.Start(); err != nil {
		return err
	}
	if err := a.grpc.Start(); err != nil {
		return err
	}
	a.log.Info().Msg("started application")
	return nil
}

func (a *Application) MustStop() {
	if err := a.Stop(); err != nil {
		panic(err)
	}
}

func (a *Application) Stop() error {
	a.log.Info().Msg("stopping application")
	a.grpc.Stop()
	if err := a.http.Stop(); err != nil {
		return err
	}
	a.log.Info().Msg("stopped application")
	return nil
}
