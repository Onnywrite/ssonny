package users

import (
	"context"
	"errors"

	"github.com/Onnywrite/ssonny/internal/lib/erix"
	"github.com/Onnywrite/ssonny/internal/storage/repo"

	"github.com/google/uuid"
)

func (s *Service) GetProfile(ctx context.Context,
	userId uuid.UUID,
) (*Profile, error) {
	log := s.log.With().Stringer("user_id", userId).Logger()

	user, err := s.repo.UserById(ctx, userId)
	switch {
	case errors.Is(err, repo.ErrEmptyResult):
		log.Debug().Err(err).Msg("user not found")

		return nil, erix.Wrap(err, erix.CodeNotFound, ErrUserNotFound)
	case err != nil:
		log.Error().Err(err).Msg("failed to get user")

		return nil, erix.Wrap(err, erix.CodeInternalServerError, ErrInternal)
	}

	profile := mapProfile(user)

	return &profile, nil
}

func (s *Service) PutProfile(ctx context.Context,
	userId uuid.UUID,
	data UpdateProfileData,
) (*Profile, error) {
	log := s.log.With().Any("profile", data).Logger()

	updated, err := s.repo.UpdateAndGetUser(ctx, userId, map[string]any{
		"user_nickname": data.Nickname,
		"user_gender":   data.Gender,
		"user_birthday": data.Birthday,
	})
	if err != nil {
		log.Error().Err(err).Msg("failed to update user")

		return nil, erix.Wrap(err, erix.CodeInternalServerError, ErrInternal)
	}

	profile := mapProfile(updated)

	return &profile, nil
}
