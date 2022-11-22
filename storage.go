package main

import (
	"database/sql"
	"errors"
	"log"
	"sync"

	"github.com/google/uuid"
)

// Storage is used by Client for storing and providing Pkg objects.
type Storage interface {
	StorePkg(Pkg) error
	LoadPkgs() ([]Pkg, error)
	UpdatePkgStatus(uuid.UUID, Status) error
	DeletePkg(uuid.UUID) error
}

// DBStorage is a wrapper on provided *sql.DB fulfilling Storage interface.
type DBStorage struct {
	sync.Mutex

	db *sql.DB
}

var _ Storage = (*DBStorage)(nil)

// NewDBStorage creates a DBStorage that fulfills Storage interface. It is
// checked whether provided database contains Packages(ID, Name, PickupPoint, Status)
// table.
func NewDBStorage(db *sql.DB, initIfEmpty bool) (*DBStorage, error) {
	const query = "SELECT ID, Name, PickupPoint, Status FROM Packages"
	if _, err := db.Query(query); err != nil {
		if initIfEmpty {
			_, err := db.Exec(`CREATE TABLE Packages(ID TEXT NOT NULL, Name TEXT, PickupPoint INT, Status INT)`)
			if err != nil {
				return nil, err
			}
			log.Println("INFO", "initialized table: Packages(ID, Name, PickupPoint, Status)")
		} else {
			const msg = "table not found: Packages(ID, Name, PickupPoint, Status)"
			return nil, errors.New(msg)
		}
	}
	dbs := &DBStorage{
		db: db,
	}
	return dbs, nil
}

// StorePkg saves p into a Packages table via DBStorage's underlying database
// connection.
func (dbs *DBStorage) StorePkg(p Pkg) error {
	dbs.Lock()
	defer dbs.Unlock()

	const query = "INSERT INTO Packages(ID, Name, PickupPoint, Status) VALUES (?, ?, ?, ?)"
	if _, err := dbs.db.Exec(query, p.ID, p.Name, p.PickupPoint, p.Status); err != nil {
		return err
	}
	return nil
}

// LoadPkgs returns a []Pkg slice read from DBStorage's underlying database
// connection.
func (dbs *DBStorage) LoadPkgs() ([]Pkg, error) {
	const query = "SELECT ID, Name, PickupPoint, Status FROM Packages"
	rows, err := dbs.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var (
		pkgs        []Pkg
		id          uuid.UUID
		name        string
		pickupPoint bool
		status      Status
	)
	for rows.Next() {
		if err := rows.Scan(&id, &name, &pickupPoint, &status); err != nil {
			return nil, err
		}
		p := NewPkg(name, withUUID(id), WithPickupPoint(pickupPoint), WithStatus(status))
		pkgs = append([]Pkg{p}, pkgs...)
	}

	return pkgs, nil
}

// UpdatePkgStatus changes status of a package with provided ID to status.
func (dbs *DBStorage) UpdatePkgStatus(id uuid.UUID, status Status) error {
	dbs.Lock()
	defer dbs.Unlock()

	const query = "UPDATE Packages SET Status = ? WHERE ID = ?"
	if _, err := dbs.db.Exec(query, status, id); err != nil {
		return err
	}
	return nil
}

// DeletePkg deletes a package entry with provided ID.
func (dbs *DBStorage) DeletePkg(id uuid.UUID) error {
	dbs.Lock()
	defer dbs.Unlock()

	const query = "DELETE FROM Packages WHERE ID = ?"
	if _, err := dbs.db.Exec(query, id); err != nil {
		return err
	}
	return nil
}
