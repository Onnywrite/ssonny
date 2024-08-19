package storage

import "github.com/Onnywrite/ssonny/internal/storage/postgres"

type Storage struct {
	*postgres.PgStorage
}

func New(postgresConn string) (*Storage, error) {
	pg, err := postgres.New(postgresConn)
	if err != nil {
		return nil, err
	}

	return &Storage{
		PgStorage: pg,
	}, nil
}
