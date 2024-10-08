package auth_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/Onnywrite/ssonny/internal/domain/models"
	"github.com/Onnywrite/ssonny/internal/lib/erix"
	"github.com/Onnywrite/ssonny/internal/lib/tests"
	"github.com/Onnywrite/ssonny/internal/services/auth"
	"github.com/Onnywrite/ssonny/internal/storage/repo"
	authmocks "github.com/Onnywrite/ssonny/mocks/auth"
	repomocks "github.com/Onnywrite/ssonny/mocks/repo"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"golang.org/x/crypto/bcrypt"
)

type LoginWithPasswordSuite struct {
	suite.Suite
	logger zerolog.Logger

	mu   *authmocks.UserRepo
	mt   *repomocks.Transactor
	mtok *authmocks.TokenRepo
	ms   *authmocks.TokenSigner
	s    *auth.Service
	ctx  context.Context
	data auth.LoginWithPasswordData
	user *models.User
}

func (s *LoginWithPasswordSuite) SetupSuite() {
	s.logger = zerolog.New(os.Stderr).Level(zerolog.Disabled)
}

func (s *LoginWithPasswordSuite) SetupTest() {
	s.mu = authmocks.NewUserRepo(s.T())
	s.mt = repomocks.NewTransactor(s.T())
	s.mtok = authmocks.NewTokenRepo(s.T())
	s.ms = authmocks.NewTokenSigner(s.T())
	s.s = auth.NewService(&s.logger, s.mu, nil, s.mtok, s.ms)
	s.ctx = context.Background()
	s.data = s.validLoginWithPasswordData()
	s.user = s.registeredUser(s.data)
}

func (s *LoginWithPasswordSuite) TestWithEmailAndNickname() {
	s.mt.EXPECT().Commit().Return(nil).Once()
	s.mt.EXPECT().Rollback().Return(nil).Once()
	s.ms.EXPECT().SignRefresh(s.user.Id, (*uint64)(nil), "self", uint64(0), uint64(52)).Return("refresh_token", nil).Once()
	s.mtok.EXPECT().SaveToken(mock.Anything, mock.Anything).Return(52, s.mt, nil).Once()
	s.ms.EXPECT().SignAccess(s.user.Id, (*uint64)(nil), "self", "*").Return("access_token", nil).Once()
	s.mu.EXPECT().UserByEmail(mock.Anything, mock.Anything).Return(s.user, nil).Once()

	_, err := s.s.LoginWithPassword(s.ctx, s.data)
	s.NoError(err)
}

func (s *LoginWithPasswordSuite) TestWithEmail() {
	s.data.Nickname = nil

	s.mt.EXPECT().Commit().Return(nil).Once()
	s.mt.EXPECT().Rollback().Return(nil).Once()
	s.ms.EXPECT().SignRefresh(s.user.Id, (*uint64)(nil), "self", uint64(0), uint64(52)).Return("refresh_token", nil).Once()
	s.mtok.EXPECT().SaveToken(mock.Anything, mock.Anything).Return(52, s.mt, nil).Once()
	s.ms.EXPECT().SignAccess(s.user.Id, (*uint64)(nil), "self", "*").Return("access_token", nil).Once()
	s.mu.EXPECT().UserByEmail(mock.Anything, mock.Anything).Return(s.user, nil).Once()

	_, err := s.s.LoginWithPassword(s.ctx, s.data)
	s.NoError(err)
}

func (s *LoginWithPasswordSuite) TestWithNickname() {
	s.data.Email = nil

	s.mt.EXPECT().Commit().Return(nil).Once()
	s.mt.EXPECT().Rollback().Return(nil).Once()
	s.ms.EXPECT().SignRefresh(s.user.Id, (*uint64)(nil), "self", uint64(0), uint64(52)).Return("refresh_token", nil).Once()
	s.mtok.EXPECT().SaveToken(mock.Anything, mock.Anything).Return(52, s.mt, nil).Once()
	s.ms.EXPECT().SignAccess(s.user.Id, (*uint64)(nil), "self", "*").Return("access_token", nil).Once()
	s.mu.EXPECT().UserByNickname(mock.Anything, mock.Anything).Return(s.user, nil).Once()

	_, err := s.s.LoginWithPassword(s.ctx, s.data)
	s.NoError(err)
}

func (s *LoginWithPasswordSuite) TestWithNothing() {
	s.data.Email = nil
	s.data.Nickname = nil

	_, err := s.s.LoginWithPassword(s.ctx, s.data)
	if s.ErrorIs(err, auth.ErrInvalidData) {
		s.Equal(erix.CodeBadRequest, erix.HttpCode(err))
	}
}

func (s *LoginWithPasswordSuite) TestUserNotFound() {
	s.mu.EXPECT().UserByEmail(mock.Anything, mock.Anything).Return(nil, repo.ErrEmptyResult).Once()

	_, err := s.s.LoginWithPassword(s.ctx, s.data)
	if s.ErrorIs(err, auth.ErrInvalidCredentials) {
		s.Equal(erix.CodeNotFound, erix.HttpCode(err))
	}
}

func (s *LoginWithPasswordSuite) TestUserRepoError() {
	s.mu.EXPECT().UserByEmail(mock.Anything, mock.Anything).Return(nil, repo.ErrInternal).Once()

	_, err := s.s.LoginWithPassword(s.ctx, s.data)
	if s.ErrorIs(err, auth.ErrInternal) {
		s.Equal(erix.CodeInternalServerError, erix.HttpCode(err))
	}
}

func (s *LoginWithPasswordSuite) TestWrongPassword() {
	s.data.Password = "wrongPassword"

	s.mu.EXPECT().UserByEmail(mock.Anything, mock.Anything).Return(s.user, nil).Once()

	_, err := s.s.LoginWithPassword(s.ctx, s.data)
	if s.ErrorIs(err, auth.ErrInvalidCredentials) {
		s.Equal(erix.CodeNotFound, erix.HttpCode(err))
	}
}

func (s *LoginWithPasswordSuite) validLoginWithPasswordData() auth.LoginWithPasswordData {
	return auth.LoginWithPasswordData{
		Email:    tests.Ptr(gofakeit.Email()),
		Nickname: tests.Ptr(gofakeit.Username()),
		Password: gofakeit.Password(true, true, true, true, true, 16),
		UserInfo: auth.UserInfo{
			Platform: gofakeit.AppName(),
			Agent:    gofakeit.AppName(),
		},
	}
}

func (s *LoginWithPasswordSuite) registeredUser(data auth.LoginWithPasswordData) *models.User {
	nick := gofakeit.Username()
	email := gofakeit.Email()
	if data.Email != nil {
		email = *data.Email
	}
	if data.Nickname != nil {
		nick = *data.Nickname
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(data.Password), bcrypt.DefaultCost)
	s.Require().NoError(err)
	return &models.User{
		Id:           uuid.New(),
		Nickname:     &nick,
		Email:        email,
		Verified:     gofakeit.Bool(),
		Gender:       tests.Ptr(gofakeit.Gender()),
		Birthday:     tests.Ptr(gofakeit.DateRange(time.Date(1945, time.September, 2, 0, 0, 0, 0, time.UTC), time.Now())),
		PasswordHash: tests.Ptr(string(hashedPassword)),
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
		DeletedAt:    nil,
	}
}

func TestLoginWithPassword(t *testing.T) {
	suite.Run(t, &LoginWithPasswordSuite{})
}
