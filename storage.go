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

type dbStorage struct {
	conn *sql.DB
}

var _ Storage = (*dbStorage)(nil)

func newDBStorage(conn *sql.DB) *dbStorage {
	return &dbStorage{
		conn: conn,
	}
}

func (db *dbStorage) StorePkg(p Pkg) error {
	tx, err := db.conn.Begin()
	if err != nil {
		return errors.WithStack(err)
	}

	const query = "INSERT INTO Packages(Name, Inpost, Status) VALUES (?, ?, ?)"
	stmt, err := tx.Prepare(query)
	if err != nil {
		return WithRollback(err, tx.Rollback())
	}

	if _, err := stmt.Exec(p.Name, p.Inpost, p.Status); err != nil {
		return WithRollback(err, tx.Rollback())
	}

	if err := tx.Commit(); err != nil {
		return WithRollback(err, tx.Rollback())
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
