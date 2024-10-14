package users

import (
	"context"
	"errors"

	"github.com/Onnywrite/ssonny/internal/lib/erix"
	"github.com/Onnywrite/ssonny/internal/storage/repo"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func (s *Service) PutProfilePassword(ctx context.Context,
	userId uuid.UUID,
	data UpdatePasswordData,
) error {
	log := s.log.With().Stringer("user_id", userId).Logger()

	user, err := s.repo.UserById(ctx, userId)
	switch {
	case errors.Is(err, repo.ErrEmptyResult):
		log.Debug().Err(err).Msg("user not found")

		return erix.Wrap(err, erix.CodeNotFound, ErrInvalidCredentials)
	case err != nil:
		log.Error().Err(err).Msg("failed to get user")

		return erix.Wrap(err, erix.CodeInternalServerError, ErrInternal)
	}

	if user.PasswordHash != nil {
		err = bcrypt.CompareHashAndPassword([]byte(*user.PasswordHash), []byte(data.CurrentPassword))
		if err != nil {
			log.Info().Err(err).Msg("invalid password")

			return erix.Wrap(err, erix.CodeNotFound, ErrInvalidCredentials)
		}
	}

	newHash, err := bcrypt.GenerateFromPassword([]byte(data.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		log.Error().Err(err).Msg("failed to hash password")

		return erix.Wrap(err, erix.CodeInternalServerError, ErrInternal)
	}

	err = s.repo.UpdateUser(ctx, userId, map[string]any{
		"user_password_hash": string(newHash),
	})
	if err != nil {
		log.Error().Err(err).Msg("failed to update user")

		return erix.Wrap(err, erix.CodeInternalServerError, ErrInternal)
	}

	return nil
}
