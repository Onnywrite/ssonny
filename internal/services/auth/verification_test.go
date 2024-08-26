package auth_test

import (
	"context"
	"os"
	"sync"
	"testing"

	"github.com/Onnywrite/ssonny/internal/domain/models"
	"github.com/Onnywrite/ssonny/internal/lib/tests"
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
	s.s = auth.NewService(&s.logger, s.mu, nil, nil, nil)
	var err error
	s.Require().NoError(err)
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

type IsVerifiedSuite struct {
	suite.Suite
	logger zerolog.Logger
	mu     *authmocks.UserRepo
	s      *auth.Service
}

func (s *IsVerifiedSuite) SetupSuite() {
	s.logger = zerolog.New(os.Stderr).Level(zerolog.Disabled)
}

func (s *IsVerifiedSuite) SetupTest() {
	s.mu = authmocks.NewUserRepo(s.T())
	s.s = auth.NewService(&s.logger, s.mu, nil, nil, nil)
	var err error
	s.Require().NoError(err)
}

func (s *IsVerifiedSuite) TestHappyPath() {
	userId := uuid.New()
	s.mu.EXPECT().UserById(mock.Anything, userId).Return(&models.User{
		Id:         userId,
		IsVerified: true,
	}, nil).Once()

	ctx := context.Background()
	verified, err := s.s.IsVerified(ctx, userId)
	s.NoError(err)
	s.True(verified)
}

func (s *IsVerifiedSuite) TestUserRepoError() {
	userId := uuid.New()
	s.mu.EXPECT().UserById(mock.Anything, userId).Return(nil, repo.ErrEmptyResult).Once()

	ctx := context.Background()
	verified, err := s.s.IsVerified(ctx, userId)
	s.ErrorIs(err, auth.ErrUserNotFound)
	s.False(verified)
}

func TestVarificationAll(t *testing.T) {
	wg := sync.WaitGroup{}
	tests.RunSuitsParallel(t, &wg, new(VerifyEmailSuite), new(IsVerifiedSuite))
	wg.Wait()
}
