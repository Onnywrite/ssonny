package cuteql

import (
	"context"
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
)

func ExecuteNamed[TArg any](ctx context.Context,
	db *sqlx.DB,
	namedQuery string,
	arg TArg,
) (*sqlx.Tx, error) {
	query, args, err := sqlx.BindNamed(sqlx.DOLLAR, namedQuery, arg)
	if err != nil {
		return nil, fmt.Errorf("%w: could not bind named query: %w", ErrInternal, err)
	}

	return Execute(ctx, db, query, args...)
}

func ExecuteSquirreled(ctx context.Context,
	db *sqlx.DB,
	builder squirrel.Sqlizer,
) (*sqlx.Tx, error) {
	query, args, err := builder.ToSql()
	if err != nil {
		return nil, fmt.Errorf("%w: could not build query: %w", ErrInternal, err)
	}

	return Execute(ctx, db, query, args...)
}

func Execute(ctx context.Context,
	db *sqlx.DB,
	query string,
	args ...any,
) (*sqlx.Tx, error) {
	tx, err := getTransaction(ctx, db)
	if err != nil {
		return nil, err
	}

	stmt, err := tx.PrepareContext(ctx, query)
	if err != nil {
		_ = tx.Rollback()

		return nil, fmt.Errorf("%w: could not prepare statement: %w", ErrInternal, err)
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, args...)
	if err != nil {
		_ = tx.Rollback()

		return nil, fmt.Errorf("%w: could not execute statement: %w", mapError(err), err)
	}

	return tx, nil
}
