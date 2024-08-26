package apps_test

import (
	"context"
	"os"
	"strings"
	"testing"

	"github.com/Onnywrite/ssonny/internal/domain/models"
	"github.com/Onnywrite/ssonny/internal/lib/erix"
	"github.com/Onnywrite/ssonny/internal/services/apps"
	"github.com/Onnywrite/ssonny/internal/storage/repo"
	appsmocks "github.com/Onnywrite/ssonny/mocks/apps"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type CreateAppSuite struct {
	suite.Suite
	logger zerolog.Logger

	service     *apps.Service
	repo        *appsmocks.AppRepo
	authService *appsmocks.AuthService
	gen         *appsmocks.RandomGenerator

	ownerId   uuid.UUID
	validData apps.CreateAppData
}

func (s *CreateAppSuite) SetupSuite() {
	s.logger = zerolog.New(os.Stdout).Level(zerolog.Disabled)
}

func (s *CreateAppSuite) SetupTest() {
	s.repo = appsmocks.NewAppRepo(s.T())
	s.gen = appsmocks.NewRandomGenerator(s.T())
	s.authService = appsmocks.NewAuthService(s.T())
	s.service = apps.NewService(&s.logger, s.repo, s.authService, s.gen)
	s.ownerId = uuid.New()
	s.validData = validCreateAppData()
}

func (s *CreateAppSuite) TestHappyPath() {
	s.authService.EXPECT().IsVerified(mock.Anything, s.ownerId).Return(true, nil).Once()
	s.gen.EXPECT().RandomString().Return("randomSecret", nil).Once()
	s.repo.EXPECT().SaveApp(mock.Anything, mock.MatchedBy(func(a models.App) bool {
		return a.Name == s.validData.Name &&
			a.Description == s.validData.Description &&
			a.OwnerId == s.ownerId
	})).Return(nil).Once()

	_, err := s.service.CreateApp(context.Background(), s.ownerId, s.validData)
	s.NoError(err)
}

func (s *CreateAppSuite) TestInvalidData() {
	s.validData.Name = ""

	_, err := s.service.CreateApp(context.Background(), s.ownerId, s.validData)
	if s.ErrorIs(err, apps.ErrInvalidData) {
		s.Equal(erix.CodeBadRequest, erix.HttpCode(err))
	}
}

func (s *CreateAppSuite) TestRandomSecretError() {
	s.authService.EXPECT().IsVerified(mock.Anything, s.ownerId).Return(true, nil).Once()
	s.gen.EXPECT().RandomString().Return("", gofakeit.Error()).Once()

	_, err := s.service.CreateApp(context.Background(), s.ownerId, s.validData)
	if s.ErrorIs(err, apps.ErrInternal) {
		s.Equal(erix.CodeInternalServerError, erix.HttpCode(err))
	}
}

func (s *CreateAppSuite) TestLongSecret() {
	s.authService.EXPECT().IsVerified(mock.Anything, s.ownerId).Return(true, nil).Once()
	s.gen.EXPECT().RandomString().Return(strings.Repeat("secret", 100), nil).Once()
	s.repo.EXPECT().SaveApp(mock.Anything, mock.MatchedBy(func(a models.App) bool {
		return a.Name == s.validData.Name &&
			a.Description == s.validData.Description &&
			a.OwnerId == s.ownerId
	})).Return(nil).Once()

	_, err := s.service.CreateApp(context.Background(), s.ownerId, s.validData)
	s.NoError(err)
}

func (s *CreateAppSuite) TestSaveAppError() {
	s.authService.EXPECT().IsVerified(mock.Anything, s.ownerId).Return(true, nil).Once()
	s.gen.EXPECT().RandomString().Return("randomSecret", nil).Once()
	s.repo.EXPECT().SaveApp(mock.Anything, mock.Anything).Return(repo.ErrUnique).Once()

	_, err := s.service.CreateApp(context.Background(), s.ownerId, s.validData)
	if s.ErrorIs(err, apps.ErrAppAlreadyExists) {
		s.Equal(erix.CodeConflict, erix.HttpCode(err))
	}
}

func (s *CreateAppSuite) TestUserUnverifiedError() {
	s.authService.EXPECT().IsVerified(mock.Anything, s.ownerId).Return(false, nil).Once()

	_, err := s.service.CreateApp(context.Background(), s.ownerId, s.validData)
	if s.ErrorIs(err, apps.ErrUserUnverified) {
		s.Equal(erix.CodeForbidden, erix.HttpCode(err))
	}
}

func (s *CreateAppSuite) TestIsVerifiedError() {
	s.authService.EXPECT().IsVerified(mock.Anything, s.ownerId).Return(false, gofakeit.Error()).Once()

	_, err := s.service.CreateApp(context.Background(), s.ownerId, s.validData)
	s.Error(err)
}

func TestCreateAppSuite(t *testing.T) {
	suite.Run(t, new(CreateAppSuite))
}

func validCreateAppData() apps.CreateAppData {
	return apps.CreateAppData{
		Name:        gofakeit.AppName(),
		Description: gofakeit.Sentence(10),
	}
}
