package tokensinvalidator

import (
	"context"
	"time"

	"github.com/Onnywrite/ssonny/internal/config"

	"github.com/rs/zerolog"
)

// Service is the worker responsible for invalidating outdated tokens.
type Service struct {
	logger     zerolog.Logger
	repo       TokenRepo
	frequency  time.Duration
	refreshTtl time.Duration
}

type TokenRepo interface {
	DeleteTokensWithRotatedAtBefore(context.Context, time.Time) ([]uint64, error)
}

type Dependencies struct {
	TokenRepo TokenRepo
}

type Config struct {
	Dependencies
	Frequency  time.Duration
	RefreshTtl time.Duration
}

// New creates a new Service with dependencies and configuration from the application config.
func New(logger zerolog.Logger, deps Dependencies) *Service {
	conf := config.Get()

	return &Service{
		logger:     logger,
		repo:       deps.TokenRepo,
		frequency:  conf.Limits.TokensInvalidationFrequency,
		refreshTtl: conf.Tokens.RefreshTtl,
	}
}

// NewWithConfig creates a new Service with the provided configuration.
func NewWithConfig(logger zerolog.Logger, c Config) *Service {
	return &Service{
		logger:     logger,
		repo:       c.TokenRepo,
		frequency:  c.Frequency,
		refreshTtl: c.RefreshTtl,
	}
}

// Run starts the token invalidation process (blocking).
// It periodically deletes tokens that have been rotated before the configured refresh TTL.
func (s *Service) Run(ctx context.Context) {
	log := s.logger.With().Dur("frequency", s.frequency).Logger()
	ticker := time.NewTicker(s.frequency)

	for {
		// Create a new context with timeout for each iteration
		// to prevent potential stuck operations.
		dbCtx, cancel := context.WithTimeout(ctx, s.frequency)
		defer cancel()

		ids, err := s.repo.DeleteTokensWithRotatedAtBefore(dbCtx, time.Now().Add(-s.refreshTtl))
		if err != nil {
			log.Error().Err(err).Msg("failed to delete tokens")
			goto Wait
		}
		if len(ids) == 0 {
			log.Debug().Msg("no tokens deleted")
			goto Wait
		}

		log.Info().Array("ids", UInt64Array(ids)).Msg("deleted tokens")

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

// UInt64Array is a helper type for logging uint64 arrays with zerolog.
type UInt64Array []uint64

// MarshalZerologArray implements [zerolog.Marshaler] for [UInt64Array].
func (a UInt64Array) MarshalZerologArray(arr *zerolog.Array) {
	for _, id := range a {
		arr.Uint64(id)
	}
}
