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
	UpdateStatus(uuid.UUID, Status) error
	DeletePkg(uuid.UUID) error
}

// DBStorage is a wrapper on provided *sql.DB fulfilling Storage interface.
type DBStorage struct {
	db *sql.DB
}

var _ Storage = (*DBStorage)(nil)

// NewDBStorage creates a DBStorage that fulfills Storage interface. It is
// checked whether provided database contains Packages(ID, Name, Inpost, Status)
// table.
func NewDBStorage(db *sql.DB) (*DBStorage, error) {
	const query = "SELECT ID, Name, Inpost, Status FROM Packages"
	if _, err := db.Query(query); err != nil {
		const msg = "table not found: Packages(ID, Name, Inpost, Status)"
		return nil, errors.New(msg)
	}
	dbs := &DBStorage{
		db: db,
	}
	return dbs, nil
}

// StorePkg saves p into a Packages table via DBStorage's underlying database
// connection.
func (dbs *DBStorage) StorePkg(p Pkg) error {
	const query = "INSERT INTO Packages(ID, Name, Inpost, Status) VALUES (?, ?, ?, ?)"
	if _, err := dbs.db.Exec(query, p.ID, p.Name, p.Inpost, p.Status); err != nil {
		return errors.WithStack(err)
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

// UpdateStatus changes status of a package with provided ID to status.
func (dbs *DBStorage) UpdateStatus(id uuid.UUID, status Status) error {
	const query = "UPDATE Packages SET Status = ? WHERE ID = ?"
	if _, err := dbs.db.Exec(query, status, id); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

// DeletePkg deletes a package entry with provided ID.
func (dbs *DBStorage) DeletePkg(id uuid.UUID) error {
	const query = "DELETE FROM Packages WHERE ID = ?"
	if _, err := dbs.db.Exec(query, id); err != nil {
		return errors.WithStack(err)
	}
	return nil
}
