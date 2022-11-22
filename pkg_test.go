package main

import (
	"strconv"
	"testing"

	"github.com/google/uuid"
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
			nil,
			Pkg{Name: "none", Status: Ordered},
		},
		{
			"uuid",
			[]PkgOption{
				withUUID(uuid.MustParse("00000000-0000-0000-0000-000000000000")),
			},
			Pkg{
				ID:     uuid.MustParse("00000000-0000-0000-0000-000000000000"),
				Name:   "uuid",
				Status: Ordered,
			},
		},
		{
			"inpost",
			[]PkgOption{WithInpost(true)},
			Pkg{Name: "inpost", Inpost: true, Status: Ordered},
		},
		{
			"status_1",
			[]PkgOption{WithStatus(Delivered)},
			Pkg{Name: "status_1", Status: Delivered},
		},
		{
			"status_2",
			[]PkgOption{WithStatus(Shipped)},
			Pkg{Name: "status_2", Status: Shipped},
		},
		{
			"inpost_status",
			[]PkgOption{WithInpost(true), WithStatus(Ordered)},
			Pkg{Name: "inpost_status", Inpost: true, Status: Ordered},
		},
		{
			"uuid_inpost_status_uuid",
			[]PkgOption{
				withUUID(uuid.MustParse("00000000-0000-0000-0000-000000000000")),
				WithInpost(true),
				WithStatus(Ordered),
			},
			Pkg{
				ID:     uuid.MustParse("00000000-0000-0000-0000-000000000000"),
				Name:   "uuid_inpost_status_uuid",
				Inpost: true,
				Status: Ordered,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewPkg(tt.name, tt.options...)
			assert.Equal(t, tt.want.Name, p.Name)
			assert.Equal(t, tt.want.Inpost, p.Inpost)
			assert.Equal(t, tt.want.Status, p.Status)
		})
	}
}

// TODO: make UUID tests deterministic
func TestUUIDOpt_apply(t *testing.T) {
	tests := []uuid.UUID{
		uuid.New(),
		uuid.MustParse("00000000-0000-0000-0000-000000000000"),
	}

	for _, tt := range tests {
		t.Run(tt.String(), func(t *testing.T) {
			var p Pkg
			uuidOpt(tt).apply(&p)
			assert.Equal(t, tt, p.ID)
		})
	}
}

func TestUUIDOpt_withUUID(t *testing.T) {
	tests := []uuid.UUID{
		uuid.New(),
		uuid.MustParse("00000000-0000-0000-0000-000000000000"),
	}

	for _, tt := range tests {
		t.Run(tt.String(), func(t *testing.T) {
			var p Pkg
			withUUID(tt).apply(&p)
			assert.Equal(t, tt, p.ID)
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
			var p Pkg
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
			var p Pkg
			WithInpost(tt).apply(&p)
			assert.Equal(t, tt, p.Inpost)
		})
	}
}

func TestStatusOpt_apply(t *testing.T) {
	tests := []Status{
		Ordered,
		Shipped,
		Delivered,
	}

	for _, tt := range tests {
		t.Run(strconv.Itoa(int(tt)), func(t *testing.T) {
			var p Pkg
			statusOpt(tt).apply(&p)
			assert.Equal(t, tt, p.Status)
		})
	}
}

func TestStatusOpt_WithStatus(t *testing.T) {
	tests := []Status{
		Ordered,
		Shipped,
		Delivered,
	}

	for _, tt := range tests {
		t.Run(strconv.Itoa(int(tt)), func(t *testing.T) {
			var p Pkg
			WithStatus(tt).apply(&p)
			assert.Equal(t, tt, p.Status)
		})
	}
}
