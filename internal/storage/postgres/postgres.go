package postgres

import (
	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/jmoiron/sqlx"
)

type PgStorage struct {
	db *sqlx.DB
}

func New(conn string) (*PgStorage, error) {
	db, err := sqlx.Connect("pgx", conn)
	if err != nil {
		return nil, err
	}

	return &PgStorage{
		db: db,
	}, nil
}

func (pg *PgStorage) Disconnect() error {
	return pg.db.Close()
}
