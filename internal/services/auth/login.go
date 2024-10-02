package auth

import (
	"context"
	"errors"

	"github.com/Onnywrite/ssonny/internal/domain/models"
	"github.com/Onnywrite/ssonny/internal/lib/erix"
	"github.com/Onnywrite/ssonny/internal/storage/repo"

	"github.com/rs/zerolog"
	"golang.org/x/crypto/bcrypt"
)

func (s *Service) LoginWithPassword(ctx context.Context, data LoginWithPasswordData,
) (*AuthenticatedUser, error) {
	if err := data.Validate(s.validate); err != nil {
		s.log.Debug().Err(err).Msg("invalid data, bad request")

		return nil, erix.Wrap(err, erix.CodeBadRequest, ErrInvalidData)
	}

	var (
		user *models.User
		err  error
		log  zerolog.Logger
	)

	if data.Email != nil {
		user, err = s.repo.UserByEmail(ctx, *data.Email)
		log = s.log.With().Str("login", *data.Email).Logger()
	} else if data.Nickname != nil {
		log = s.log.With().Str("login", *data.Nickname).Logger()
		user, err = s.repo.UserByNickname(ctx, *data.Nickname)
	}

	switch {
	case errors.Is(err, repo.ErrEmptyResult):
		log.Debug().Err(err).Msg("empty result when getting user")

		return nil, erix.Wrap(err, erix.CodeNotFound, ErrInvalidCredentials)
	case err != nil:
		log.Error().Err(err).Msg("error while getting user")

		return nil, erix.Wrap(err, erix.CodeInternalServerError, ErrInternal)
	}

	if err := bcrypt.CompareHashAndPassword(
		[]byte(*user.PasswordHash),
		[]byte(data.Password)); err != nil {
		log.Debug().Msg("invalid password")

		return nil, erix.Wrap(err, erix.CodeNotFound, ErrInvalidCredentials)
	}

	return s.generateAndSaveTokens(ctx, *user, data.UserInfo)
}
