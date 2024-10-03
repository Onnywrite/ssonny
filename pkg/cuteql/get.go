package cuteql

import (
	"context"
	"fmt"
)

func Get[T any](ctx context.Context, opts ...SqlOpt) (*T, error) {
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

	obj := new(T)

	for _, arg := range params.Args {
		err = stmt.GetContext(ctx, obj, arg...)
		if err != nil {
			return nil, fmt.Errorf("could not execute statement: %w", mapError(err))
		}
	}

	return obj, nil
}
