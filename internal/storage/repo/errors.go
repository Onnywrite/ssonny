package repo

import "errors"

var (
	ErrEmptyResult      = errors.New("empty result")
	ErrUnique           = errors.New("unique constraint violation")
	ErrChecked          = errors.New("check constraint violation")
	ErrFK               = errors.New("foreign key constraint violation")
	ErrNull             = errors.New("not null constraint violation")
	ErrDataInconsistent = errors.New("data inconsistent")
	ErrInternal         = errors.New("internal error")
)
