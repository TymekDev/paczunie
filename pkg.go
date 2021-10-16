package main

import "github.com/google/uuid"

// Status denotes current status of a Pkg.
type Status int

const (
	// Ordered Status means that an order has been placed.
	Ordered Status = iota + 1
	// Shipped Status means that a package has been posted.
	Shipped
	// Delivered Status means that a package has been delivered.
	Delivered
)

// Pkg is a structure representing a single package entry.
type Pkg struct {
	// ID is an UUID assigned to package during its creation.
	ID uuid.UUID
	// Name is a package name given by the user.
	Name string
	// Inpost a flag whether package will arrive at Inpost parcel locker.
	Inpost bool
	// Status denotes current status of a Pkg.
	Status Status
}

// NewPkg creates is a Pkg struct constructor.
func NewPkg(name string, options ...PkgOption) Pkg {
	p := Pkg{ID: uuid.New(), Name: name}
	for _, o := range options {
		o.apply(&p)
	}
	return p
}

// PkgOption is an interface used to pass options to NewPkg constructor.
type PkgOption interface {
	apply(*Pkg)
}

type uuidOpt uuid.UUID

var _ PkgOption = uuidOpt{}

func (o uuidOpt) apply(p *Pkg) {
	p.ID = uuid.UUID(o)
}

// WithUUID returns a PkgOption setting ID field in Pkg struct.
func withUUID(x uuid.UUID) PkgOption {
	return uuidOpt(x)
}

type inpostOpt bool

var _ PkgOption = inpostOpt(false)

func (o inpostOpt) apply(p *Pkg) {
	p.Inpost = bool(o)
}

// WithInpost returns a PkgOption setting Inpost field in Pkg struct.
func WithInpost(x bool) PkgOption {
	return inpostOpt(x)
}

type statusOpt Status

var _ PkgOption = statusOpt(0)

func (o statusOpt) apply(p *Pkg) {
	p.Status = Status(o)
}

// WithStatus returns a PkgOption setting Status field in Pkg struct to a
// provided Status value.
func WithStatus(x Status) PkgOption {
	return statusOpt(x)
}
