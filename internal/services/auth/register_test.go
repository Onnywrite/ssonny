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
	mt     *mocks.Transactor
	s      *auth.Service
}

func (s *RegisterWithPassword) SetupSuite() {
	s.logger = zerolog.New(os.Stderr).With().Logger().Level(zerolog.Disabled)
}

func (s *RegisterWithPassword) SetupTest() {
	s.mu = mocks.NewUserRepo(s.T())
	s.ms = mocks.NewSessionRepo(s.T())
	s.mt = mocks.NewTransactor(s.T())
	s.s = auth.NewService(&s.logger, s.mu, s.ms)
}

func (s *RegisterWithPassword) TestHappyPath() {
	s.mt.EXPECT().Commit().Return(nil)
	s.mt.EXPECT().Rollback().Return(nil)
	s.mu.EXPECT().SaveUser(mock.Anything, mock.Anything).Return(&models.User{}, nil)
	s.ms.EXPECT().SaveSession(mock.Anything, mock.Anything).Return(s.mt, nil)

	ctx := context.Background()
	authUser, err := s.s.RegisterWithPassword(ctx, validRegisterWithPasswordData())
	if s.Nil(err) {
		s.NotNil(authUser.Profile.Id)
	}
}

func (s *RegisterWithPassword) TestNonRepoError() {
	ctx := context.Background()

	data := validRegisterWithPasswordData()
	data.Password = strings.Repeat("123", 25) // length greater than 72 bytes
	_, err := s.s.RegisterWithPassword(ctx, data)
	s.NotNil(err)

	data.Password = strings.Repeat("жыза", 19) // length greater than 72 bytes
	_, err = s.s.RegisterWithPassword(ctx, data)
	s.NotNil(err)
}

func (s *RegisterWithPassword) TestRepoUniqueError() {
	data := validRegisterWithPasswordData()

	s.mt.EXPECT().Commit().Return(nil).Once()
	s.mt.EXPECT().Rollback().Return(nil).Once()
	s.ms.EXPECT().SaveSession(mock.Anything, mock.Anything).Return(s.mt, nil).Once()
	s.mu.EXPECT().SaveUser(mock.Anything, mock.AnythingOfType("models.User")).Return(&models.User{}, nil).Once()
	s.mu.EXPECT().SaveUser(mock.Anything, mock.AnythingOfType("models.User")).Return(nil, repo.ErrUnique).Once()

	ctx := context.Background()
	_, err := s.s.RegisterWithPassword(ctx, data)
	s.Nil(err)
	_, err = s.s.RegisterWithPassword(ctx, data)
	if s.NotNil(err) {
		s.ErrorIs(err, auth.ErrUserAlreadyExists)
		s.Equal(erix.CodeConflict, erix.HttpCode(err))
	}
}

func (s *RegisterWithPassword) TestRepoSessionError() {
	data := validRegisterWithPasswordData()
	s.ms.EXPECT().SaveSession(mock.Anything, mock.Anything).Return(nil, repo.ErrInternal).Once()
	s.mu.EXPECT().SaveUser(mock.Anything, mock.AnythingOfType("models.User")).Return(&models.User{}, nil).Once()

	ctx := context.Background()
	_, err := s.s.RegisterWithPassword(ctx, data)
	if s.NotNil(err) {
		s.ErrorIs(err, auth.ErrInternal)
		s.Equal(erix.CodeInternalServerError, erix.HttpCode(err))
	}
}

func (s *RegisterWithPassword) TestRepoSessionCommitError() {
	data := validRegisterWithPasswordData()
	someCommitError := fmt.Errorf("commit failed")
	s.mt.EXPECT().Commit().Return(someCommitError).Once()
	s.mt.EXPECT().Rollback().Return(nil).Once()
	s.ms.EXPECT().SaveSession(mock.Anything, mock.Anything).Return(s.mt, nil).Once()
	s.mu.EXPECT().SaveUser(mock.Anything, mock.AnythingOfType("models.User")).Return(&models.User{}, nil).Once()

	ctx := context.Background()
	_, err := s.s.RegisterWithPassword(ctx, data)
	if s.NotNil(err) {
		s.ErrorIs(err, auth.ErrInternal)
		s.Equal(erix.CodeInternalServerError, erix.HttpCode(err))
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
