package auth

import (
	"context"

	"github.com/google/uuid"
)

func (s *Service) VerifyEmail(ctx context.Context, userId uuid.UUID) error {
	log := s.log.With().Stringer("user_id", userId).Logger()

	err := s.repo.UpdateUser(ctx, userId, map[string]any{
		"user_verified": true,
	})
	if err != nil {
		return userFailed(log, err)
	}

	return nil
}
