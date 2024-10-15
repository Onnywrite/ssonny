package auth_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/Onnywrite/ssonny/internal/domain/models"
	"github.com/Onnywrite/ssonny/internal/lib/erix"
	"github.com/Onnywrite/ssonny/internal/lib/tokens"
	"github.com/Onnywrite/ssonny/internal/services/auth"
	"github.com/Onnywrite/ssonny/internal/storage/repo"
	authmocks "github.com/Onnywrite/ssonny/mocks/auth"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type RefreshSuite struct {
	suite.Suite
	log zerolog.Logger

	ctx   context.Context
	token tokens.Refresh
	mt    *authmocks.TokenRepo
	ms    *authmocks.TokenSigner
	s     *auth.Service
}

func (s *RefreshSuite) SetupSuite() {
	s.log = zerolog.New(os.Stderr).Level(zerolog.Disabled)
}

func (s *RefreshSuite) SetupTest() {
	s.ctx = context.Background()
	s.mt = authmocks.NewTokenRepo(s.T())
	s.ms = authmocks.NewTokenSigner(s.T())
	s.s = auth.NewService(s.log, auth.Config{
		UserRepo:     nil,
		EmailService: nil,
		TokenRepo:    s.mt,
		TokensSigner: s.ms,
	})

	s.token = tokens.Refresh{
		Issuer:          "issuer",
		Subject:         uuid.New(),
		Audience:        nil,
		AuthorizedParty: "party!",
		ExpiresAt:       time.Now().Add(time.Hour).Unix(),
		Id:              1,
		Rotation:        2,
	}
}

func (s *RefreshSuite) TestHappyPath() {
	s.ms.EXPECT().SignAccess(s.token.Subject, s.token.Audience, "self", "*").Return("access_token", nil)
	s.ms.EXPECT().SignRefresh(s.token.Subject, s.token.Audience, "self", s.token.Rotation+1, s.token.Id).Return("refresh_token", nil)
	s.mt.EXPECT().Token(mock.Anything, s.token.Id).Return(&models.Token{
		Id:       1,
		UserId:   s.token.Subject,
		AppId:    s.token.Audience,
		Rotation: 2,
	}, nil).Once()
	s.mt.EXPECT().UpdateToken(mock.Anything, s.token.Id, mock.MatchedBy(func(t map[string]any) bool {
		rot := t["token_rotation"].(uint64)
		return rot == 3
	})).Return(nil).Once()

	_, err := s.s.Refresh(s.ctx, s.token)
	s.NoError(err)
}

func (s *RefreshSuite) TestTokenRepoErrors() {
	s.mt.EXPECT().Token(mock.Anything, s.token.Id).Return(nil, repo.ErrEmptyResult).Once()

	_, err := s.s.Refresh(s.ctx, s.token)
	if s.ErrorIs(err, tokens.ErrExpired) {
		s.Equal(erix.CodeUnauthorized, erix.HttpCode(err))
	}

	s.mt.EXPECT().Token(mock.Anything, s.token.Id).Return(nil, gofakeit.Error()).Once()

	_, err = s.s.Refresh(s.ctx, s.token)
	s.ErrorIs(err, auth.ErrInternal)
}

func (s *RefreshSuite) TestActualRotationGreater() {
	s.mt.EXPECT().Token(mock.Anything, s.token.Id).Return(&models.Token{
		UserId:   s.token.Subject,
		AppId:    s.token.Audience,
		Rotation: 10,
	}, nil).Once()
	s.mt.EXPECT().DeleteTokens(mock.Anything, s.token.Subject, s.token.Audience).Return(nil).Once()

	_, err := s.s.Refresh(s.ctx, s.token)
	if s.ErrorIs(err, auth.ErrInvalidTokenRotation) {
		s.Equal(erix.CodeUnauthorized, erix.HttpCode(err))
	}
}

func (s *RefreshSuite) TestActualRotationLess() {
	s.mt.EXPECT().Token(mock.Anything, s.token.Id).Return(&models.Token{
		UserId:   s.token.Subject,
		AppId:    s.token.Audience,
		Rotation: 1,
	}, nil).Once()
	s.mt.EXPECT().DeleteTokens(mock.Anything, s.token.Subject, s.token.Audience).Return(nil).Once()

	_, err := s.s.Refresh(s.ctx, s.token)
	if s.ErrorIs(err, auth.ErrInvalidTokenRotation) {
		s.Equal(erix.CodeUnauthorized, erix.HttpCode(err))
	}
}

func (s *RefreshSuite) TestDeletionError() {
	s.mt.EXPECT().Token(mock.Anything, s.token.Id).Return(&models.Token{
		UserId:   s.token.Subject,
		AppId:    s.token.Audience,
		Rotation: 222,
	}, nil).Once()
	s.mt.EXPECT().DeleteTokens(mock.Anything, s.token.Subject, s.token.Audience).Return(gofakeit.Error()).Once()

	_, err := s.s.Refresh(s.ctx, s.token)
	if s.ErrorIs(err, auth.ErrInvalidTokenRotation) {
		s.Equal(erix.CodeUnauthorized, erix.HttpCode(err))
	}
}

func (s *RefreshSuite) TestUpdateTokenError() {
	s.mt.EXPECT().Token(mock.Anything, s.token.Id).Return(&models.Token{
		Id:       1,
		Rotation: 2,
	}, nil).Once()
	s.mt.EXPECT().UpdateToken(mock.Anything, s.token.Id, mock.MatchedBy(func(t map[string]any) bool {
		_, hasRotation := t["token_rotation"]
		_, hasRotatedAt := t["token_rotated_at"]
		return hasRotation && hasRotatedAt && len(t) == 2
	})).Return(gofakeit.Error()).Once()

	_, err := s.s.Refresh(s.ctx, s.token)
	if s.ErrorIs(err, auth.ErrInternal) {
		s.Equal(erix.CodeUnauthorized, erix.HttpCode(err))
	}
}

func (s *RefreshSuite) TestSignAccessError() {
	s.mt.EXPECT().Token(mock.Anything, s.token.Id).Return(&models.Token{
		Id:       1,
		UserId:   s.token.Subject,
		AppId:    s.token.Audience,
		Rotation: 2,
	}, nil).Once()
	s.mt.EXPECT().UpdateToken(mock.Anything, s.token.Id, mock.MatchedBy(func(t map[string]any) bool {
		_, hasRotation := t["token_rotation"]
		_, hasRotatedAt := t["token_rotated_at"]
		return hasRotation && hasRotatedAt && len(t) == 2
	})).Return(nil).Once()
	s.ms.EXPECT().SignAccess(s.token.Subject, s.token.Audience, "self", "*").Return("", gofakeit.Error()).Once()

	_, err := s.s.Refresh(s.ctx, s.token)
	if s.ErrorIs(err, auth.ErrInternal) {
		s.Equal(erix.CodeInternalServerError, erix.HttpCode(err))
	}
}

func (s *RefreshSuite) TestSignRefreshError() {
	s.mt.EXPECT().Token(mock.Anything, s.token.Id).Return(&models.Token{
		Id:       1,
		UserId:   s.token.Subject,
		AppId:    s.token.Audience,
		Rotation: 2,
	}, nil).Once()
	s.mt.EXPECT().UpdateToken(mock.Anything, s.token.Id, mock.MatchedBy(func(t map[string]any) bool {
		_, hasRotation := t["token_rotation"]
		_, hasRotatedAt := t["token_rotated_at"]
		return hasRotation && hasRotatedAt && len(t) == 2
	})).Return(nil).Once()
	s.ms.EXPECT().SignAccess(s.token.Subject, s.token.Audience, "self", "*").Return("access_token", nil).Once()
	s.ms.EXPECT().SignRefresh(s.token.Subject, s.token.Audience, "self", s.token.Rotation+1, s.token.Id).Return("", gofakeit.Error()).Once()

	_, err := s.s.Refresh(s.ctx, s.token)
	if s.ErrorIs(err, auth.ErrInternal) {
		s.Equal(erix.CodeInternalServerError, erix.HttpCode(err))
	}
}

func TestRefresh(t *testing.T) {
	suite.Run(t, &RefreshSuite{})
}
