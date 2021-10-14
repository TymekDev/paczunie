package main

import "sync"

// Storage is used by Client for storing and providing Pkg objects.
type Storage interface {
	StorePkg(Pkg) error
	LoadPkgs() ([]Pkg, error)
}

type sliceStorage struct {
	sync.Mutex

	pkgs []Pkg
}

var _ Storage = (*sliceStorage)(nil)

func (s *sliceStorage) StorePkg(p Pkg) error {
	s.Lock()
	defer s.Unlock()
	s.pkgs = append(s.pkgs, p)
	return nil
}

func (s *sliceStorage) LoadPkgs() ([]Pkg, error) {
	s.Lock()
	defer s.Unlock()
	return s.pkgs, nil
}
