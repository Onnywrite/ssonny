package usersdeleter

import (
	"context"
	"time"

	"github.com/Onnywrite/ssonny/internal/config"
	"github.com/Onnywrite/ssonny/internal/services/email"
	"github.com/google/uuid"

	"github.com/rs/zerolog"
)

// Service is the worker responsible for invalidating outdated tokens.
type Service struct {
	logger    zerolog.Logger
	repo      UserRepo
	frequency time.Duration
	ttl       time.Duration
	sendEmail bool
}

type UserRepo interface {
	DeleteUnverifiedUsersWithCreatedAtBefore(context.Context, time.Time) ([]uuid.UUID, error)
}

type EmailService interface {
	SendNotification(context.Context, email.Notification) error
}

type Dependencies struct {
	TokenRepo UserRepo
}

type Config struct {
	Dependencies
	Frequency            time.Duration
	UnverifiedAccountTtl time.Duration
	SendEmail            bool
}

// New creates a new Service with dependencies and configuration from the application config.
func New(logger zerolog.Logger, deps Dependencies) *Service {
	conf := config.Get().Limits.UnverifiedAccount

	return &Service{
		logger:    logger,
		repo:      deps.TokenRepo,
		frequency: conf.Frequency,
		ttl:       conf.Ttl,
		sendEmail: conf.SendEmail,
	}
}

// NewWithConfig creates a new Service with the provided configuration.
func NewWithConfig(logger zerolog.Logger, c Config) *Service {
	return &Service{
		logger:    logger,
		repo:      c.TokenRepo,
		frequency: c.Frequency,
		ttl:       c.UnverifiedAccountTtl,
		sendEmail: c.SendEmail,
	}
}

// Run starts the user deletion process (blocking).
// It periodically deletes unverified users that have been registered
// before the specified TTL and sends a message to the deleted users.
func (s *Service) Run(ctx context.Context) {
	log := s.logger.With().Dur("frequency", s.frequency).Logger()
	ticker := time.NewTicker(s.frequency)

	for {
		// Create a new context with timeout for each iteration
		// to prevent potential stuck operations.
		dbCtx, cancel := context.WithTimeout(ctx, s.frequency)
		defer cancel()

		ids, err := s.repo.DeleteUnverifiedUsersWithCreatedAtBefore(dbCtx, time.Now().Add(-s.ttl))
		if err != nil {
			log.Error().Err(err).Msg("failed to delete users")
			goto Wait
		}
		if len(ids) == 0 {
			log.Debug().Msg("no users deleted")
			goto Wait
		}

		log.Info().Array("ids", UuidArray(ids)).Msg("deleted users")

		// Sorry, I really need to use lables
	Wait:
		// Select comes after the logic because there is the case when
		// it is stuck forever.
		// For example, if service restarts more often, than one tick happens.
		select {
		case <-ctx.Done():
			log.Info().Msg("stopping event processing")
			ticker.Stop()
			return
		case <-ticker.C:
		}
	}
}

// UuidArray is a helper type for logging [uuid.UUID] arrays with zerolog.
type UuidArray []uuid.UUID

// MarshalZerologArray implements [zerolog.Marshaler] for [Uuid].
func (a UuidArray) MarshalZerologArray(arr *zerolog.Array) {
	for _, id := range a {
		arr.Str(id.String())
	}
}
