package app

import (
	"crypto/rand"
	"crypto/rsa"
	"os"
	"time"

	grpcapp "github.com/Onnywrite/ssonny/internal/app/grpc"
	httpapp "github.com/Onnywrite/ssonny/internal/app/http"
	"github.com/Onnywrite/ssonny/internal/config"
	"github.com/Onnywrite/ssonny/internal/lib/tokens"
	"github.com/Onnywrite/ssonny/internal/services/auth"
	"github.com/Onnywrite/ssonny/internal/services/email"
	"github.com/Onnywrite/ssonny/internal/storage"
	"github.com/rs/zerolog"
)

type Application struct {
	cfg  *config.Config
	log  *zerolog.Logger
	http *httpapp.App
	grpc *grpcapp.App
	db   *storage.Storage
}

func New(cfg *config.Config) *Application {
	// setting up the logger
	l := zerolog.New(os.Stdout).
		Hook(zerolog.HookFunc(
			func(e *zerolog.Event, level zerolog.Level, message string) {
				e.Timestamp().Caller(4)
			}))

	// connecting to a database
	db, err := storage.New(cfg.PostgresConn)
	if err != nil {
		l.Fatal().Err(err).Msg("error while connecting to database")
	}

	// initializing all services and its dependencies
	// TODO: nice config
	key, err := rsa.GenerateKey(rand.Reader, 1024)
	if err != nil {
		l.Fatal().Err(err).Msg("error while generating rsa key")
	}
	tokensGenerator := tokens.NewWithKeys("sso.onnywrite.com", time.Minute*5, time.Hour*240, time.Hour*24, &key.PublicKey, key)
	if err != nil {
		l.Fatal().Err(err).Msg("error while creating tokens generator")
	}

	emailService, err := email.New(&l)
	if err != nil {
		l.Fatal().Err(err).Msg("error while creating tokens generator")
	}

	authService := auth.NewService(&l, db, emailService, db, tokensGenerator)

	// creating grpc instance
	grpc := grpcapp.NewGRPC(&l, grpcapp.Options{
		Port:           cfg.Grpc.Port,
		Timeout:        cfg.Grpc.Timeout,
		CurrentService: "ssonny",
	}, grpcapp.Dependecies{})

	// creating http instance
	http := httpapp.New(&l, httpapp.Options{
		Port: cfg.Https.Port,
	}, httpapp.Dependecies{
		AuthService: authService,
	})

	return &Application{
		cfg:  cfg,
		log:  &l,
		http: http,
		grpc: grpc,
		db:   db,
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
	if err := a.db.Disconnect(); err != nil {
		return err
	}
	if err := a.http.Stop(); err != nil {
		return err
	}
	a.log.Info().Msg("stopped application")
	return nil
}
