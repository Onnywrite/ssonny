package apps

import (
	"context"
	"errors"
	"math"
	"math/rand/v2"

	"github.com/Onnywrite/ssonny/internal/domain/models"
	"github.com/Onnywrite/ssonny/internal/lib/erix"
	"github.com/Onnywrite/ssonny/internal/storage/repo"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"golang.org/x/crypto/bcrypt"
)

// CreateApp creates a new application for a given user, generates a secret for it,
// and optionally links provided domains to the application.
func (s *Service) CreateApp(
	ctx context.Context,
	ownerId uuid.UUID,
	data CreateAppData,
) (*AppCreated, error) {
	log := s.log.With().
		Stringer("owner_id", ownerId).
		Any("app_info", data).
		Logger()

	if err := data.Validate(s.validate); err != nil {
		log.Debug().Err(err).Msg("invalid data, bad request")

		return nil, erix.Wrap(err, erix.CodeBadRequest, ErrInvalidData)
	}

	verified, err := s.auth.IsVerified(ctx, ownerId)
	if err != nil {
		return nil, err
	}

	if !verified {
		log.Debug().Msg("user is not verified")

		return nil, erix.Wrap(ErrUserUnverified, erix.CodeForbidden, ErrUserUnverified)
	}

	randomSecret, hashed, err := s.randomSecretAndHash(&log)
	if err != nil {
		return nil, err
	}

	//nolint: gosec
	appId := uint64(rand.Int64N(math.MaxInt64-1) + 1)
	//nolint: exhaustruct
	app, tx, err := s.repo.SaveApp(ctx, models.App{
		Id:          appId,
		OwnerId:     ownerId,
		Name:        data.Name,
		Description: data.Description,
		SecretHash:  hashed,
	})

	switch {
	case errors.Is(err, repo.ErrUnique):
		log.Debug().Err(err).Msg("duplicating app")

		return nil, erix.Wrap(err, erix.CodeConflict, ErrAppAlreadyExists)
	case err != nil:
		log.Error().Err(err).Msg("error while saving app")

		return nil, erix.Wrap(err, erix.CodeInternalServerError, ErrInternal)
	}

	ctx = repo.WithTransactor(ctx, tx)
	defer tx.Rollback()

	err = s.LinkDomainsToApp(ctx, appId, data.DomainsIds)
	if err != nil {
		return nil, err
	}

	log.Debug().Msg("app created")

	return &AppCreated{
		Id:     app.Id,
		Secret: randomSecret,
	}, nil
}

func (s *Service) CreateDomains(
	ctx context.Context,
	ownerId uuid.UUID,
	domains []string,
) ([]models.Domain, error) {
	log := s.log.With().
		Stringer("owner_id", ownerId).
		Any("domains", domains).
		Logger()

	domainModels := make([]models.Domain, len(domains))
	//nolint: exhaustruct
	for i, d := range domains {
		domainModels[i] = models.Domain{
			Name:       d,
			IsVerified: false,
		}
	}

	savedDomains, err := s.repo.SaveDomains(ctx, ownerId, domainModels)

	switch {
	case errors.Is(err, repo.ErrUnique):
		log.Debug().Err(err).Msg("error while saving domains")

		return nil, erix.Wrap(err, erix.CodeConflict, ErrDomainAlreadyExists)
	case err != nil:
		log.Error().Err(err).Msg("error while saving domains")

		return nil, erix.Wrap(ErrInternal, erix.CodeInternalServerError, ErrInternal)
	}

	return savedDomains, nil
}

// LinkDomainsToApp links existent domains, owned by the user, to an application.
// It is idempotent.
func (s *Service) LinkDomainsToApp(
	ctx context.Context,
	appId uint64,
	domainsIds []uint64,
) error {
	if len(domainsIds) == 0 {
		return nil
	}

	log := s.log.With().
		Uint64("app_id", appId).
		Any("domains_ids", domainsIds).
		Logger()

	err := s.repo.TieDomainsToApp(ctx, appId, domainsIds)

	switch {
	case errors.Is(err, repo.ErrFK):
		log.Debug().Err(err).Msg("domains do not exist")

		return erix.Wrap(err, erix.CodeNotFound, ErrDomainNotFound)
	case err != nil && !errors.Is(err, repo.ErrUnique):
		log.Error().Err(err).Msg("error while linking domains to app")

		return erix.Wrap(err, erix.CodeInternalServerError, ErrInternal)
	}

	return nil
}

func (s *Service) randomSecretAndHash(log *zerolog.Logger) (string, string, error) {
	const maxLengthBcryptTakes = 72

	randomSecret, err := s.generator.RandomString()
	if err != nil {
		log.Error().Err(err).Msg("error while generating random secret")

		return "", "", erix.Wrap(err, erix.CodeInternalServerError, ErrInternal)
	}

	hashed, err := bcrypt.GenerateFromPassword(
		[]byte(randomSecret[:min(maxLengthBcryptTakes, len(randomSecret))]),
		bcrypt.DefaultCost)
	if err != nil {
		log.Error().Err(err).Msg("error while hashing app secret")

		return "", "", erix.Wrap(err, erix.CodeInternalServerError, ErrInternal)
	}

	return randomSecret, string(hashed), nil
}
