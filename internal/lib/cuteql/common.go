package cuteql

import (
	"errors"

	"github.com/Masterminds/squirrel"
	"github.com/Onnywrite/ssonny/internal/storage/repo"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/rotisserie/eris"
)

// -----------------------------------------------
//
// Building squirrel
//
// -----------------------------------------------

func buildSquirrel(builder squirrel.Sqlizer) (string, []any, error) {
	switch sq := builder.(type) {
	case squirrel.SelectBuilder:
		query, args, err := sq.PlaceholderFormat(squirrel.Dollar).ToSql()
		if err != nil {
			return "", nil, eris.Wrap(repo.ErrInternal, "could not build squirrel query: "+err.Error())
		}
		return query, args, nil
	case squirrel.UpdateBuilder:
		query, args, err := sq.PlaceholderFormat(squirrel.Dollar).ToSql()
		if err != nil {
			return "", nil, eris.Wrap(repo.ErrInternal, "could not build squirrel query: "+err.Error())
		}
		return query, args, nil
	case squirrel.DeleteBuilder:
		query, args, err := sq.PlaceholderFormat(squirrel.Dollar).ToSql()
		if err != nil {
			return "", nil, eris.Wrap(repo.ErrInternal, "could not build squirrel query: "+err.Error())
		}
		return query, args, nil
	case squirrel.InsertBuilder:
		query, args, err := sq.PlaceholderFormat(squirrel.Dollar).ToSql()
		if err != nil {
			return "", nil, eris.Wrap(repo.ErrInternal, "could not build squirrel query: "+err.Error())
		}
		return query, args, nil
	default:
		return "", nil, eris.Wrap(repo.ErrInternal, "could not build squirrel query: unsupported builder type")
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
