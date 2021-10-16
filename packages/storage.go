package packages

import (
	"database/sql"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

// Storage is used by Client for storing and providing Pkg objects.
type Storage interface {
	StorePkg(Pkg) error
	LoadPkgs() ([]Pkg, error)
}

// DBStorage is a wrapper on provided *sql.DB fulfilling Storage interface.
type DBStorage struct {
	db *sql.DB
}

var _ Storage = (*DBStorage)(nil)

// NewDBStorage creates a DBStorage that fulfills Storage interface.
func NewDBStorage(db *sql.DB) *DBStorage {
	return &DBStorage{
		db: db,
	}
}

// StorePkg saves p into a Packages table via DBStorage's underlying database
// connection.
func (dbs *DBStorage) StorePkg(p Pkg) error {
	tx, err := dbs.db.Begin()
	if err != nil {
		return errors.WithStack(err)
	}

	const query = "INSERT INTO Packages(ID, Name, Inpost, Status) VALUES (?, ?, ?, ?)"
	stmt, err := tx.Prepare(query)
	if err != nil {
		return errors.WithStack(withRollback(err, tx.Rollback()))
	}

	if _, err := stmt.Exec(p.ID, p.Name, p.Inpost, p.Status); err != nil {
		return errors.WithStack(withRollback(err, tx.Rollback()))
	}

	if err := tx.Commit(); err != nil {
		return errors.WithStack(withRollback(err, tx.Rollback()))
	}

	return nil
}

// LoadPkgs returns a []Pkg slice read from DBStorage's underlying database
// connection.
func (dbs *DBStorage) LoadPkgs() ([]Pkg, error) {
	const query = "SELECT ID, Name, Inpost, Status FROM Packages"
	rows, err := dbs.db.Query(query)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	defer rows.Close()

	var (
		pkgs   []Pkg
		id     uuid.UUID
		name   string
		inpost bool
		status Status
	)
	for rows.Next() {
		if err := rows.Scan(&id, &name, &inpost, &status); err != nil {
			return nil, errors.WithStack(err)
		}
		p := NewPkg(name, withUUID(id), WithInpost(inpost), WithStatus(status))
		pkgs = append(pkgs, p)
	}

	return pkgs, nil
}
