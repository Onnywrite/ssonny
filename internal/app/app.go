package app

import (
	"os"

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
	logger := zerolog.New(os.Stdout).
		Hook(zerolog.HookFunc(
			func(e *zerolog.Event, _ zerolog.Level, message string) {
				const skipFrames = 4

				e.Timestamp().Caller(skipFrames)
			}))

	// connecting to a database
	db, err := storage.New(cfg.Containerless.PostgresConn)
	if err != nil {
		logger.Fatal().Err(err).Msg("error while connecting to database")
	}

	tokensGenerator := tokens.New(
		cfg.Tokens.Issuer,
		cfg.Containerless.SecretString,
		cfg.Tokens.AccessTtl,
		cfg.Tokens.RefreshTtl,
		cfg.Tokens.IdTtl,
	)
	if err != nil {
		logger.Fatal().Err(err).Msg("error while creating tokens generator")
	}

	emailService, err := email.New(&logger)
	if err != nil {
		logger.Fatal().Err(err).Msg("error while creating tokens generator")
	}

	authService := auth.NewService(&logger, db, emailService, db, tokensGenerator)

	// creating grpc instance
	grpc := grpcapp.NewGRPC(&logger, grpcapp.Options{
		Port:           uint16(cfg.Grpc.Port),
		UseTLS:         cfg.Grpc.UseTls,
		CertPath:       cfg.Containerless.TlsCertPath,
		KeyPath:        cfg.Containerless.TlsKeyPath,
		Timeout:        cfg.Grpc.Timeout,
		CurrentService: cfg.Tokens.Issuer,
	}, grpcapp.Dependecies{})

	// creating http instance
	http := httpapp.New(&logger, httpapp.Options{
		Port:     uint16(cfg.Http.Port),
		UseTLS:   cfg.Http.UseTls,
		CertPath: cfg.Containerless.TlsCertPath,
		KeyPath:  cfg.Containerless.TlsKeyPath,
	}, httpapp.Dependecies{
		AuthService: authService,
		TokenParser: tokensGenerator,
	})

	return &Application{
		cfg:  cfg,
		log:  &logger,
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
