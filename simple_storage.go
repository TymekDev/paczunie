package main

import "sync"

type simpleStorage struct {
	sync.Mutex

	pkgs []Pkg
}

var _ Storage = (*simpleStorage)(nil)

func (s *simpleStorage) Store(p Pkg) {
	s.Lock()
	defer s.Unlock()
	s.pkgs = append(s.pkgs, p)
}

func (s *simpleStorage) Values() []Pkg {
	s.Lock()
	defer s.Unlock()
	return s.pkgs
}
