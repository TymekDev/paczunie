package packages

import (
	"database/sql"
	"testing"

	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDBStorage_NewDBStorage(t *testing.T) {
	db, err := sql.Open("sqlite3", ":memory:")
	require.NoError(t, err)

	_, err = NewDBStorage(db)
	assert.EqualError(t, err, "table not found: Packages(ID, Name, Inpost, Status)")

	const query = "CREATE TABLE Packages(ID, Name, Inpost, Status)"
	_, err = db.Exec(query)
	require.NoError(t, err)

	dbs, err := NewDBStorage(db)
	require.NoError(t, err)
	assert.Equal(t, db, dbs.db)
}
