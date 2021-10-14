package main

import (
	"database/sql"
	"sync"

	_ "github.com/mattn/go-sqlite3"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
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

type dbStorage struct {
	conn *sql.DB
}

var _ Storage = (*dbStorage)(nil)

func (db *dbStorage) StorePkg(p Pkg) error {
	tx, err := db.conn.Begin()
	if err != nil {
		return errors.WithStack(err)
	}

	const query = "INSERT INTO Packages(Name, Inpost, Status) VALUES (?, ?, ?)"
	stmt, err := tx.Prepare(query)
	if err != nil {
		_ = tx.Rollback() // TODO: handle error
		return errors.WithStack(err)
	}

	r, err := stmt.Exec(p.Name, p.Inpost, p.Status)
	if err != nil {
		_ = tx.Rollback() // TODO: handle error
		return errors.WithStack(err)
	}
	log.Debug().Interface("result", r).Send()

	if err := tx.Commit(); err != nil {
		_ = tx.Rollback() // TODO: handle error
		return errors.WithStack(err)
	}

	return nil
}

func (db *dbStorage) LoadPkgs() ([]Pkg, error) {
	const query = "SELECT Name, Inpost, Status FROM Packages"
	rows, err := db.conn.Query(query)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	defer rows.Close()

	var (
		pkgs   []Pkg
		name   string
		inpost bool
		status Status
	)
	for rows.Next() {
		if err := rows.Scan(&name, &inpost, &status); err != nil {
			return nil, errors.WithStack(err)
		}
		p := NewPkg(name, WithInpost(inpost), WithStatus(status))
		pkgs = append(pkgs, p)
	}

	return pkgs, nil
}
