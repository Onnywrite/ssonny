package postgres_test

import (
	"context"
	"testing"
	"time"

	"github.com/Onnywrite/ssonny/internal/domain/models"
	"github.com/Onnywrite/ssonny/internal/storage/postgres"
	"github.com/Onnywrite/ssonny/internal/storage/repo"
	"github.com/brianvoe/gofakeit/v7"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
	pg "github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

const (
	dbUser           = "user"
	dbDatabase       = "database"
	dbPassword       = "password"
	dbImage          = "postgres:16.2-alpine3.19"
	dbMigrationsPath = "file://../../../migrations"
)

type SaveUserSuite struct {
	suite.Suite
	_pgcontainer *pg.PostgresContainer

	pg *postgres.PgStorage
}

func (s *SaveUserSuite) SetupSuite() {
	s.pg, s._pgcontainer = postgresUp(&s.Suite)
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
		Nickname:     ptr(gofakeit.Username()),
		Email:        gofakeit.Email(),
		IsVerified:   gofakeit.Bool(),
		Gender:       ptr(gofakeit.Gender()),
		PasswordHash: ptr(gofakeit.Password(true, true, true, false, false, 60)),
		Birthday:     ptr(time.Date(2024, time.August, 1, 0, 0, 0, 0, time.UTC)),
	}
}

func ptr[T any](t T) *T {
	return &t
}

func postgresUp(s *suite.Suite) (*postgres.PgStorage, *pg.PostgresContainer) {
	ctx := context.Background()
	container, err := pg.Run(ctx,
		dbImage,
		pg.WithUsername(dbUser),
		pg.WithDatabase(dbDatabase),
		pg.WithPassword(dbPassword),
		testcontainers.WithWaitStrategyAndDeadline(time.Second*10,
			wait.ForLog("database system is ready to accept connections").WithOccurrence(2),
		))
	s.Require().NoError(err)

	conn, err := container.ConnectionString(ctx, "sslmode=disable")
	s.Require().NoError(err)

	m, err := migrate.New(dbMigrationsPath, conn)
	s.Require().NoError(err)
	err = m.Up()
	s.Require().NoError(err)

	pg, err := postgres.New(conn)
	s.Require().NoError(err)

	return pg, container
}
