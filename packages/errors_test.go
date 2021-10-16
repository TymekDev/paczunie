package packages

import (
	"strconv"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func TestRollbackError_Error(t *testing.T) {
	tests := []struct {
		err         error
		rollbackErr error
		want        string
	}{
		{
			nil,
			nil,
			"",
		},
		{
			errors.New("some error"),
			nil,
			"some error",
		},
		{
			errors.New("some error"),
			errors.New("rollback failure"),
			"some error (additionally, rollback failed: rollback failure)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			err := rollbackError{err: tt.err, rollbackErr: tt.rollbackErr}
			assert.Equal(t, tt.want, err.Error())
		})
	}
}

func TestRollbackError_WithRollback(t *testing.T) {
	tests := []struct {
		err         error
		rollbackErr error
	}{
		{
			nil,
			nil,
		},
		{
			errors.New("some error"),
			nil,
		},
		{
			errors.New("some error"),
			errors.New("rollback failure"),
		},
	}

	for i, tt := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			err := withRollback(tt.err, tt.rollbackErr)
			if tt.rollbackErr == nil {
				assert.Equal(t, tt.err, err)
			} else {
				assert.Equal(t, rollbackError{tt.err, tt.rollbackErr}, err)
			}
		})
	}
}
