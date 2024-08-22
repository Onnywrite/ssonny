package auth

import (
	"context"
	"errors"
	"time"

	"github.com/Onnywrite/ssonny/internal/lib/erix"
	"github.com/Onnywrite/ssonny/internal/lib/tokens"
	"github.com/Onnywrite/ssonny/internal/storage/repo"
)

type Tokens struct {
	Refresh string
	Access  string
}

func (s *Service) Refresh(ctx context.Context, refresh tokens.Refresh) (*Tokens, error) {
	log := s.log.With().
		Uint64("jwt_id", refresh.Id).
		Any("app_id", refresh.Audience).
		Stringer("user_id", refresh.Subject).Logger()

	token, err := s.tokenRepo.Token(ctx, refresh.Id)

	switch {
	case errors.Is(err, repo.ErrEmptyResult):
		log.Debug().Err(err).Msg("empty result while getting token")
		return nil, erix.Wrap(err, erix.CodeUnauthorized, tokens.ErrExpired)
	case err != nil:
		log.Error().Err(err).Msg("error while getting token")
		return nil, erix.Wrap(err, erix.CodeInternalServerError, ErrInternal)
	}

	if token.Rotation != refresh.Rotation {
		log.Warn().Msg("invalid rotation number. Invalidating")

		if err = s.tokenRepo.DeleteTokens(ctx, token.UserId, token.AppId); err != nil {
			log.Error().Err(err).Msg("could not invalidate sus tokens")
			return nil, erix.Wrap(err, erix.CodeUnauthorized, ErrInvalidTokenRotation)
		}
		return nil, erix.Wrap(ErrInvalidTokenRotation, erix.CodeUnauthorized, ErrInvalidTokenRotation)
	}

	newRotation := token.Rotation + 1

	err = s.tokenRepo.UpdateToken(ctx, token.Id, map[string]any{
		"token_rotation":   newRotation,
		"token_rotated_at": time.Now(),
	})
	if err != nil {
		log.Error().Err(err).Msg("could not update token rotation")
		return nil, erix.Wrap(err, erix.CodeUnauthorized, ErrInternal)
	}

	newAccess, err := s.signer.SignAccess(token.UserId, token.AppId, "self", "*")
	if err != nil {
		log.Error().Err(err).Msg("error while signing access token")
		return nil, erix.Wrap(err, erix.CodeInternalServerError, ErrInternal)
	}

	newRefresh, err := s.signer.SignRefresh(token.UserId, token.AppId, "self", newRotation, token.Id)
	if err != nil {
		log.Error().Err(err).Msg("error while signing refresh token")
		return nil, erix.Wrap(err, erix.CodeInternalServerError, ErrInternal)
	}

	return &Tokens{
		Refresh: newRefresh,
		Access:  newAccess,
	}, nil
}
