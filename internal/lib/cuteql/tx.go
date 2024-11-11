package cuteql

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

func WithTransactor(parent context.Context, tx Transactor) context.Context {
	return context.WithValue(parent, txKey, tx)
}

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

type txKeyStruct struct{}

//nolint: gochecknoglobals
var txKey = txKeyStruct{}
