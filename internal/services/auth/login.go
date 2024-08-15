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

type LoginWithPasswordData struct {
	// one of
	Email    *string
	Nickname *string
	Password string
	UserInfo UserInfo
}

func (s *Service) LoginWithPassword(ctx context.Context, data LoginWithPasswordData) (*AuthenticatedUser, error) {
	// TODO: validate data
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
	} else {
		s.log.Debug().Msg("email and nickname are nil")
		return nil, erix.Wrap(ErrInvalidData, erix.CodeBadRequest, ErrInvalidData)
	}
	switch {
	case errors.Is(err, repo.ErrEmptyResult):
		log.Debug().Err(err).Msg("empty result when getting user")
		return nil, erix.Wrap(err, erix.CodeNotFound, ErrInvalidCredentials)
	case err != nil:
		log.Error().Err(err).Msg("error while getting user")
		return nil, erix.Wrap(err, erix.CodeInternalServerError, ErrInternal)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(*user.PasswordHash), []byte(data.Password)); err != nil {
		log.Debug().Msg("invalid password")
		return nil, erix.Wrap(err, erix.CodeNotFound, ErrInvalidCredentials)
	}

	return s.generateAndSaveTokens(ctx, *user, data.UserInfo)
}
