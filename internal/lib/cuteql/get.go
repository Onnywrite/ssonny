package cuteql

import (
	"context"

	"github.com/Onnywrite/ssonny/internal/storage/repo"

	"github.com/Masterminds/squirrel"
	"github.com/blockloop/scan/v2"
	"github.com/jmoiron/sqlx"
	"github.com/rotisserie/eris"
)

func GetNamed[TArg any, T any](ctx context.Context,
	db *sqlx.DB,
	t *sqlx.Tx,
	namedQuery string,
	arg TArg,
) (*T, *sqlx.Tx, error) {
	query, args, err := sqlx.BindNamed(sqlx.DOLLAR, namedQuery, arg)
	if err != nil {
		return nil, nil, eris.Wrap(repo.ErrInternal, "could not bind named query: "+err.Error())
	}

	return Get[T](ctx, db, t, query, args...)
}

func GetSquirreled[T any](ctx context.Context,
	db *sqlx.DB,
	transaction *sqlx.Tx,
	builder squirrel.Sqlizer,
) (*T, *sqlx.Tx, error) {
	query, args, err := buildSquirrel(builder)
	if err != nil {
		return nil, nil, err
	}

	return Get[T](ctx, db, transaction, query, args...)
}

func Get[T any](ctx context.Context,
	db *sqlx.DB,
	transaction *sqlx.Tx,
	query string,
	args ...any,
) (*T, *sqlx.Tx, error) {
	var (
		tx  *sqlx.Tx
		err error
	)

	if transaction == nil {
		tx, err = db.BeginTxx(ctx, nil)
		if err != nil {
			return nil, nil, eris.Wrap(repo.ErrInternal, "could not begin tx: "+err.Error())
		}
	} else {
		tx = transaction
	}

	stmt, err := tx.PrepareContext(ctx, query)
	if err != nil {
		_ = tx.Rollback()

		return nil, nil, eris.Wrap(repo.ErrInternal, "could not prepare statement: "+err.Error())
	}

	rows, err := stmt.QueryContext(ctx, args...)
	if err != nil {
		_ = tx.Rollback()

		return nil, nil, eris.Wrap(mapError(err), "could not execute statement: "+err.Error())
	}

	obj := new(T)

	err = scan.Row(obj, rows)
	if err != nil {
		_ = tx.Rollback()

		return nil, nil, eris.Wrap(mapError(err), "could not scan result: "+err.Error())
	}

	return obj, tx, nil
}
