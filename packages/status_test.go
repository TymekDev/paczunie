package packages

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStatusFromInt_Convert(t *testing.T) {
	tests := []int{1, 2, 3, 4}
	for _, tt := range tests {
		t.Run(strconv.Itoa(tt), func(t *testing.T) {
			status, err := StatusFromInt(tt)
			assert.NoError(t, err)
			assert.Equal(t, Status(tt), status)
		})
	}
}

func TestStatusFromInt_Error(t *testing.T) {
	tests := []int{-999, -100, -1, 0, 5, 100, 999}

	for _, tt := range tests {
		t.Run(strconv.Itoa(tt), func(t *testing.T) {
			status, err := StatusFromInt(tt)
			assert.Equal(t, Status(0), status)
			assert.Error(t, err)
		})
	}
}

func TestStatusFromString_Convert(t *testing.T) {
	tests := []int{1, 2, 3, 4}
	for _, tt := range tests {
		t.Run(strconv.Itoa(tt), func(t *testing.T) {
			status, err := StatusFromString(strconv.Itoa(tt))
			assert.NoError(t, err)
			assert.Equal(t, Status(tt), status)
		})
	}
}

func TestStatusFromString_Error(t *testing.T) {
	tests := []string{"", "0", "asdf", "xyz", "1.0"}

	for _, tt := range tests {
		t.Run(tt, func(t *testing.T) {
			status, err := StatusFromString(tt)
			assert.Equal(t, Status(0), status)
			assert.Error(t, err)
		})
	}
}
