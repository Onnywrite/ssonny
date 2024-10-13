package app

import (
	"context"

	grpcapp "github.com/Onnywrite/ssonny/internal/app/grpc"
	httpapp "github.com/Onnywrite/ssonny/internal/app/http"
	"github.com/Onnywrite/ssonny/internal/config"
	"github.com/Onnywrite/ssonny/internal/storage"
	"github.com/Onnywrite/ssonny/pkg/must"

	"github.com/rs/zerolog"
)

// Application represents the top-level application structure.
type Application struct {
	cfg  config.Config
	log  zerolog.Logger
	http *httpapp.App
	grpc *grpcapp.App
	db   *storage.Storage
}

// New creates a new Application instance with loaded configuration.
func New() *Application {
	cfg := must.Ok2(config.Load("sso.yaml", "/etc/sso/sso.yaml"))
	config.Set(cfg)

	return NewWithConfig(cfg)
}

// NewWithConfig creates a new Application instance with the provided configuration.
func NewWithConfig(cfg config.Config) *Application {
	logger := newLogger(cfg)

	db, err := newStorage(cfg)
	if err != nil {
		logger.Fatal().Err(err).Msg("error while connecting to database")
	}

	grpc, http := newApps(logger, cfg, db)

	return &Application{
		cfg:  cfg,
		log:  logger,
		http: http,
		grpc: grpc,
		db:   db,
	}
}

// Run starts the application and listens for shutdown signals.
func (a *Application) Run(ctx context.Context) error {
	if err := a.start(); err != nil {
		return err
	}

	<-ctx.Done()

	a.log.Info().Msg("shutting down application")

	if err := a.shutdown(); err != nil {
		a.log.Error().Err(err).Msg("error while shutting down application")

		return err
	}

	return nil
}

// start starts the application components.
func (a *Application) start() error {
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

// shutdown gracefully stops the application components.
func (a *Application) shutdown() error {
	a.log.Info().Msg("stopping application")

	a.grpc.Stop()

	a.log.Debug().Msg("disconnecting from database")

	if err := a.db.Disconnect(); err != nil {
		return err
	}

	a.log.Debug().Msg("disconnected from database")

	if err := a.http.Stop(); err != nil {
		return err
	}

	a.log.Info().Msg("stopped application")

	return nil
}

func newStorage(cfg config.Config) (*storage.Storage, error) {
	pgconf := cfg.Secrets.Postgres
	rdconf := cfg.Secrets.Redis

	db, err := storage.New(storage.PostgresConfig{
		Host:     pgconf.Host,
		Port:     pgconf.Port,
		Username: pgconf.User,
		Password: pgconf.Password,
		Database: pgconf.Database,
		SslMode:  pgconf.SslMode,
	}, storage.RedisConfig{
		Host:     rdconf.Host,
		Port:     rdconf.Port,
		Username: rdconf.User,
		Password: rdconf.Password,
		Db:       rdconf.Db,
	})
	if err != nil {
		return nil, err
	}

	return db, nil
}
