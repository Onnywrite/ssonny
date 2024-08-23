package auth_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/Onnywrite/ssonny/internal/lib/isitjwt"
	"github.com/Onnywrite/ssonny/internal/lib/tokens"
	"github.com/Onnywrite/ssonny/internal/services/auth"
	"github.com/Onnywrite/ssonny/internal/storage/repo"
	authmocks "github.com/Onnywrite/ssonny/mocks/auth"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type VerifyEmailSuite struct {
	suite.Suite
	logger     zerolog.Logger
	mu         *authmocks.UserRepo
	s          *auth.Service
	validToken string
}

func (s *VerifyEmailSuite) SetupSuite() {
	s.logger = zerolog.New(os.Stderr).Level(zerolog.Disabled)
}

func (s *VerifyEmailSuite) SetupTest() {
	s.mu = authmocks.NewUserRepo(s.T())
	s.s = auth.NewService(&s.logger, s.mu, nil, nil, tokens.New("", "secret", time.Hour, time.Hour, time.Hour))
	var err error
	s.validToken, err = isitjwt.Sign(isitjwt.TODOSecret, uuid.New(), auth.SubjectEmail, time.Hour)
	s.Require().NoError(err)
}

func (s *VerifyEmailSuite) TestHappyPath() {
	s.mu.EXPECT().UpdateUser(mock.Anything, mock.Anything, mock.MatchedBy(func(u map[string]any) bool {
		return u["user_verified"].(bool)
	})).Return(nil).Once()

	ctx := context.Background()
	err := s.s.VerifyEmail(ctx, s.validToken)
	s.NoError(err)
}

func (s *VerifyEmailSuite) TestVerificationError() {
	ctx := context.Background()
	err := s.s.VerifyEmail(ctx, "invalidToken")
	s.Error(err)
}

func (s *VerifyEmailSuite) TestUserUpdateError() {
	s.mu.EXPECT().UpdateUser(mock.Anything, mock.Anything, mock.Anything).Return(repo.ErrEmptyResult).Once()

	ctx := context.Background()
	err := s.s.VerifyEmail(ctx, s.validToken)
	s.ErrorIs(err, auth.ErrUserNotFound)
}

func TestVerifyEmail(t *testing.T) {
	suite.Run(t, &VerifyEmailSuite{})
}
