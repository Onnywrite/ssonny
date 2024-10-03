package cuteql

import "fmt"

// finishTransaction provides a centralized mechanism for
// gracefully concluding a database transaction.
// It ensures data integrity by rolling back changes if errors occur
// and provides flexibility by allowing conditional commits and rollbacks.
//
// It automatically commits changes if ShouldCommit is true
// and rollbacks if ShouldRollback is true and an error occurred.
//
// It is always being used in the defer statements of every SQL func wrapper.
func finishTransaction(errPtr *error, params *SqlParams) {
	if *errPtr != nil && params.ShouldRollback {
		_ = params.Tx.Rollback()
		return
	}

	if *errPtr != nil || !params.ShouldCommit {
		return
	}

	commitErr := params.Tx.Commit()
	if commitErr != nil {
		*errPtr = fmt.Errorf("could not commit transaction: %w", commitErr)
	}
}
