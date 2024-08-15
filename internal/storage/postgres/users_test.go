package postgres_test

import (
	"context"
	"testing"
	"time"

	"github.com/Onnywrite/ssonny/internal/domain/models"
	"github.com/Onnywrite/ssonny/internal/lib/tests"
	"github.com/Onnywrite/ssonny/internal/storage/postgres"
	"github.com/Onnywrite/ssonny/internal/storage/repo"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/suite"
)

type SaveUserSuite struct {
	suite.Suite
	_pgcontainer tests.Terminator

	pg *postgres.PgStorage
}

func (s *SaveUserSuite) SetupSuite() {
	s.pg, s._pgcontainer = tests.PostgresUp(&s.Suite)
}

func (s *SaveUserSuite) SetupTest() {
	// TODO: truncate table
}

func (s *SaveUserSuite) TestHappyPath() {
	user := validUser()
	saved, tx, err := s.pg.SaveUser(context.Background(), user)
	if s.NoError(err) {
		err = tx.Commit()
		s.NoError(err)

		s.Equal(user.Nickname, saved.Nickname)
		s.Equal(user.Email, saved.Email)
		s.Equal(user.IsVerified, saved.IsVerified)
		s.Equal(user.Gender, saved.Gender)
		s.Equal(user.PasswordHash, saved.PasswordHash)
		s.Equal(user.Birthday, saved.Birthday)
	}
}

func (s *SaveUserSuite) TestMyError() {
	user := validUser()
	_, tx, err := s.pg.SaveUser(context.Background(), user)
	if s.NoError(err) {
		err = tx.Commit()
		s.NoError(err)
	}

	_, _, err = s.pg.SaveUser(context.Background(), user)
	s.ErrorIs(err, repo.ErrUnique)
}

func (s *SaveUserSuite) TearDownSuite() {
	err := s.pg.Disconnect()
	s.Require().NoError(err)
	err = s._pgcontainer.Terminate(context.Background())
	s.Require().NoError(err)
}

func TestSaveUserSuite(t *testing.T) {
	suite.Run(t, new(SaveUserSuite))
}

func validUser() models.User {
	return models.User{
		Nickname:     tests.Ptr(gofakeit.Username()),
		Email:        gofakeit.Email(),
		IsVerified:   gofakeit.Bool(),
		Gender:       tests.Ptr(gofakeit.Gender()),
		PasswordHash: tests.Ptr(gofakeit.Password(true, true, true, false, false, 60)),
		Birthday:     tests.Ptr(time.Date(2024, time.August, 1, 0, 0, 0, 0, time.UTC)),
	}
}
