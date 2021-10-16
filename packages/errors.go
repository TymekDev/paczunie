package packages

import (
	"fmt"
)

// RollbackError is a wrapper on an error for storing an additionall rollback
// error.
type rollbackError struct {
	err         error
	rollbackErr error
}

var _ error = rollbackError{}

func (e rollbackError) Error() string {
	if e.rollbackErr == nil {
		if e.err == nil {
			return ""
		}
		return e.err.Error()
	}
	const format = "%s (additionally, rollback failed: %s)"
	return fmt.Sprintf(format, e.err, e.rollbackErr)
}

// WithRollback attaches rollbackErr error to err. If rollback is nil, then err
// is returned.
func withRollback(err, rollbackErr error) error {
	if rollbackErr == nil {
		return err
	}
	return rollbackError{
		err:         err,
		rollbackErr: rollbackErr,
	}
}
