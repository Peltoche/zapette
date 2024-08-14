package sqlstorage

import (
	"database/sql"
	"testing"

	"github.com/Peltoche/zapette/internal/migrations"
	"github.com/stretchr/testify/require"
)

func NewTestStorage(t *testing.T) *sql.DB {
	cfg := Config{Path: ":memory:"}

	db, err := NewSQliteClient(&cfg)
	require.NoError(t, err)

	err = db.Ping()
	require.NoError(t, err)

	err = migrations.Run(db, nil)
	require.NoError(t, err)

	return db
}
