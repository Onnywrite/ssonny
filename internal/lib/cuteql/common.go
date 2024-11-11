package cuteql

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jmoiron/sqlx"
)

var (
	ErrNull        = errors.New("null constraint")
	ErrFK          = errors.New("foreign key constraint")
	ErrUnique      = errors.New("unique constraint")
	ErrChecked     = errors.New("check constraint")
	ErrEmptyResult = errors.New("empty result set")
	ErrInternal    = errors.New("internal sql error")
)

func Commit(tx *sqlx.Tx) error {
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("%w: could not commit tx: %w", ErrInternal, err)
	}

	return nil
}

func getTransaction(ctx context.Context, db *sqlx.DB) (*sqlx.Tx, error) {
	var err error

	tx := EjectSqlxTransaction(ctx)
	if tx == nil {
		tx, err = db.BeginTxx(ctx, nil)
	}

	if err != nil {
		return nil, fmt.Errorf("%w: could not begin transaction: %w", ErrInternal, err)
	}

	return tx, nil
}

// copied from https://github.com/jackc/pgerrcode/blob/master/errcode.go
const (
	notNullViolation    = "23502"
	foreignKeyViolation = "23503"
	uniqueViolation     = "23505"
	checkViolation      = "23514"
)

// copied from database/sql package.
const (
	sqlErrNoRows = "sql: no rows in result set"
)

// nolint: gochecknoglobals
var errorsMap = map[string]error{
	notNullViolation:    ErrNull,
	foreignKeyViolation: ErrFK,
	uniqueViolation:     ErrUnique,
	checkViolation:      ErrChecked,
	sqlErrNoRows:        ErrEmptyResult,
}

func mapError(err error) error {
	pgErr := new(pgconn.PgError)
	stringErr := err.Error()

	if errors.As(err, &pgErr) {
		stringErr = pgErr.Code
	}

	doneErr, ok := errorsMap[stringErr]
	if !ok {
		doneErr = ErrInternal
	}

	return doneErr
}
