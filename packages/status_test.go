package packages

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestToStatus_Convert(t *testing.T) {
	tests := []int{1, 2, 3}
	for _, tt := range tests {
		t.Run(strconv.Itoa(tt), func(t *testing.T) {
			status, err := ToStatus(tt)
			assert.NoError(t, err)
			assert.Equal(t, Status(tt), status)
		})
	}
}

func TestToStatus_Error(t *testing.T) {
	tests := []int{-999, -100, -1, 0, 4, 5, 100, 999}

	for _, tt := range tests {
		t.Run(strconv.Itoa(tt), func(t *testing.T) {
			status, err := ToStatus(tt)
			assert.Equal(t, Status(0), status)
			assert.Error(t, err)
		})
	}
}
