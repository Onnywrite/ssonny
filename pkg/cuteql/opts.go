package cuteql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
)

type SqlOpt func(context.Context, SqlParams) SqlParams

type SqlParams struct {
	Tx             *sqlx.Tx
	Query          string
	Args           [][]any
	ShouldCommit   bool
	ShouldRollback bool
	Err            error
}

var (
	errOptDbAndTxUsed   = errors.New("using both WithTx and WithDb not allowed")
	errOptDbOrTxNotUsed = errors.New("neither WithTx nor WithDb was used")
	errOptMultipleArgs  = errors.New("args already set, maybe you need WithBatch only")
)

// WithDb is used to provide a sqlx.DB database connection object.
// It also implicitly starts a new database transaction, in which all queries
// will be executed.
// The transaction will be committed or rollbacked automatically.
// You really should NOT change this behaviour by [WithNoRollback].
//
// Do not combine with [WithTx] or [WithNewTx].
func WithDb(db *sqlx.DB) SqlOpt {
	return func(ctx context.Context, sqlp SqlParams) SqlParams {
		if sqlp.Tx != nil {
			sqlp.Err = errOptDbAndTxUsed
			return sqlp
		}

		sqlp.Tx, sqlp.Err = db.BeginTxx(ctx, nil)
		sqlp.ShouldCommit = true

		return sqlp
	}
}

// WithTx is used to provide an existing sqlx.Tx transaction object.
// The transaction will NOT be committed automatically,
// but will be rollbacked if something went wrong.
// You can change this behaviour using [WithCommit] or [WithNoRollback].
//
// Do not combine with [WithDb] or [WithNewTx].
func WithTx(tx *sqlx.Tx) SqlOpt {
	return func(_ context.Context, sqlp SqlParams) SqlParams {
		if sqlp.Tx != nil {
			sqlp.Err = errOptDbAndTxUsed
		}

		sqlp.Tx = tx

		return sqlp
	}
}

func WithNewTx(db *sqlx.DB, isolationLevel sql.IsolationLevel, readOnly bool) SqlOpt {
	return func(ctx context.Context, sqlp SqlParams) SqlParams {
		if sqlp.Tx != nil {
			sqlp.Err = errOptDbAndTxUsed
			return sqlp
		}

		sqlp.Tx, sqlp.Err = db.BeginTxx(ctx, &sql.TxOptions{
			Isolation: isolationLevel,
			ReadOnly:  readOnly,
		})

		return sqlp
	}
}

func WithCommit() SqlOpt {
	return func(_ context.Context, sqlp SqlParams) SqlParams {
		sqlp.ShouldCommit = true

		return sqlp
	}
}

func WithNoRollback() SqlOpt {
	return func(ctx context.Context, sp SqlParams) SqlParams {
		sp.ShouldRollback = false

		return sp
	}
}

func WithQuery(query string) SqlOpt {
	return func(_ context.Context, sqlp SqlParams) SqlParams {
		sqlp.Query = query

		return sqlp
	}
}

type Sqlizer interface {
	ToSql() (string, []interface{}, error)
}

func WithBuilder(sq Sqlizer) SqlOpt {
	return func(_ context.Context, sqlp SqlParams) SqlParams {
		if len(sqlp.Args) > 0 {
			sqlp.Err = errOptMultipleArgs
			return sqlp
		}

		query, args, err := sq.ToSql()
		if err != nil {
			sqlp.Err = err
			return sqlp
		}

		sqlp.Query = query
		sqlp.Args = [][]any{args}

		return sqlp
	}
}

func WithArgs(args ...any) SqlOpt {
	return func(_ context.Context, sqlp SqlParams) SqlParams {
		if len(sqlp.Args) > 0 {
			sqlp.Err = errOptMultipleArgs
			return sqlp
		}

		sqlp.Args = [][]any{args}

		return sqlp
	}
}

// does not work with [WithBuilder].
func WithBatch(args ...[]any) SqlOpt {
	return func(_ context.Context, sqlp SqlParams) SqlParams {
		if len(sqlp.Args) > 0 {
			sqlp.Err = errOptMultipleArgs
			return sqlp
		}

		sqlp.Args = args

		return sqlp
	}
}

// fromOpts processes a variable number of SqlOpt functions to
// configure a default SqlParams value.
//
// It returns an error if no tx provided.
func fromOpts(ctx context.Context, opts ...SqlOpt) (SqlParams, error) {
	params := SqlParams{
		Tx:             nil,
		Query:          "",
		Args:           [][]any{},
		ShouldCommit:   false,
		ShouldRollback: true,
		Err:            nil,
	}

	for _, opt := range opts {
		params = opt(ctx, params)
		if params.Err != nil {
			return params, fmt.Errorf("cuteql opts error: %w", params.Err)
		}
	}

	if params.Tx == nil {
		return params, errOptDbOrTxNotUsed
	}

	return params, nil
}
