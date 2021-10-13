package main

// Pkg is a structure representing a single package entry.
type Pkg struct {
	// Name is a package name given by the user.
	Name string

	// Inpost a flag whether package will arrive at Inpost parcel locker.
	Inpost bool
}

// NewPkg creates is a Pkg struct constructor.
func NewPkg(name string, options ...PkgOption) Pkg {
	p := Pkg{Name: name}
	for _, o := range options {
		o.apply(&p)
	}
	return p
}

// PkgOption is an interface used to pass options to NewPkg constructor.
type PkgOption interface {
	apply(*Pkg)
}

type inpostOpt bool

var _ PkgOption = inpostOpt(false)

func (o inpostOpt) apply(p *Pkg) {
	p.Inpost = bool(o)
}

// WithInpost returns a PkgOption setting Inpost field in Pkg struct to true.
func WithInpost(x bool) PkgOption {
	return inpostOpt(x)
}
