package main

import "github.com/google/uuid"

// Pkg is a structure representing a single package entry.
type Pkg struct {
	// ID is an UUID assigned to package during its creation.
	ID uuid.UUID
	// Name is a package name given by the user.
	Name string
	// PickupPoint a flag whether package will arrive at PickupPoint parcel locker.
	PickupPoint bool
	// Status denotes current status of a Pkg.
	Status Status
}

// NewPkg creates is a Pkg struct constructor.
func NewPkg(name string, options ...PkgOption) Pkg {
	p := Pkg{ID: uuid.New(), Name: name, Status: Ordered}
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

type pickupPointOpt bool

var _ PkgOption = pickupPointOpt(false)

func (o pickupPointOpt) apply(p *Pkg) {
	p.PickupPoint = bool(o)
}

// WithPickupPoint returns a PkgOption setting PickupPoint field in Pkg struct.
func WithPickupPoint(x bool) PkgOption {
	return pickupPointOpt(x)
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
