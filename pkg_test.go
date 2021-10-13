package main

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewPkg(t *testing.T) {
	tests := []struct {
		name    string
		options []PkgOption
		want    Pkg
	}{
		{
			"none",
			[]PkgOption{},
			Pkg{Name: "none"},
		},
		{
			"inpost",
			[]PkgOption{WithInpost(true)},
			Pkg{Name: "inpost", Inpost: true},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewPkg(tt.name, tt.options...)
			assert.Equal(t, tt.want, p)
		})
	}
}

func TestInpostOpt_apply(t *testing.T) {
	tests := []bool{
		true,
		false,
	}

	for _, tt := range tests {
		t.Run(strconv.FormatBool(tt), func(t *testing.T) {
			p := Pkg{}
			inpostOpt(tt).apply(&p)
			assert.Equal(t, tt, p.Inpost)
		})
	}
}

func TestInpostOpt_WithInpost(t *testing.T) {
	tests := []bool{
		true,
		false,
	}

	for _, tt := range tests {
		t.Run(strconv.FormatBool(tt), func(t *testing.T) {
			p := Pkg{}
			WithInpost(tt).apply(&p)
			assert.Equal(t, tt, p.Inpost)
		})
	}
}
