package cuteql

import (
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/Onnywrite/ssonny/internal/storage/repo"
	"github.com/jmoiron/sqlx"
	"github.com/rotisserie/eris"
)

func ExecuteNamed[TArg any](ctx context.Context,
	db *sqlx.DB,
	t *sqlx.Tx,
	namedQuery string,
	arg TArg) (*sqlx.Tx, error) {
	query, args, err := sqlx.BindNamed(sqlx.DOLLAR, namedQuery, arg)
	if err != nil {
		return nil, eris.Wrap(repo.ErrInternal, "could not bind named query: "+err.Error())
	}

	return Execute(ctx, db, t, query, args...)
}

func ExecuteSquirreled(ctx context.Context,
	db *sqlx.DB,
	t *sqlx.Tx,
	builder squirrel.Sqlizer) (*sqlx.Tx, error) {
	query, args, err := buildSquirrel(builder)
	if err != nil {
		return nil, err
	}
	return Execute(ctx, db, t, query, args...)
}

func Execute(ctx context.Context,
	db *sqlx.DB,
	t *sqlx.Tx,
	query string,
	args ...any) (tx *sqlx.Tx, err error) {
	if t == nil {
		tx, err = db.BeginTxx(ctx, nil)
		if err != nil {
			return nil, eris.Wrap(repo.ErrInternal, "could not begin tx: "+err.Error())
		}
	} else {
		tx = t
	}

	stmt, err := tx.PrepareContext(ctx, query)
	if err != nil {
		return nil, eris.Wrap(repo.ErrInternal, "could not prepare statement: "+err.Error())
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, args...)
	if err != nil {
		return nil, eris.Wrap(mapError(err), "could not execute statement: "+err.Error())
	}

	return
}
