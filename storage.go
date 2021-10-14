package main

import "sync"

type Storage interface {
	Store(Pkg)
	Values() []Pkg
}

type sliceStorage struct {
	sync.Mutex

	pkgs []Pkg
}

var _ Storage = (*sliceStorage)(nil)

func (s *sliceStorage) Store(p Pkg) {
	s.Lock()
	defer s.Unlock()
	s.pkgs = append(s.pkgs, p)
}

func (s *sliceStorage) Values() []Pkg {
	s.Lock()
	defer s.Unlock()
	return s.pkgs
}
