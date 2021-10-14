package main

import "sync"

// Storage is used by Client for storing and providing Pkg objects.
type Storage interface {
	StorePkg(Pkg)
	LoadPkgs() []Pkg
}

type sliceStorage struct {
	sync.Mutex

	pkgs []Pkg
}

var _ Storage = (*sliceStorage)(nil)

func (s *sliceStorage) StorePkg(p Pkg) {
	s.Lock()
	defer s.Unlock()
	s.pkgs = append(s.pkgs, p)
}

func (s *sliceStorage) LoadPkgs() []Pkg {
	s.Lock()
	defer s.Unlock()
	return s.pkgs
}
