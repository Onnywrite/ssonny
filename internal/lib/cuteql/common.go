package cuteql

import (
	"errors"

	"github.com/Masterminds/squirrel"
	"github.com/Onnywrite/ssonny/internal/storage/repo"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jmoiron/sqlx"
	"github.com/rotisserie/eris"
)

// -----------------------------------------------
//
// Commit util
//
// -----------------------------------------------

func Commit(tx *sqlx.Tx) error {
	if err := tx.Commit(); err != nil {
		return eris.Wrap(repo.ErrInternal, "could not commit tx: "+err.Error())
	}

	return nil
}

// -----------------------------------------------
//
// Building squirrel
//
// -----------------------------------------------

func buildSquirrel(builder squirrel.Sqlizer) (string, []any, error) {
	switch squir := builder.(type) {
	case squirrel.SelectBuilder:
		query, args, err := squir.PlaceholderFormat(squirrel.Dollar).ToSql()
		if err != nil {
			return "", nil, eris.Wrap(repo.ErrInternal, "could not build squirrel query: "+err.Error())
		}

		return query, args, nil
	case squirrel.UpdateBuilder:
		query, args, err := squir.PlaceholderFormat(squirrel.Dollar).ToSql()
		if err != nil {
			return "", nil, eris.Wrap(repo.ErrInternal, "could not build squirrel query: "+err.Error())
		}

		return query, args, nil
	case squirrel.DeleteBuilder:
		query, args, err := squir.PlaceholderFormat(squirrel.Dollar).ToSql()
		if err != nil {
			return "", nil, eris.Wrap(repo.ErrInternal, "could not build squirrel query: "+err.Error())
		}

		return query, args, nil
	case squirrel.InsertBuilder:
		query, args, err := squir.PlaceholderFormat(squirrel.Dollar).ToSql()
		if err != nil {
			return "", nil, eris.Wrap(repo.ErrInternal, "could not build squirrel query: "+err.Error())
		}
		return query, args, nil
	default:
		return "", nil, eris.Wrap(
			repo.ErrInternal,
			"could not build squirrel query: unsupported builder type")
	}
}

// -----------------------------------------------
//
// Errors mapping
//
// -----------------------------------------------

// copied from https://github.com/jackc/pgerrcode/blob/master/errcode.go
const (
	notNullViolation    = "23502"
	foreignKeyViolation = "23503"
	uniqueViolation     = "23505"
	checkViolation      = "23514"
)

// copied from database/sql package.
const (
	sqlErrNooRows = "sql: no rows in result set"
)

// nolint: gochecknoglobals
var errorsMap = map[string]error{
	notNullViolation:    repo.ErrNull,
	foreignKeyViolation: repo.ErrFK,
	uniqueViolation:     repo.ErrUnique,
	checkViolation:      repo.ErrChecked,
	sqlErrNooRows:       repo.ErrEmptyResult,
}

func mapError(err error) error {
	var (
		pgErr     = new(pgconn.PgError)
		stringErr string
	)

	if errors.As(err, &pgErr) {
		stringErr = pgErr.Code
	} else {
		stringErr = err.Error()
	}

	doneErr, ok := errorsMap[stringErr]
	if !ok {
		doneErr = repo.ErrInternal
	}

	return doneErr
}
