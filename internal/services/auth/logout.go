package auth

import (
	"context"
	"errors"

	"github.com/Onnywrite/ssonny/internal/lib/erix"
	"github.com/Onnywrite/ssonny/internal/storage/repo"
)

func (s *Service) Logout(ctx context.Context, jwtId uint64) error {
	log := s.log.With().Uint64("jwt_id", jwtId).Logger()

	err := s.tokenRepo.DeleteToken(ctx, jwtId)
	if errors.Is(err, repo.ErrEmptyResult) {
		log.Info().Msg("empty result while deleting token")
		return nil
	}
	if err != nil {
		log.Error().Err(err).Msg("error while deleting token")
		return erix.Wrap(err, erix.CodeInternalServerError, ErrInternal)
	}

	return nil
}
