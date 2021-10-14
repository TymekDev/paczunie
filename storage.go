package main

import (
	"database/sql"
	"sync"

	_ "github.com/mattn/go-sqlite3"
	"github.com/pkg/errors"
)

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

type dbStorage sql.DB

var _ Storage = (*dbStorage)(nil)

func (db *dbStorage) StorePkg(p Pkg) error {
	return errors.New("not implemented")
}

func (db *dbStorage) LoadPkgs() ([]Pkg, error) {
	return nil, errors.New("not implemented")
}
