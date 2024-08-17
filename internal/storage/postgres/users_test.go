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
	s.pg.TruncateTableUsers(context.Background())
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
	u, err := s.pg.UserById(context.Background(), saved.Id)
	if s.NoError(err) {
		s.Equal(*saved, *u)
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

type UpdateUserSuite struct {
	suite.Suite
	_pgcontainer tests.Terminator

	pg   *postgres.PgStorage
	user models.User
}

func (s *UpdateUserSuite) SetupSuite() {
	s.pg, s._pgcontainer = tests.PostgresUp(&s.Suite)
}

func (s *UpdateUserSuite) SetupTest() {
	s.pg.TruncateTableUsers(context.Background())
	saved, tx, err := s.pg.SaveUser(context.Background(), validUser())
	if s.NoError(err) {
		err = tx.Commit()
		s.NoError(err)
	}
	s.user = *saved
}

func (s *UpdateUserSuite) TestHappyPath() {
	err := s.pg.UpdateUser(context.Background(), s.user.Id, map[string]any{
		"user_nickname": "_new_nickname",
		"user_email":    "_new_email@golang.test",
		"user_birthday": time.Date(2000, time.August, 1, 0, 0, 0, 0, time.UTC),
	})
	if s.NoError(err) {
		u, err := s.pg.UserById(context.Background(), s.user.Id)
		if s.NoError(err) {
			s.user.Nickname = tests.Ptr("_new_nickname")
			s.user.Email = "_new_email@golang.test"
			s.user.Birthday = tests.Ptr(time.Date(2000, time.August, 1, 0, 0, 0, 0, time.UTC))
			s.Equal(s.user, *u)
		}
	}
}

func (s *UpdateUserSuite) TestMyErrors() {
	err := s.pg.UpdateUser(context.Background(), s.user.Id, map[string]any{
		"user_nickname": "_new_nickname",
		"user_email":    nil,
	})
	s.ErrorIs(err, repo.ErrNull)

	anotherUser, tx, err := s.pg.SaveUser(context.Background(), validUser())
	if s.NoError(err) {
		err = tx.Commit()
		s.NoError(err)
	}

	err = s.pg.UpdateUser(context.Background(), s.user.Id, map[string]any{
		"user_gender": anotherUser.Gender,
	})
	s.NoError(err)

	err = s.pg.UpdateUser(context.Background(), s.user.Id, map[string]any{
		"user_nickname": anotherUser.Nickname,
	})
	s.ErrorIs(err, repo.ErrUnique)
}

func (s *UpdateUserSuite) TestInexistentField() {
	err := s.pg.UpdateUser(context.Background(), s.user.Id, map[string]any{
		"user_nickname": "_new_nickname",
		"__email__":     "newemail@golang.test",
	})
	s.ErrorIs(err, repo.ErrInternal)
}

func (s *UpdateUserSuite) TestEmptyFields() {
	err := s.pg.UpdateUser(context.Background(), s.user.Id, map[string]any{})
	s.ErrorIs(err, repo.ErrEmptyResult)
	u, err := s.pg.UserById(context.Background(), s.user.Id)
	if s.NoError(err) {
		s.Equal(s.user, *u)
	}
}

func (s *UpdateUserSuite) TearDownSuite() {
	err := s.pg.Disconnect()
	s.Require().NoError(err)
	err = s._pgcontainer.Terminate(context.Background())
	s.Require().NoError(err)
}

func TestUpdateUserSuite(t *testing.T) {
	suite.Run(t, new(UpdateUserSuite))
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

type GetUserSuite struct {
	suite.Suite
	_pgcontainer tests.Terminator

	pg   *postgres.PgStorage
	user models.User
}

func (s *GetUserSuite) SetupSuite() {
	s.pg, s._pgcontainer = tests.PostgresUp(&s.Suite)
}

func (s *GetUserSuite) SetupTest() {
	s.pg.TruncateTableUsers(context.Background())
	s.user = validUser()
}

func (s *GetUserSuite) TestHappyPath() {
	ctx, c := context.WithTimeout(context.Background(), time.Second)
	defer c()
	saved, tx, err := s.pg.SaveUser(ctx, s.user)
	s.Require().NoError(err)
	err = tx.Commit()
	s.Require().NoError(err)

	u, err := s.pg.UserByEmail(ctx, s.user.Email)
	if s.NoError(err) {
		s.Equal(*saved, *u)
	}

	u, err = s.pg.UserByNickname(ctx, *s.user.Nickname)
	if s.NoError(err) {
		s.Equal(*saved, *u)
	}

	u, err = s.pg.UserById(ctx, saved.Id)
	if s.NoError(err) {
		s.Equal(*saved, *u)
	}
}

func (s *GetUserSuite) TestEmptyResult() {
	_, err := s.pg.UserByEmail(context.Background(), s.user.Email)
	s.ErrorIs(err, repo.ErrEmptyResult)
	_, err = s.pg.UserByNickname(context.Background(), *s.user.Nickname)
	s.ErrorIs(err, repo.ErrEmptyResult)
	_, err = s.pg.UserById(context.Background(), s.user.Id)
	s.ErrorIs(err, repo.ErrEmptyResult)
}

func TestGetUserSuite(t *testing.T) {
	suite.Run(t, new(GetUserSuite))
}
