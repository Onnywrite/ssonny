package postgres_test

import (
	"context"
	"reflect"
	"sync"
	"testing"
	"time"

	"github.com/Onnywrite/ssonny/internal/domain/models"
	"github.com/Onnywrite/ssonny/internal/lib/tests"
	"github.com/Onnywrite/ssonny/internal/storage/repo"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/google/uuid"
)

type SaveTokenSuite struct {
	tokensSuiteBase
}

func (s *SaveTokenSuite) TestHappyPath() {
	id, tx, err := s.Pg.SaveToken(s.ctx, s.token)
	s.Require().NoError(err)
	err = tx.Commit()
	s.Require().NoError(err)

	tkn, err := s.Pg.Token(s.ctx, id)
	if s.NoError(err) {
		s.token.Id = id
		// postgres stores time less precise(
		s.token.RotatedAt = tkn.RotatedAt
		s.Equal(s.token, *tkn)
	}
}

func (s *SaveTokenSuite) TestInexistentUser() {
	s.token.UserId = uuid.New()
	_, _, err := s.Pg.SaveToken(s.ctx, s.token)
	s.ErrorIs(err, repo.ErrFK)
}

//	func (s *SaveTokenSuite) TestInexistentApp() {
//		s.token.AppId = uint64(gofakeit.Int64())
//		_, _, err := s.pg.SaveToken(context.Background(), s.token)
//		s.ErrorIs(err, repo.ErrFK)
//	}

func (s *SaveTokenSuite) TestDuplicatingToken() {
	id, tx, err := s.Pg.SaveToken(s.ctx, s.token)
	s.Require().NoError(err)
	err = tx.Commit()
	s.Require().NoError(err)

	tkn, err := s.Pg.Token(s.ctx, id)
	if s.NoError(err) {
		s.token.Id = id
		// postgres stores time less precise(
		s.token.RotatedAt = tkn.RotatedAt
		s.Equal(s.token, *tkn)
	}

	id, tx, err = s.Pg.SaveToken(s.ctx, s.token)
	s.Require().NoError(err)
	err = tx.Commit()
	s.Require().NoError(err)

	tkn, err = s.Pg.Token(s.ctx, id)
	if s.NoError(err) {
		s.token.Id = id
		// postgres stores time less precise(
		s.token.RotatedAt = tkn.RotatedAt
		s.Equal(s.token, *tkn)
	}
}

type UpdateTokenSuite struct {
	tokensSuiteBase
}

func (s *UpdateTokenSuite) SetupTest() {
	s.tokensSuiteBase.SetupTest()

	id, tx, err := s.Pg.SaveToken(s.ctx, s.token)
	s.Require().NoError(err)
	err = tx.Commit()
	s.Require().NoError(err)
	s.token.Id = id
}

func (s *UpdateTokenSuite) TestHappyPath() {
	s.token.Rotation = uint64(gofakeit.Int64())
	s.token.Platform = gofakeit.AppName()

	err := s.Pg.UpdateToken(s.ctx, s.token.Id, map[string]any{
		"token_rotation": s.token.Rotation,
		"token_platform": s.token.Platform,
	})
	s.Require().NoError(err)

	tkn, err := s.Pg.Token(s.ctx, s.token.Id)
	if s.NoError(err) {
		// postgres stores time less precise(
		s.token.RotatedAt = tkn.RotatedAt
		s.Equal(s.token, *tkn)
	}
}

func (s *UpdateTokenSuite) TestMyErrors() {
	err := s.Pg.UpdateToken(s.ctx, s.token.Id, map[string]any{
		"token_rotation": nil,
		"token_platform": gofakeit.AppName(),
	})
	s.ErrorIs(err, repo.ErrNull)

	err = s.Pg.UpdateToken(s.ctx, s.token.Id, map[string]any{})
	s.ErrorIs(err, repo.ErrEmptyResult)

}

func (s *UpdateTokenSuite) TestRestrictedIdField() {
	ctx, c := context.WithTimeout(context.Background(), time.Second)
	defer c()

	fieldId, ok := reflect.TypeFor[models.Token]().FieldByName("Id")
	s.Require().True(ok)
	fieldIdName := fieldId.Tag.Get("db")
	s.Require().NotEmpty(fieldIdName)

	err := s.Pg.UpdateToken(ctx, s.token.Id, map[string]any{
		fieldIdName: uint64(gofakeit.Int64()),
	})
	s.ErrorIs(err, repo.ErrInternal)
}

func (s *UpdateTokenSuite) TestInexistentField() {
	ctx, c := context.WithTimeout(context.Background(), time.Second)
	defer c()

	err := s.Pg.UpdateToken(ctx, s.token.Id, map[string]any{
		"__token_rotation__": s.token.Rotation + 1,
		"platform":           gofakeit.AppName(),
	})
	s.ErrorIs(err, repo.ErrInternal)
}

type GetTokenSuite struct {
	tokensSuiteBase
}

func (s *GetTokenSuite) SetupTest() {
	s.tokensSuiteBase.SetupTest()

	id, tx, err := s.Pg.SaveToken(s.ctx, s.token)
	s.Require().NoError(err)
	err = tx.Commit()
	s.Require().NoError(err)
	s.token.Id = id
}

func (s *GetTokenSuite) TestHappyPath() {
	tkn, err := s.Pg.Token(s.ctx, s.token.Id)
	if s.NoError(err) {
		// postgres stores time less precise(
		s.token.RotatedAt = tkn.RotatedAt
		s.Equal(s.token, *tkn)
	}
}

func (s *GetTokenSuite) TestEmptyResult() {
	_, err := s.Pg.Token(s.ctx, uint64(gofakeit.Int64()))
	s.ErrorIs(err, repo.ErrEmptyResult)
}

type DeleteTokensSuite struct {
	tests.PostgresSuite
	_cancel context.CancelFunc

	userIds [2]uuid.UUID
	appIds  [2]uint64
	tokens  [3]models.Token
	ctx     context.Context
}

func (s *DeleteTokensSuite) SetupSuite() {
	s.PostgresSuite.SetupSuite()

	for i := range s.userIds {
		saved, tx, err := s.Pg.SaveUser(context.Background(), validUser())
		s.Require().NoError(err)
		err = tx.Commit()
		s.Require().NoError(err)
		s.userIds[i] = saved.Id
	}

	for i := range s.appIds {
		s.appIds[i] = uint64(gofakeit.Int64())
	}
}

func (s *DeleteTokensSuite) SetupTest() {
	for i := range s.tokens {
		s.tokens[i] = validToken()
		s.tokens[i].UserId = s.userIds[i%2]
		s.tokens[i].AppId = &s.appIds[i%2]

		id, tx, err := s.Pg.SaveToken(context.Background(), s.tokens[i])
		s.Require().NoError(err)
		err = tx.Commit()
		s.Require().NoError(err)
		s.tokens[i].Id = id
	}

	s.ctx, s._cancel = context.WithTimeout(context.Background(), time.Second*2)
}

func (s *DeleteTokensSuite) TearDownTest() {
	err := s.Pg.TruncateTableTokens(s.ctx)
	s.Require().NoError(err)
	s._cancel()
}

func (s *DeleteTokensSuite) TestHappyPath() {
	err := s.Pg.DeleteTokens(s.ctx, s.userIds[0], &s.appIds[0])
	s.Require().NoError(err)

	count, err := s.Pg.CountTokens(s.ctx, s.userIds[0], &s.appIds[0])
	s.Require().NoError(err)
	s.Equal(uint64(0), count)

	_, err = s.Pg.Token(s.ctx, s.tokens[0].Id)
	s.ErrorIs(err, repo.ErrEmptyResult)

	_, err = s.Pg.Token(s.ctx, s.tokens[1].Id)
	s.NoError(err)
	s.Equal(*s.tokens[1].AppId, s.appIds[1])
	s.Equal(s.tokens[1].UserId, s.userIds[1])
}

func (s *DeleteTokensSuite) TestEmptyResult() {
	err := s.Pg.DeleteTokens(s.ctx, uuid.New(), tests.Ptr(uint64(gofakeit.Int64())))
	s.Require().NoError(err)
}

func TestAllToken(t *testing.T) {
	wg := sync.WaitGroup{}
	tests.RunSuitsParallel(&wg, t,
		new(SaveTokenSuite),
		new(UpdateTokenSuite),
		new(GetTokenSuite),
		new(DeleteTokensSuite),
	)
	wg.Wait()
}

type tokensSuiteBase struct {
	tests.PostgresSuite

	userId uuid.UUID
	appId  uint64
	token  models.Token
	ctx    context.Context

	_cancel context.CancelFunc
}

func (s *tokensSuiteBase) SetupSuite() {
	s.PostgresSuite.SetupSuite()

	user, tx, err := s.Pg.SaveUser(context.Background(), validUser())
	s.Require().NoError(err)
	err = tx.Commit()
	s.Require().NoError(err)
	s.userId = user.Id

	s.appId = uint64(gofakeit.Int64())
}

func (s *tokensSuiteBase) SetupTest() {
	s.token = validToken()
	s.token.UserId = s.userId
	s.token.AppId = &s.appId

	s.ctx, s._cancel = context.WithTimeout(context.Background(), time.Second*2)
}

func (s *tokensSuiteBase) TearDownTest() {
	s._cancel()
	err := s.Pg.TruncateTableTokens(context.Background())
	s.Require().NoError(err)
}

func validToken() models.Token {
	return models.Token{
		Rotation:  uint64(gofakeit.Int64()),
		RotatedAt: gofakeit.Date(),
		Platform:  gofakeit.AppName(),
		Agent:     gofakeit.AdjectiveInterrogative(),
	}
}
