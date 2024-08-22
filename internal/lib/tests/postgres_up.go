package tests

import (
	"context"
	"time"

	"github.com/Onnywrite/ssonny/internal/storage/postgres"

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

type Terminator interface {
	Terminate(ctx context.Context) error
}

func PostgresUp(s *suite.Suite) (*postgres.PgStorage, Terminator) {
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

type PostgresSuite struct {
	suite.Suite
	pgcontainer Terminator

	Pg *postgres.PgStorage
}

func (pgs *PostgresSuite) SetupSuite() {
	pgs.Pg, pgs.pgcontainer = PostgresUp(&pgs.Suite)
}

func (pgs *PostgresSuite) TearDownSuite() {
	err := pgs.Pg.Disconnect()
	pgs.Require().NoError(err)
	err = pgs.pgcontainer.Terminate(context.Background())
	pgs.Require().NoError(err)
}
