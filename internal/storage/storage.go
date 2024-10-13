package storage

import (
	"fmt"

	"github.com/Onnywrite/ssonny/internal/storage/postgres"
	"github.com/Onnywrite/ssonny/internal/storage/redis"
)

type Storage struct {
	*postgres.PgStorage
	*redis.RedisStorage
}

type PostgresConfig struct {
	Host     string
	Port     int
	Username string
	Password string
	Database string
	SslMode  string
}

type RedisConfig struct {
	Host     string
	Port     int
	Username string
	Password string
	Db       int
}

func New(pgconf PostgresConfig, rdconf RedisConfig) (*Storage, error) {
	pgconn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		pgconf.Host, pgconf.Port, pgconf.Username, pgconf.Password, pgconf.Database, pgconf.SslMode)

	postgresDatabase, err := postgres.New(pgconn)
	if err != nil {
		return nil, err
	}

	redisDatabase, err := redis.New(
		fmt.Sprintf("%s:%d", rdconf.Host, rdconf.Port),
		rdconf.Username, rdconf.Password, rdconf.Db)
	if err != nil {
		return nil, err
	}

	return &Storage{
		PgStorage:    postgresDatabase,
		RedisStorage: redisDatabase,
	}, nil
}

func (s *Storage) Disconnect() error {
	defer s.RedisStorage.Disconnect()

	if err := s.PgStorage.Disconnect(); err != nil {
		return err
	}

	return nil
}
