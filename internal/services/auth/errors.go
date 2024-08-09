package auth

import (
	"errors"

	"github.com/Onnywrite/ssonny/internal/lib/erix"
	"github.com/Onnywrite/ssonny/internal/storage/repo"
	"github.com/rs/zerolog"
)

var (
	ErrUserNotFound       = errors.New("user not found")
	ErrUserAlreadyExists  = errors.New("user already exists")
	ErrInvalidData        = errors.New("user has invalid data")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrDependencyNotFound = errors.New("data, user depends on, not found")
	ErrEmailUnverified    = errors.New("email cannot be verified whatsoever")
	ErrInternal           = errors.New("internal error")
)

var (
	ErrSessionNotFound           = errors.New("session not found")
	ErrSessionOpened             = errors.New("session is already opened")
	ErrSessionInvalidData        = errors.New("session has invalid data")
	ErrSessionDependencyNotFound = errors.New("user or application not found")
)

func userFailed(log *zerolog.Logger, err error) error {
	switch {
	case errors.Is(err, repo.ErrEmptyResult):
		log.Info().Err(err).Msg("empty result when getting user")
		return erix.Wrap(err, erix.CodeNotFound, ErrUserNotFound)
	case errors.Is(err, repo.ErrUnique):
		log.Info().Err(err).Msg("cannot create user")
		return erix.Wrap(err, erix.CodeConflict, ErrUserAlreadyExists)
	case errors.Is(err, repo.ErrChecked):
		log.Info().Err(err).Msg("some data does not satisfy repo, maybe validation is not up to date with database schema")
		return erix.Wrap(err, erix.CodeBadRequest, ErrInvalidData)
	case errors.Is(err, repo.ErrFK):
		log.Info().Err(err).Msg("FK constraint violation, either row with this PK does not exist or you forgot 'delete cascade'")
		return erix.Wrap(err, erix.CodeNotFound, ErrDependencyNotFound)
	case errors.Is(err, repo.ErrNull):
		log.Info().Err(err).Msg("cannot create user")
		return erix.Wrap(err, erix.CodeBadRequest, ErrInvalidData)
		// repo.ErrDataInconsistent and repo.ErrInternal are just ErrInternal
	case err != nil:
		log.Error().Err(err).Msg("error while operating on user")
		return erix.Wrap(err, erix.CodeInternalServerError, ErrInternal)
	}
	log.Error().Msg("not nil error passed")
	panic("not nil error passed, check log for details")
}
