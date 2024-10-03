package cuteql

import (
	"context"
	"fmt"
)

// add WithResultSize and WithIgnoreEmptyResult
func Query[T any](ctx context.Context, opts ...SqlOpt) ([]T, error) {
	params, err := fromOpts(ctx, opts...)
	if err != nil {
		return nil, err
	}

	defer finishTransaction(&err, &params)

	stmt, err := params.Tx.Preparex(params.Query)
	if err != nil {
		return nil, fmt.Errorf("could not prepare statement: %w", err)
	}
	defer stmt.Close()

	resultSlice := make([]T, 0)
	for _, arg := range params.Args {
		err = stmt.SelectContext(ctx, &resultSlice, arg...)
		if err != nil {
			return nil, fmt.Errorf("could not execute statement: %w", mapError(err))
		}
	}

	if err != nil {
		return nil, fmt.Errorf("could not scan result: %w", err)
	}

	return resultSlice, nil
}
