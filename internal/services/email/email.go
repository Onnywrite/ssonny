package email

import (
	"context"

	"github.com/rs/zerolog"
)

type FakeEmailService struct {
	log *zerolog.Logger
}

func New(logger *zerolog.Logger) (*FakeEmailService, error) {
	return &FakeEmailService{
		log: logger,
	}, nil
}

func (s *FakeEmailService) SendVerificationEmail(ctx context.Context, data VerificationEmail,
) error {
	s.log.Info().Str("recipient", data.Recipient).
		Str("user_nickname", data.UserNickname).
		Str("token", data.Token).
		Msg("sending verification email..")

	return nil
}
