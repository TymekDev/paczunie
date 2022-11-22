package main

import (
	"database/sql"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	_ "modernc.org/sqlite"
)

func TestDBStorage_NewDBStorage(t *testing.T) {
	db, err := sql.Open("sqlite", ":memory:")
	require.NoError(t, err)

	_, err = NewDBStorage(db, false)
	assert.EqualError(t, err, "table not found: Packages(ID, Name, Inpost, Status)")

	const query = "CREATE TABLE Packages(ID, Name, Inpost, Status)"
	_, err = db.Exec(query)
	require.NoError(t, err)

	dbs, err := NewDBStorage(db, false)
	require.NoError(t, err)
	assert.Equal(t, db, dbs.db)
}

func TestDBStorage_StorePkg(t *testing.T) {
	// TODO: add tests
}

func TestDBStorage_LoadPkgs(t *testing.T) {
	// TODO: add tests
}
