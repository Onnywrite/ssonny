package auth_test

import (
	"context"
	"os"
	"testing"

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
	logger zerolog.Logger
	mu     *authmocks.UserRepo
	s      *auth.Service
}

func (s *VerifyEmailSuite) SetupSuite() {
	s.logger = zerolog.New(os.Stderr).Level(zerolog.Disabled)
}

func (s *VerifyEmailSuite) SetupTest() {
	s.mu = authmocks.NewUserRepo(s.T())
	s.s = auth.NewService(s.logger, auth.Config{
		UserRepo:     s.mu,
		EmailService: nil,
		TokenRepo:    nil,
		TokensSigner: tokens.NewWithConfig(tokens.Config{}),
	})
}

func (s *VerifyEmailSuite) TestHappyPath() {
	userId := uuid.New()

	s.mu.EXPECT().UpdateUser(mock.Anything, userId, mock.MatchedBy(func(u map[string]any) bool {
		return u["user_verified"].(bool)
	})).Return(nil).Once()

	ctx := context.Background()
	err := s.s.VerifyEmail(ctx, userId)
	s.NoError(err)
}

func (s *VerifyEmailSuite) TestUserUpdateError() {
	userId := uuid.New()
	s.mu.EXPECT().UpdateUser(mock.Anything, userId, mock.Anything).Return(repo.ErrEmptyResult).Once()

	ctx := context.Background()
	err := s.s.VerifyEmail(ctx, userId)
	s.ErrorIs(err, auth.ErrUserNotFound)
}

func TestVerifyEmail(t *testing.T) {
	suite.Run(t, &VerifyEmailSuite{})
}
