package cuteql

import (
	"context"
	"fmt"
)

// Execute executes a plain text query without returning any row in a transaction
// tx, which must be NOT nil (else panic).
// It can return [ErrNull], [ErrFK], [ErrUnique], [ErrChecked], [ErrEmptyResult]
// and any other err, so always check it on nil.
//
// It does not commit transaction, but it will rollback tx, if not nil err returned.
func Execute(ctx context.Context, opts ...SqlOpt) error {
	params, err := fromOpts(ctx, opts...)
	if err != nil {
		return err
	}

	defer finishTransaction(&err, &params)

	stmt, err := params.Tx.Preparex(params.Query)
	if err != nil {
		return fmt.Errorf("could not prepare statement: %w", err)
	}
	defer stmt.Close()

	for _, arg := range params.Args {
		_, err = stmt.ExecContext(ctx, arg...)
		if err != nil {
			return fmt.Errorf("could not execute statement: %w", mapError(err))
		}
	}

	return nil
}
