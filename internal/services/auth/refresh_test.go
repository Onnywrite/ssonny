package auth_test

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"os"
	"testing"
	"time"

	"github.com/Onnywrite/ssonny/internal/domain/models"
	"github.com/Onnywrite/ssonny/internal/lib/erix"
	"github.com/Onnywrite/ssonny/internal/lib/tokens"
	"github.com/Onnywrite/ssonny/internal/services/auth"
	"github.com/Onnywrite/ssonny/internal/storage/repo"
	"github.com/Onnywrite/ssonny/mocks"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type RefreshSuite struct {
	suite.Suite
	log    zerolog.Logger
	rsaKey *rsa.PrivateKey

	ctx        context.Context
	validToken string
	mt         *mocks.TokenRepo
	s          *auth.Service
}

func (s *RefreshSuite) SetupSuite() {
	s.log = zerolog.New(os.Stderr).Level(zerolog.Disabled)
	rsaKey, err := rsa.GenerateKey(rand.Reader, 1024)
	s.Require().NoError(err)
	s.rsaKey = rsaKey
}

func (s *RefreshSuite) SetupTest() {
	s.ctx = context.Background()
	tokensGen := tokens.NewWithKeys(
		"",
		time.Hour,
		time.Hour,
		time.Hour,
		&s.rsaKey.PublicKey,
		s.rsaKey)
	validToken, err := tokensGen.SignRefresh(uuid.New(), 2, 0, 1, "0")
	s.Require().NoError(err)
	s.validToken = validToken
	s.mt = mocks.NewTokenRepo(s.T())
	s.s = auth.NewService(&s.log, nil, nil, s.mt, tokensGen)
}

func (s *RefreshSuite) TestHappyPath() {
	s.mt.EXPECT().Token(mock.Anything, uint64(1)).Return(&models.Token{
		Rotation: 2,
	}, nil).Once()
	s.mt.EXPECT().UpdateToken(mock.Anything, mock.MatchedBy(func(t models.Token) bool {
		return t.Rotation == 3
	})).Return(nil).Once()

	_, err := s.s.Refresh(s.ctx, s.validToken)
	s.NoError(err)
}

func (s *RefreshSuite) TestInvalidToken() {
	_, err := s.s.Refresh(s.ctx, "invalidToken")
	if s.ErrorIs(err, tokens.ErrInvalidToken) {
		s.Equal(erix.CodeUnauthorized, erix.HttpCode(err))
	}
}

func (s *RefreshSuite) TestGetTokenErrors() {
	s.mt.EXPECT().Token(mock.Anything, uint64(1)).Return(nil, repo.ErrEmptyResult).Once()

	_, err := s.s.Refresh(s.ctx, s.validToken)
	if s.ErrorIs(err, tokens.ErrExpired) {
		s.Equal(erix.CodeUnauthorized, erix.HttpCode(err))
	}

	s.mt.EXPECT().Token(mock.Anything, uint64(1)).Return(nil, gofakeit.Error()).Once()

	_, err = s.s.Refresh(s.ctx, s.validToken)
	s.Error(err)
}

func (s *RefreshSuite) TestActualRotationGreater() {
	s.mt.EXPECT().Token(mock.Anything, uint64(1)).Return(&models.Token{
		Rotation: 10,
	}, nil).Once()
	s.mt.EXPECT().DeleteTokens(mock.Anything, mock.Anything, mock.Anything).Return(nil).Once()

	_, err := s.s.Refresh(s.ctx, s.validToken)
	if s.ErrorIs(err, auth.ErrInvalidTokenRotation) {
		s.Equal(erix.CodeUnauthorized, erix.HttpCode(err))
	}
}

func (s *RefreshSuite) TestActualRotationLess() {
	s.mt.EXPECT().Token(mock.Anything, uint64(1)).Return(&models.Token{
		Rotation: 1,
	}, nil).Once()
	s.mt.EXPECT().DeleteTokens(mock.Anything, mock.Anything, mock.Anything).Return(nil).Once()

	_, err := s.s.Refresh(s.ctx, s.validToken)
	if s.ErrorIs(err, auth.ErrInvalidTokenRotation) {
		s.Equal(erix.CodeUnauthorized, erix.HttpCode(err))
	}
}

func (s *RefreshSuite) TestDeletionError() {
	s.mt.EXPECT().Token(mock.Anything, uint64(1)).Return(&models.Token{
		Rotation: 222,
	}, nil).Once()
	s.mt.EXPECT().DeleteTokens(mock.Anything, mock.Anything, mock.Anything).Return(gofakeit.Error()).Once()

	_, err := s.s.Refresh(s.ctx, s.validToken)
	if s.ErrorIs(err, auth.ErrInvalidTokenRotation) {
		s.Equal(erix.CodeUnauthorized, erix.HttpCode(err))
	}
}

func (s *RefreshSuite) TestUpdateTokenError() {
	s.mt.EXPECT().Token(mock.Anything, mock.Anything).Return(&models.Token{
		Rotation: 2,
	}, nil).Once()
	s.mt.EXPECT().UpdateToken(mock.Anything, mock.Anything).Return(gofakeit.Error()).Once()

	_, err := s.s.Refresh(s.ctx, s.validToken)
	if s.ErrorIs(err, auth.ErrInternal) {
		s.Equal(erix.CodeUnauthorized, erix.HttpCode(err))
	}
}

func TestRefresh(t *testing.T) {
	suite.Run(t, &RefreshSuite{})
}
