package auth_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/Onnywrite/ssonny/internal/lib/isitjwt"
	"github.com/Onnywrite/ssonny/internal/services/auth"
	"github.com/Onnywrite/ssonny/internal/storage/repo"
	"github.com/Onnywrite/ssonny/mocks"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type VerifyEmailSuite struct {
	suite.Suite
	logger     zerolog.Logger
	mu         *mocks.UserRepo
	s          *auth.Service
	validToken string
}

func (s *VerifyEmailSuite) SetupSuite() {
	s.logger = zerolog.New(os.Stderr).Level(zerolog.Disabled)
}

func (s *VerifyEmailSuite) SetupTest() {
	s.mu = mocks.NewUserRepo(s.T())
	s.s = auth.NewService(&s.logger, s.mu, nil, nil)
	var err error
	s.validToken, err = isitjwt.Sign(isitjwt.TODOSecret, uuid.New(), auth.SubjectEmail, time.Hour)
	s.Require().Nil(err)
}

func (s *VerifyEmailSuite) TestHappyPath() {
	s.mu.EXPECT().UpdateUser(mock.Anything, mock.Anything).Return(nil).Once()

	ctx := context.Background()
	err := s.s.VerifyEmail(ctx, s.validToken)
	s.Nil(err)
}

func (s *VerifyEmailSuite) TestVerificationError() {
	ctx := context.Background()
	err := s.s.VerifyEmail(ctx, "invalidToken")
	s.NotNil(err)
}

func (s *VerifyEmailSuite) TestUserUpdateError() {
	s.mu.EXPECT().UpdateUser(mock.Anything, mock.Anything).Return(repo.ErrEmptyResult).Once()

	ctx := context.Background()
	err := s.s.VerifyEmail(ctx, s.validToken)
	s.ErrorIs(err, auth.ErrUserNotFound)
}

func TestVerifyEmail(t *testing.T) {
	suite.Run(t, &VerifyEmailSuite{})
}
