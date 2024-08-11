package auth

import (
	"context"
	"errors"
	"time"

	"github.com/Onnywrite/ssonny/internal/domain/models"
	"github.com/Onnywrite/ssonny/internal/lib/erix"
	"github.com/Onnywrite/ssonny/internal/lib/tokens"
	"github.com/Onnywrite/ssonny/internal/storage/repo"
)

type Tokens struct {
	Refresh string
	Access  string
}

func (s *Service) Refresh(ctx context.Context, refreshToken string) (*Tokens, error) {
	refresh, err := s.tokens.ParseRefresh(refreshToken)
	if err != nil {
		s.log.Info().Err(err).Msg("error while parsing refresh token")
		return nil, erix.Wrap(err, erix.CodeUnauthorized, tokens.ErrInvalidToken)
	}

	log := s.log.With().Uint64("jwt_id", refresh.Id).Logger()

	token, err := s.tokenRepo.Token(ctx, refresh.Id)
	switch {
	case errors.Is(err, repo.ErrEmptyResult):
		log.Info().Err(err).Msg("empty result while getting token")
		return nil, erix.Wrap(err, erix.CodeUnauthorized, tokens.ErrExpired)
	case err != nil:
		log.Error().Err(err).Msg("error while getting token")
		return nil, erix.Wrap(err, erix.CodeInternalServerError, ErrInternal)
	}

	if token.Rotation != refresh.Rotation {
		log.Info().Msg("invalid rotation number. Invalidating")
		if err = s.tokenRepo.DeleteTokens(ctx, token.UserId, token.AppId); err != nil {
			log.Error().Err(err).Msg("could not invalidate sus tokens")
			return nil, erix.Wrap(err, erix.CodeUnauthorized, ErrInvalidTokenRotation)
		}
		return nil, erix.Wrap(ErrInvalidTokenRotation, erix.CodeUnauthorized, ErrInvalidTokenRotation)
	}

	rotatedToken := models.Token{
		Id:        refresh.Id,
		UserId:    token.UserId,
		AppId:     token.AppId,
		Rotation:  token.Rotation + 1,
		RotatedAt: time.Now(),
		Platform:  token.Platform,
		Agent:     token.Agent,
	}
	err = s.tokenRepo.UpdateToken(ctx, rotatedToken)
	if err != nil {
		log.Error().Err(err).Msg("could not update token rotation")
		return nil, erix.Wrap(err, erix.CodeUnauthorized, ErrInternal)
	}

	newAccess, err := s.tokens.SignAccess(rotatedToken.UserId, rotatedToken.AppId, "0", "*")
	if err != nil {
		log.Error().Err(err).Msg("error while signing access token")
		return nil, erix.Wrap(err, erix.CodeInternalServerError, ErrInternal)
	}

	newRefresh, err := s.tokens.SignRefresh(rotatedToken.UserId, rotatedToken.AppId, rotatedToken.Rotation, rotatedToken.Id, "0")
	if err != nil {
		log.Error().Err(err).Msg("error while signing refresh token")
		return nil, erix.Wrap(err, erix.CodeInternalServerError, ErrInternal)
	}

	return &Tokens{
		Refresh: newRefresh,
		Access:  newAccess,
	}, nil
}
