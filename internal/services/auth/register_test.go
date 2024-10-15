package auth_test

import (
	"context"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/Onnywrite/ssonny/internal/domain/models"
	"github.com/Onnywrite/ssonny/internal/lib/erix"
	"github.com/Onnywrite/ssonny/internal/lib/tests"
	"github.com/Onnywrite/ssonny/internal/services/auth"
	"github.com/Onnywrite/ssonny/internal/services/email"
	"github.com/Onnywrite/ssonny/internal/storage/repo"
	authmocks "github.com/Onnywrite/ssonny/mocks/auth"
	repomocks "github.com/Onnywrite/ssonny/mocks/repo"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type RegisterWithPassword struct {
	suite.Suite
	logger zerolog.Logger

	mu   *authmocks.UserRepo
	mt   *repomocks.Transactor
	mtok *authmocks.TokenRepo
	me   *authmocks.EmailService
	ms   *authmocks.TokenSigner
	s    *auth.Service
	data auth.RegisterWithPasswordData
}

func (s *RegisterWithPassword) SetupSuite() {
	s.logger = zerolog.New(os.Stderr).Level(zerolog.Disabled)
}

func (s *RegisterWithPassword) SetupTest() {
	s.data = validRegisterWithPasswordData()
	s.mu = authmocks.NewUserRepo(s.T())
	s.mt = repomocks.NewTransactor(s.T())
	s.mtok = authmocks.NewTokenRepo(s.T())
	s.me = authmocks.NewEmailService(s.T())
	s.ms = authmocks.NewTokenSigner(s.T())
	s.s = auth.NewService(s.logger, s.mu, s.me, s.mtok, s.ms)
}

func (s *RegisterWithPassword) TestHappyPath() {
	s.me.EXPECT().SendVerificationEmail(mock.Anything, mock.MatchedBy(func(message email.VerificationEmail) bool {
		return message.Recipient == s.data.Email && message.UserNickname == *s.data.Nickname
	})).Return(nil)
	s.mt.EXPECT().Commit().Return(nil).Twice()
	s.mt.EXPECT().Rollback().Return(nil).Twice()
	userId := uuid.New()
	s.mu.EXPECT().SaveUser(mock.Anything, mock.Anything).Return(&models.User{
		Id: userId,
	}, s.mt, nil)
	s.ms.EXPECT().SignEmail(userId).Return("email_token", nil).Once()
	s.ms.EXPECT().SignAccess(userId, (*uint64)(nil), "self", "*").Return("access_token", nil).Once()
	s.mtok.EXPECT().SaveToken(mock.Anything, mock.Anything).Return(52, s.mt, nil).Once()
	s.ms.EXPECT().SignRefresh(userId, (*uint64)(nil), "self", uint64(0), uint64(52)).Return("refresh_token", nil).Once()

	ctx := context.Background()
	_, err := s.s.RegisterWithPassword(ctx, s.data)
	s.NoError(err)
}

func (s *RegisterWithPassword) TestSignAccessError() {
	s.me.EXPECT().SendVerificationEmail(mock.Anything, mock.MatchedBy(func(message email.VerificationEmail) bool {
		return message.Recipient == s.data.Email && message.UserNickname == strings.Split(s.data.Email, "@")[0]
	})).Return(nil)
	s.mt.EXPECT().Commit().Return(nil).Once()
	s.mt.EXPECT().Rollback().Return(nil).Once()
	userId := uuid.New()
	s.mu.EXPECT().SaveUser(mock.Anything, mock.Anything).Return(&models.User{
		Id: userId,
	}, s.mt, nil)
	s.ms.EXPECT().SignEmail(userId).Return("email_token", nil).Once()
	s.ms.EXPECT().SignAccess(userId, (*uint64)(nil), "self", "*").Return("", gofakeit.Error()).Once()

	s.data.Nickname = nil
	ctx := context.Background()
	_, err := s.s.RegisterWithPassword(ctx, s.data)
	s.ErrorIs(err, auth.ErrInternal)
}

func (s *RegisterWithPassword) TestSignRefreshError() {
	s.me.EXPECT().SendVerificationEmail(mock.Anything, mock.MatchedBy(func(message email.VerificationEmail) bool {
		return message.Recipient == s.data.Email && message.UserNickname == strings.Split(s.data.Email, "@")[0]
	})).Return(nil)
	s.mt.EXPECT().Commit().Return(nil).Once()
	s.mt.EXPECT().Rollback().Return(nil).Twice()
	userId := uuid.New()
	s.mu.EXPECT().SaveUser(mock.Anything, mock.Anything).Return(&models.User{
		Id: userId,
	}, s.mt, nil)
	s.ms.EXPECT().SignEmail(userId).Return("email_token", nil).Once()
	s.ms.EXPECT().SignAccess(userId, (*uint64)(nil), "self", "*").Return("access_token", nil).Once()
	s.mtok.EXPECT().SaveToken(mock.Anything, mock.Anything).Return(52, s.mt, nil).Once()
	s.ms.EXPECT().SignRefresh(userId, (*uint64)(nil), "self", uint64(0), uint64(52)).Return("", gofakeit.Error()).Once()

	s.data.Nickname = nil
	ctx := context.Background()
	_, err := s.s.RegisterWithPassword(ctx, s.data)
	s.ErrorIs(err, auth.ErrInternal)
}

func (s *RegisterWithPassword) TestNoNickname() {
	s.me.EXPECT().SendVerificationEmail(mock.Anything, mock.MatchedBy(func(message email.VerificationEmail) bool {
		return message.Recipient == s.data.Email && message.UserNickname == strings.Split(s.data.Email, "@")[0]
	})).Return(nil)
	s.mt.EXPECT().Commit().Return(nil).Twice()
	s.mt.EXPECT().Rollback().Return(nil).Twice()
	userId := uuid.New()
	s.mu.EXPECT().SaveUser(mock.Anything, mock.Anything).Return(&models.User{
		Id: userId,
	}, s.mt, nil)
	s.ms.EXPECT().SignEmail(userId).Return("email_token", nil).Once()
	s.ms.EXPECT().SignAccess(userId, (*uint64)(nil), "self", "*").Return("access_token", nil).Once()
	s.mtok.EXPECT().SaveToken(mock.Anything, mock.Anything).Return(52, s.mt, nil).Once()
	s.ms.EXPECT().SignRefresh(userId, (*uint64)(nil), "self", uint64(0), uint64(52)).Return("refresh_token", nil).Once()

	s.data.Nickname = nil
	ctx := context.Background()
	authUser, err := s.s.RegisterWithPassword(ctx, s.data)
	if s.NoError(err) {
		s.NotNil(authUser.Profile.Id)
	}
}

func (s *RegisterWithPassword) TestNonRepoError() {
	ctx := context.Background()

	s.data.Password = strings.Repeat("123", 25) // length greater than 72 bytes
	_, err := s.s.RegisterWithPassword(ctx, s.data)
	s.Error(err)

	s.data.Password = strings.Repeat("жыза", 19) // length greater than 72 bytes
	_, err = s.s.RegisterWithPassword(ctx, s.data)
	s.Error(err)
}

func (s *RegisterWithPassword) TestUserRepoError() {
	s.mu.EXPECT().SaveUser(mock.Anything, mock.AnythingOfType("models.User")).Return(nil, nil, repo.ErrUnique).Once()

	ctx := context.Background()
	_, err := s.s.RegisterWithPassword(ctx, s.data)
	if s.Error(err) {
		s.ErrorIs(err, auth.ErrUserAlreadyExists)
		s.Equal(erix.CodeConflict, erix.HttpCode(err))
	}
}

func (s *RegisterWithPassword) TestUserRepoCommitError() {
	someCommitError := fmt.Errorf("commit failed")
	s.mt.EXPECT().Commit().Return(someCommitError).Once()
	s.mt.EXPECT().Rollback().Return(nil).Once()
	s.mu.EXPECT().SaveUser(mock.Anything, mock.AnythingOfType("models.User")).Return(&models.User{}, s.mt, nil).Once()
	s.ms.EXPECT().SignEmail(mock.AnythingOfType("uuid.UUID")).Return("email_token", nil).Once()
	s.me.EXPECT().SendVerificationEmail(mock.Anything, mock.Anything).Return(nil)

	ctx := context.Background()
	_, err := s.s.RegisterWithPassword(ctx, s.data)
	if s.Error(err) {
		s.ErrorIs(err, auth.ErrInternal)
		s.Equal(erix.CodeInternalServerError, erix.HttpCode(err))
	}
}

func (s *RegisterWithPassword) TestTokenRepoError() {
	s.me.EXPECT().SendVerificationEmail(mock.Anything, mock.Anything).Return(nil)
	s.mt.EXPECT().Commit().Return(nil).Once()
	s.mt.EXPECT().Rollback().Return(nil).Once()
	userId := uuid.New()
	s.mu.EXPECT().SaveUser(mock.Anything, mock.Anything).Return(&models.User{
		Id: userId,
	}, s.mt, nil)
	s.ms.EXPECT().SignEmail(userId).Return("email_token", nil).Once()
	s.ms.EXPECT().SignAccess(userId, (*uint64)(nil), "self", "*").Return("access_token", nil).Once()
	s.mtok.EXPECT().SaveToken(mock.Anything, mock.Anything).Return(0, nil, repo.ErrInternal).Once()

	ctx := context.Background()
	_, err := s.s.RegisterWithPassword(ctx, s.data)
	if s.Error(err) {
		s.ErrorIs(err, auth.ErrInternal)
		s.Equal(erix.CodeInternalServerError, erix.HttpCode(err))
	}
}

func (s *RegisterWithPassword) TestTokenRepoCommitError() {
	userTransactor := repomocks.NewTransactor(s.T())
	userTransactor.EXPECT().Commit().Return(nil).Once()
	userTransactor.EXPECT().Rollback().Return(nil).Once()

	someCommitError := fmt.Errorf("commit failed")
	s.mt.EXPECT().Commit().Return(someCommitError).Once()
	s.mt.EXPECT().Rollback().Return(nil).Once()
	userId := uuid.New()
	s.mu.EXPECT().SaveUser(mock.Anything, mock.Anything).Return(&models.User{
		Id: userId,
	}, userTransactor, nil)
	s.me.EXPECT().SendVerificationEmail(mock.Anything, mock.Anything).Return(nil)
	s.ms.EXPECT().SignEmail(userId).Return("email_token", nil).Once()
	s.ms.EXPECT().SignAccess(userId, (*uint64)(nil), "self", "*").Return("access_token", nil).Once()
	s.mtok.EXPECT().SaveToken(mock.Anything, mock.Anything).Return(52, s.mt, nil).Once()
	s.ms.EXPECT().SignRefresh(userId, (*uint64)(nil), "self", uint64(0), uint64(52)).Return("refresh_token", nil)

	ctx := context.Background()
	_, err := s.s.RegisterWithPassword(ctx, s.data)
	if s.Error(err) {
		s.ErrorIs(err, auth.ErrInternal)
		s.Equal(erix.CodeInternalServerError, erix.HttpCode(err))
	}
}

func (s *RegisterWithPassword) TestEmailUnconfirmedError() {
	s.mt.EXPECT().Rollback().Return(nil).Once()
	userId := uuid.New()
	s.mu.EXPECT().SaveUser(mock.Anything, mock.AnythingOfType("models.User")).Return(&models.User{
		Id: userId,
	}, s.mt, nil).Once()
	s.ms.EXPECT().SignEmail(userId).Return("email_token", nil).Once()
	s.me.EXPECT().SendVerificationEmail(mock.Anything, mock.Anything).Return(gofakeit.Error()).Once()

	ctx := context.Background()
	_, err := s.s.RegisterWithPassword(ctx, s.data)
	if s.Error(err) {
		s.ErrorIs(err, auth.ErrEmailUnverified)
		s.Equal(erix.CodeBadRequest, erix.HttpCode(err))
	}
}

func (s *RegisterWithPassword) TestSignEmailTokenError() {
	s.mt.EXPECT().Rollback().Return(nil).Once()
	userId := uuid.New()
	s.mu.EXPECT().SaveUser(mock.Anything, mock.Anything).Return(&models.User{
		Id: userId,
	}, s.mt, nil)
	s.ms.EXPECT().SignEmail(userId).Return("", gofakeit.Error()).Once()

	ctx := context.Background()
	_, err := s.s.RegisterWithPassword(ctx, s.data)
	if s.Error(err) {
		s.ErrorIs(err, auth.ErrInternal)
		s.Equal(erix.CodeInternalServerError, erix.HttpCode(err))
	}
}

func validRegisterWithPasswordData() auth.RegisterWithPasswordData {
	var (
		minBirthday = time.Date(1945, time.September, 2, 0, 0, 0, 0, time.UTC)
		maxBirthday = time.Now()
	)
	return auth.RegisterWithPasswordData{
		Nickname: tests.Ptr(gofakeit.Username()),
		Email:    gofakeit.Email(),
		Gender:   tests.Ptr(gofakeit.Gender()),
		Birthday: tests.Ptr(gofakeit.DateRange(minBirthday, maxBirthday)),
		Password: gofakeit.Password(true, true, true, true, true, 16),
		UserInfo: auth.UserInfo{
			Platform: gofakeit.AppName(),
			Agent:    gofakeit.AppName(),
		},
	}
}

func TestRegisterWithPassword(t *testing.T) {
	suite.Run(t, &RegisterWithPassword{})
}
