package main

import (
	"fmt"

	"github.com/pkg/errors"
)

// RollbackError is a wrapper on an error for storing an additionall rollback
// error.
type RollbackError struct {
	Err         error
	RollbackErr error
}

var _ error = RollbackError{}

func (e RollbackError) Error() string {
	if e.RollbackErr == nil {
		return e.Err.Error()
	}
	const format = "%s (additionally, rollback failed: %s)"
	return fmt.Sprintf(format, e.Err, e.RollbackErr)
}

// WithRollback attaches rollback error to err. If rollback is nil, then err is
// returned.
func WithRollback(err, rollback error) error {
	if rollback == nil {
		return errors.WithStack(err)
	}
	return RollbackError{
		Err:         err,
		RollbackErr: rollback,
	}
}
