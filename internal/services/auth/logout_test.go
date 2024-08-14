package auth_test

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"os"
	"testing"

	"github.com/Onnywrite/ssonny/internal/lib/erix"
	"github.com/Onnywrite/ssonny/internal/services/auth"
	"github.com/Onnywrite/ssonny/internal/storage/repo"
	"github.com/Onnywrite/ssonny/mocks"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type LogoutSuite struct {
	suite.Suite
	logger zerolog.Logger
	rsaKey *rsa.PrivateKey

	mtok *mocks.TokenRepo
	s    *auth.Service
}

func (s *LogoutSuite) SetupSuite() {
	s.logger = zerolog.New(os.Stderr).Level(zerolog.Disabled)
	rsaKey, err := rsa.GenerateKey(rand.Reader, 1024)
	s.Require().NoError(err)
	s.rsaKey = rsaKey
}

func (s *LogoutSuite) SetupTest() {
	s.mtok = mocks.NewTokenRepo(s.T())
	s.s = auth.NewService(&s.logger, nil, nil, s.mtok, newTokensGen(s.rsaKey))
}

func (s *LogoutSuite) TestHappyPath() {
	s.mtok.EXPECT().DeleteToken(mock.Anything, uint64(1)).Return(nil).Once()

	err := s.s.Logout(context.Background(), 1)
	s.NoError(err)
}

func (s *LogoutSuite) TestTokenRepoEmptyResult() {
	s.mtok.EXPECT().DeleteToken(mock.Anything, uint64(1)).Return(repo.ErrEmptyResult).Once()

	err := s.s.Logout(context.Background(), 1)
	s.NoError(err)
}

func (s *LogoutSuite) TestTokenRepoError() {
	s.mtok.EXPECT().DeleteToken(mock.Anything, uint64(1)).Return(gofakeit.Error()).Once()

	err := s.s.Logout(context.Background(), 1)
	if s.ErrorIs(err, auth.ErrInternal) {
		s.Equal(erix.CodeInternalServerError, erix.HttpCode(err))
	}
}

func TestLogout(t *testing.T) {
	suite.Run(t, &LogoutSuite{})
}
