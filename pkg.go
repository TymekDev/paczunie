package main

// Pkg is a structure representing a single package entry.
type Pkg struct {
	// Name is a package name given by the user.
	Name string
}

// NewPkg creates is a Pkg struct constructor.
func NewPkg(name string) Pkg {
	p := Pkg{Name: name}
	return p
}
