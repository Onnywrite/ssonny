package auth

import (
	"context"

	"github.com/Onnywrite/ssonny/internal/lib/erix"
	"github.com/Onnywrite/ssonny/internal/lib/isitjwt"
)

const (
	SubjectEmail = "email"
)

func (s *Service) VerifyEmail(ctx context.Context, token string) error {
	log := s.log.With().Str("token", token).Logger()
	userId, err := isitjwt.Verify(isitjwt.TODOSecret, SubjectEmail, token)
	if err != nil {
		log.Warn().Err(err).Msg("error while verifying email token")
		return erix.Wrap(err, erix.CodeBadRequest, isitjwt.ErrInvalidToken)
	}

	err = s.repo.UpdateUser(ctx, userId, map[string]any{
		"user_verified": true,
	})
	if err != nil {
		return userFailed(&log, err)
	}

	return nil
}
