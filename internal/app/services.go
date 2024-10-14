package app

import (
	grpcapp "github.com/Onnywrite/ssonny/internal/app/grpc"
	httpapp "github.com/Onnywrite/ssonny/internal/app/http"
	"github.com/Onnywrite/ssonny/internal/config"
	"github.com/Onnywrite/ssonny/internal/lib/tokens"
	"github.com/Onnywrite/ssonny/internal/services/auth"
	"github.com/Onnywrite/ssonny/internal/services/email"
	"github.com/Onnywrite/ssonny/internal/services/users"
	"github.com/Onnywrite/ssonny/internal/storage"

	"github.com/rs/zerolog"
)

// newApps initializes and returns the gRPC and HTTP applications with their dependencies.
func newApps(logger zerolog.Logger, cfg config.Config, db *storage.Storage,
) (*grpcapp.App, *httpapp.App) {
	tokensGenerator := tokens.New()

	emailService, err := email.New(logger)
	if err != nil {
		logger.Fatal().Err(err).Msg("error while creating email service")
	}

	authService := auth.NewService(logger, auth.Config{
		UserRepo:     db,
		EmailService: emailService,
		TokenRepo:    db,
		TokensSigner: tokensGenerator,
	})

	usersService := users.NewService(logger, users.Config{
		UserRepo:     db,
		EmailService: emailService,
	})

	// Initialize the gRPC configuration.
	grpcConfig := grpcapp.Config{
		Port:     cfg.Grpc.Port,
		UseTls:   cfg.Grpc.UseTls,
		CertPath: cfg.Tls.CertPath,
		KeyPath:  cfg.Tls.KeyPath,
		Timeout:  cfg.Grpc.Timeout,
	}

	// Initialize the HTTP configuration.
	httpConfig := httpapp.Config{
		Port:        cfg.Http.Port,
		UseTls:      cfg.Http.UseTls,
		CertPath:    cfg.Tls.CertPath,
		KeyPath:     cfg.Tls.KeyPath,
		Substorager: db,
		// Inject dependencies.
		AuthService:  authService,
		TokenParser:  tokensGenerator,
		UsersService: usersService,
	}

	return grpcapp.New(logger, grpcConfig), httpapp.New(logger, httpConfig)
}
