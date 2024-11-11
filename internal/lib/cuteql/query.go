package cuteql

import (
	"context"
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/blockloop/scan/v2"
	"github.com/jmoiron/sqlx"
)

func QueryNamed[TArg any, T any](ctx context.Context,
	db *sqlx.DB,
	namedQuery string,
	arg TArg,
) ([]T, *sqlx.Tx, error) {
	query, args, err := sqlx.BindNamed(sqlx.DOLLAR, namedQuery, arg)
	if err != nil {
		return nil, nil, fmt.Errorf("%w: could not bind named query: %w", ErrInternal, err)
	}

	return Query[T](ctx, db, query, args...)
}

func QuerySquirreled[T any](ctx context.Context,
	db *sqlx.DB,
	builder squirrel.Sqlizer,
) ([]T, *sqlx.Tx, error) {
	query, args, err := builder.ToSql()
	if err != nil {
		return nil, nil, fmt.Errorf("%w: could not build query: %w", ErrInternal, err)
	}

	return Query[T](ctx, db, query, args...)
}

func Query[T any](ctx context.Context,
	db *sqlx.DB,
	query string,
	args ...any,
) ([]T, *sqlx.Tx, error) {
	tx, err := getTransaction(ctx, db)
	if err != nil {
		return nil, nil, err
	}

	stmt, err := tx.PrepareContext(ctx, query)
	if err != nil {
		_ = tx.Rollback()

		return nil, nil, fmt.Errorf("%w: could not prepare statement: %w", ErrInternal, err)
	}

	rows, err := stmt.QueryContext(ctx, args...)
	if err != nil {
		_ = tx.Rollback()

		return nil, nil, fmt.Errorf("%w: could not execute statement: %w", mapError(err), err)
	}

	objs := make([]T, 0, 10)

	err = scan.Rows(&objs, rows)
	if err != nil {
		_ = tx.Rollback()

		return nil, nil, fmt.Errorf("%w: could not scan result: %w", mapError(err), err)
	}

	return objs, tx, nil
}
