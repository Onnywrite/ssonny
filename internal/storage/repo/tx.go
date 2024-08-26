package repo

import (
	"context"

	"github.com/jmoiron/sqlx"
)

type Transactor interface {
	Rollback() error
	Commit() error
}

type TransactionBeginner interface {
	BeginTransaction(context.Context) (Transactor, error)
}

func WithTransactor(ctx context.Context, tx Transactor) context.Context {
	return context.WithValue(ctx, txKey, tx)
}

type txKeyStruct struct{}

var txKey = txKeyStruct{}

func EjectSqlxTransaction(ctx context.Context) *sqlx.Tx {
	value := ctx.Value(txKey)
	if value == nil {
		return nil
	}

	transactor, ok := value.(Transactor)
	if !ok {
		return nil
	}

	if sqlxTx, ok := transactor.(*sqlx.Tx); ok {
		return sqlxTx
	}

	return nil
}
