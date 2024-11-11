package tests

import (
	"context"
	"fmt"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/testcontainers/testcontainers-go"
	pg "github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

const (
	dbUser           = "user"
	dbPassword       = "password"
	dbDatabase       = "database"
	dbImage          = "postgres:16.2-alpine3.19"
	dbMigrationsPath = "file:../migrations"
)

type PostgresConfig struct {
	User           string
	Password       string
	Database       string
	Image          string
	MigrationsPath string
}

type Terminator interface {
	Terminate(ctx context.Context) error
}

func PostgresUpWithConfig(conf PostgresConfig) (Terminator, error) {
	ctx := context.Background()
	container, err := pg.Run(ctx,
		conf.Image,
		pg.WithUsername(conf.User),
		pg.WithDatabase(conf.Database),
		pg.WithPassword(conf.Password),
		testcontainers.WithWaitStrategyAndDeadline(time.Second*10,
			wait.ForLog("database system is ready to accept connections").WithOccurrence(2),
		))
	if err != nil {
		return nil, fmt.Errorf("run pg container: %w", err)
	}

	conn, err := container.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		return nil, fmt.Errorf("get connection string: %w", err)
	}

	m, err := migrate.New(conf.MigrationsPath, conn)
	if err != nil {
		return nil, fmt.Errorf("init migrating: %w", err)
	}

	err = m.Up()
	if err != nil {
		return nil, fmt.Errorf("migrate up: %w", err)
	}

	return container, nil
}
