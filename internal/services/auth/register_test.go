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
	"github.com/Onnywrite/ssonny/internal/services/auth"
	"github.com/Onnywrite/ssonny/internal/storage/repo"
	"github.com/Onnywrite/ssonny/mocks"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type RegisterWithPassword struct {
	suite.Suite
	logger zerolog.Logger
	mu     *mocks.UserRepo
	ms     *mocks.SessionRepo
	me     *mocks.EmailService
	mt     *mocks.Transactor
	s      *auth.Service
	data   auth.RegisterWithPasswordData
}

func (s *RegisterWithPassword) SetupSuite() {
	s.logger = zerolog.New(os.Stderr).Level(zerolog.Disabled)
}

func (s *RegisterWithPassword) SetupTest() {
	s.data = validRegisterWithPasswordData()
	s.mu = mocks.NewUserRepo(s.T())
	s.ms = mocks.NewSessionRepo(s.T())
	s.mt = mocks.NewTransactor(s.T())
	s.me = mocks.NewEmailService(s.T())
	s.s = auth.NewService(&s.logger, s.mu, s.ms, s.me)
}

func (s *RegisterWithPassword) TestHappyPath() {
	s.me.EXPECT().SendVerificationEmail(mock.Anything, mock.Anything).Return(nil)
	s.mt.EXPECT().Commit().Return(nil).Twice()
	s.mt.EXPECT().Rollback().Return(nil).Twice()
	s.mu.EXPECT().SaveUser(mock.Anything, mock.Anything).Return(&models.User{}, s.mt, nil)
	s.ms.EXPECT().SaveSession(mock.Anything, mock.Anything).Return(s.mt, nil)

	ctx := context.Background()
	authUser, err := s.s.RegisterWithPassword(ctx, s.data)
	if s.Nil(err) {
		s.NotNil(authUser.Profile.Id)
	}
}

func (s *RegisterWithPassword) TestNonRepoError() {
	ctx := context.Background()

	s.data.Password = strings.Repeat("123", 25) // length greater than 72 bytes
	_, err := s.s.RegisterWithPassword(ctx, s.data)
	s.NotNil(err)

	s.data.Password = strings.Repeat("жыза", 19) // length greater than 72 bytes
	_, err = s.s.RegisterWithPassword(ctx, s.data)
	s.NotNil(err)
}

func (s *RegisterWithPassword) TestUserRepoError() {
	s.mu.EXPECT().SaveUser(mock.Anything, mock.AnythingOfType("models.User")).Return(nil, nil, repo.ErrUnique).Once()

	ctx := context.Background()
	_, err := s.s.RegisterWithPassword(ctx, s.data)
	if s.NotNil(err) {
		s.ErrorIs(err, auth.ErrUserAlreadyExists)
		s.Equal(erix.CodeConflict, erix.HttpCode(err))
	}
}

func (s *RegisterWithPassword) TestUserRepoCommitError() {
	someCommitError := fmt.Errorf("commit failed")
	s.mt.EXPECT().Commit().Return(someCommitError).Once()
	s.mt.EXPECT().Rollback().Return(nil).Once()
	s.mu.EXPECT().SaveUser(mock.Anything, mock.AnythingOfType("models.User")).Return(&models.User{}, s.mt, nil).Once()
	s.me.EXPECT().SendVerificationEmail(mock.Anything, mock.Anything).Return(nil)

	ctx := context.Background()
	_, err := s.s.RegisterWithPassword(ctx, s.data)
	if s.NotNil(err) {
		s.ErrorIs(err, auth.ErrInternal)
		s.Equal(erix.CodeInternalServerError, erix.HttpCode(err))
	}
}

func (s *RegisterWithPassword) TestSessionRepoError() {
	s.me.EXPECT().SendVerificationEmail(mock.Anything, mock.Anything).Return(nil)
	s.mt.EXPECT().Commit().Return(nil).Once()
	s.mt.EXPECT().Rollback().Return(nil).Once()
	s.ms.EXPECT().SaveSession(mock.Anything, mock.Anything).Return(nil, repo.ErrInternal).Once()
	s.mu.EXPECT().SaveUser(mock.Anything, mock.AnythingOfType("models.User")).Return(&models.User{}, s.mt, nil).Once()

	ctx := context.Background()
	_, err := s.s.RegisterWithPassword(ctx, s.data)
	if s.NotNil(err) {
		s.ErrorIs(err, auth.ErrInternal)
		s.Equal(erix.CodeInternalServerError, erix.HttpCode(err))
	}
}

func (s *RegisterWithPassword) TestSessionRepoCommitError() {
	userTransactor := mocks.NewTransactor(s.T())
	userTransactor.EXPECT().Commit().Return(nil).Once()
	userTransactor.EXPECT().Rollback().Return(nil).Once()

	someCommitError := fmt.Errorf("commit failed")
	s.mt.EXPECT().Commit().Return(someCommitError).Once()
	s.mt.EXPECT().Rollback().Return(nil).Once()
	s.mu.EXPECT().SaveUser(mock.Anything, mock.AnythingOfType("models.User")).Return(&models.User{}, userTransactor, nil).Once()
	s.me.EXPECT().SendVerificationEmail(mock.Anything, mock.Anything).Return(nil)
	s.ms.EXPECT().SaveSession(mock.Anything, mock.Anything).Return(s.mt, nil).Once()

	ctx := context.Background()
	_, err := s.s.RegisterWithPassword(ctx, s.data)
	if s.NotNil(err) {
		s.ErrorIs(err, auth.ErrInternal)
		s.Equal(erix.CodeInternalServerError, erix.HttpCode(err))
	}
}

func (s *RegisterWithPassword) TestEmailUnconfirmedError() {
	s.mt.EXPECT().Rollback().Return(nil).Once()
	s.mu.EXPECT().SaveUser(mock.Anything, mock.AnythingOfType("models.User")).Return(&models.User{}, s.mt, nil).Once()
	someEmailError := fmt.Errorf("email is invalid and does not exists and something went wrong")
	s.me.EXPECT().SendVerificationEmail(mock.Anything, mock.Anything).Return(someEmailError).Once()

	ctx := context.Background()
	_, err := s.s.RegisterWithPassword(ctx, s.data)
	if s.NotNil(err) {
		s.ErrorIs(err, auth.ErrEmailUnverified)
		s.Equal(erix.CodeBadRequest, erix.HttpCode(err))
	}
}

func ptr[T any](a T) *T {
	return &a
}

func validRegisterWithPasswordData() auth.RegisterWithPasswordData {
	var (
		minBirthday = time.Date(1945, time.September, 2, 0, 0, 0, 0, time.UTC)
		maxBirthday = time.Now()
	)
	return auth.RegisterWithPasswordData{
		Nickname: gofakeit.Username(),
		Email:    gofakeit.Email(),
		Gender:   gofakeit.Gender(),
		Birthday: ptr(gofakeit.DateRange(minBirthday, maxBirthday)),
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
