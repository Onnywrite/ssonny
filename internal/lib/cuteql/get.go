package cuteql

import (
	"context"
	"errors"

	"github.com/Onnywrite/ssonny/internal/storage/repo"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackskj/carta"
	"github.com/jmoiron/sqlx"
	"github.com/rotisserie/eris"
)

func GetNamed[TArg any, T any](ctx context.Context,
	db *sqlx.DB,
	t *sqlx.Tx,
	namedQuery string,
	arg TArg) (*T, *sqlx.Tx, error) {
	query, args, err := sqlx.BindNamed(sqlx.DOLLAR, namedQuery, arg)
	if err != nil {
		return nil, nil, eris.Wrap(repo.ErrInternal, "could not bind named query: "+err.Error())
	}

	return Get[T](ctx, db, t, query, args...)
}

func Get[T any](ctx context.Context,
	db *sqlx.DB,
	t *sqlx.Tx,
	query string,
	args ...any) (obj *T, tx *sqlx.Tx, err error) {
	if t == nil {
		tx, err = db.BeginTxx(ctx, nil)
		if err != nil {
			return nil, nil, eris.Wrap(repo.ErrInternal, "could not begin tx: "+err.Error())
		}
	} else {
		tx = t
	}

	stmt, err := tx.PrepareContext(ctx, query)
	if err != nil {
		return nil, nil, eris.Wrap(repo.ErrInternal, "could not prepare statement: "+err.Error())
	}
	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx, args...)
	if err != nil {
		return nil, nil, eris.Wrap(mapError(err), "could not execute statement: "+err.Error())
	}
	defer rows.Close()

	obj = new(T)
	err = carta.Map(rows, obj)
	if err != nil {
		return nil, nil, eris.Wrap(repo.ErrInternal, "could not map rows: "+err.Error())
	}

	return
}

// copied from https://github.com/jackc/pgerrcode/blob/master/errcode.go
const (
	notNullViolation    = "23502"
	foreignKeyViolation = "23503"
	uniqueViolation     = "23505"
	checkViolation      = "23514"
)

// copied from database/sql package
const (
	sqlErrNooRows = "sql: no rows in result set"
)

var errorsMap = map[string]error{
	notNullViolation:    repo.ErrNull,
	foreignKeyViolation: repo.ErrFK,
	uniqueViolation:     repo.ErrUnique,
	checkViolation:      repo.ErrChecked,
	sqlErrNooRows:       repo.ErrEmptyResult,
}

func mapError(err error) error {
	pgErr := &pgconn.PgError{}
	var stringErr string
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
