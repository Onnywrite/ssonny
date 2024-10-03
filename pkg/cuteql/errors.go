package cuteql

import (
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
)

var (
	ErrNull        = errors.New("null constraint")
	ErrFK          = errors.New("foreign key constraint")
	ErrUnique      = errors.New("unique constraint")
	ErrChecked     = errors.New("check constraint")
	ErrEmptyResult = errors.New("empty result set")
)

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

var errorsMap = map[string]error{
	notNullViolation:    ErrNull,
	foreignKeyViolation: ErrFK,
	uniqueViolation:     ErrUnique,
	checkViolation:      ErrChecked,
	sqlErrNoRows:        ErrEmptyResult,
}

// mapError acts as a translator.
// It takes a generic error as input and attempts to determine
// if it's a PostgreSQL error. If it is, it maps the PostgreSQL error code to
// a more specific and meaningful custom error defined in this file.
// This makes error handling more structured and easier to understand within your codebase.
//
// The func is internal and used after every query execution.
// It allows you to switch on nice
// [ErrNull], [ErrFK], [ErrUnique], [ErrChecked], [ErrEmptyResult] errors.
func mapError(err error) error {
	pgErr := new(pgconn.PgError)
	stringErr := err.Error()

	if errors.As(err, &pgErr) {
		stringErr = pgErr.Code
	}

	doneErr, ok := errorsMap[stringErr]
	if !ok {
		doneErr = err
	}

	return doneErr
}
