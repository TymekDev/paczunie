package main

// Pkg is a structure representing a single package entry.
type Pkg struct {
	// Name is a package name given by the user.
	Name string
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
