package postgres_test

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/Onnywrite/ssonny/internal/domain/models"
	"github.com/Onnywrite/ssonny/internal/lib/tests"
	"github.com/Onnywrite/ssonny/internal/storage/repo"
	"github.com/brianvoe/gofakeit/v7"
)

type SaveUserSuite struct {
	tests.PostgresSuite
}

func (s *SaveUserSuite) SetupTest() {
	s.Pg.TruncateTableUsers(context.Background())
}

func (s *SaveUserSuite) TestHappyPath() {
	user := validUser()
	saved, tx, err := s.Pg.SaveUser(context.Background(), user)
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
	u, err := s.Pg.UserById(context.Background(), saved.Id)
	if s.NoError(err) {
		s.Equal(*saved, *u)
	}
}

func (s *SaveUserSuite) TestMyError() {
	user := validUser()
	_, tx, err := s.Pg.SaveUser(context.Background(), user)
	if s.NoError(err) {
		err = tx.Commit()
		s.NoError(err)
	}

	_, _, err = s.Pg.SaveUser(context.Background(), user)
	s.ErrorIs(err, repo.ErrUnique)
}

type UpdateUserSuite struct {
	tests.PostgresSuite
	user models.User
}

func (s *UpdateUserSuite) SetupTest() {
	s.Pg.TruncateTableUsers(context.Background())
	saved, tx, err := s.Pg.SaveUser(context.Background(), validUser())
	if s.NoError(err) {
		err = tx.Commit()
		s.NoError(err)
	}
	s.user = *saved
}

func (s *UpdateUserSuite) TestHappyPath() {
	ctx, c := context.WithTimeout(context.Background(), time.Second)
	defer c()

	err := s.Pg.UpdateUser(ctx, s.user.Id, map[string]any{
		"user_nickname": "_new_nickname",
		"user_email":    "_new_email@golang.test",
		"user_birthday": time.Date(2000, time.August, 1, 0, 0, 0, 0, time.UTC),
	})
	if s.NoError(err) {
		u, err := s.Pg.UserById(ctx, s.user.Id)
		if s.NoError(err) {
			s.user.Nickname = tests.Ptr("_new_nickname")
			s.user.Email = "_new_email@golang.test"
			s.user.Birthday = tests.Ptr(time.Date(2000, time.August, 1, 0, 0, 0, 0, time.UTC))
			s.Equal(s.user, *u)
		}
	}
}

func (s *UpdateUserSuite) TestMyErrors() {
	ctx, c := context.WithTimeout(context.Background(), time.Second)
	defer c()
	err := s.Pg.UpdateUser(ctx, s.user.Id, map[string]any{
		"user_nickname": "_new_nickname",
		"user_email":    nil,
	})
	s.ErrorIs(err, repo.ErrNull)

	anotherUser, tx, err := s.Pg.SaveUser(ctx, validUser())
	if s.NoError(err) {
		err = tx.Commit()
		s.NoError(err)
	}

	err = s.Pg.UpdateUser(ctx, s.user.Id, map[string]any{
		"user_gender": anotherUser.Gender,
	})
	s.NoError(err)

	err = s.Pg.UpdateUser(ctx, s.user.Id, map[string]any{
		"user_nickname": anotherUser.Nickname,
	})
	s.ErrorIs(err, repo.ErrUnique)
}

func (s *UpdateUserSuite) TestInexistentField() {
	ctx, c := context.WithTimeout(context.Background(), time.Second)
	defer c()
	err := s.Pg.UpdateUser(ctx, s.user.Id, map[string]any{
		"user_nickname": "_new_nickname",
		"__email__":     "newemail@golang.test",
	})
	s.ErrorIs(err, repo.ErrInternal)
}

func (s *UpdateUserSuite) TestEmptyFields() {
	ctx, c := context.WithTimeout(context.Background(), time.Second)
	defer c()
	err := s.Pg.UpdateUser(ctx, s.user.Id, map[string]any{})
	s.ErrorIs(err, repo.ErrEmptyResult)
	u, err := s.Pg.UserById(ctx, s.user.Id)
	if s.NoError(err) {
		s.Equal(s.user, *u)
	}
}

type GetUserSuite struct {
	tests.PostgresSuite
	user models.User
}

func (s *GetUserSuite) SetupTest() {
	s.Pg.TruncateTableUsers(context.Background())
	s.user = validUser()
}

func (s *GetUserSuite) TestHappy() {
	ctx, c := context.WithTimeout(context.Background(), time.Second)
	defer c()
	saved, tx, err := s.Pg.SaveUser(ctx, s.user)
	s.Require().NoError(err)
	err = tx.Commit()
	s.Require().NoError(err)

	u, err := s.Pg.UserById(ctx, saved.Id)
	if s.NoError(err) {
		s.Equal(*saved, *u)
	}

	u, err = s.Pg.UserByNickname(ctx, *s.user.Nickname)
	if s.NoError(err) {
		s.Equal(*saved, *u)
	}

	u, err = s.Pg.UserByEmail(ctx, s.user.Email)
	if s.NoError(err) {
		s.Equal(*saved, *u)
	}
}

func (s *GetUserSuite) TestEmptyResult() {
	ctx, c := context.WithTimeout(context.Background(), time.Second)
	defer c()
	_, err := s.Pg.UserByEmail(ctx, s.user.Email)
	s.ErrorIs(err, repo.ErrEmptyResult)
	_, err = s.Pg.UserByNickname(ctx, *s.user.Nickname)
	s.ErrorIs(err, repo.ErrEmptyResult)
	_, err = s.Pg.UserById(ctx, s.user.Id)
	s.ErrorIs(err, repo.ErrEmptyResult)
}

func TestAllUser(t *testing.T) {
	wg := sync.WaitGroup{}
	tests.RunSuitsParallel(&wg, t, new(SaveUserSuite), new(UpdateUserSuite), new(GetUserSuite))
	wg.Wait()
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
