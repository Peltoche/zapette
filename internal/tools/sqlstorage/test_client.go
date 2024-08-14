package sqlstorage

import (
	"database/sql"
	"testing"

	"github.com/Peltoche/zapette/internal/migrations"
	"github.com/Peltoche/zapette/internal/tools"
	"github.com/stretchr/testify/require"
)

func NewTestStorage(t *testing.T) *sql.DB {
	cfg := Config{Path: ":memory:"}
	tools := tools.NewToolboxForTest(t)

	db, err := NewSQliteClient(&cfg, nil, tools)
	require.NoError(t, err)

	err = db.Ping()
	require.NoError(t, err)

	err = migrations.Run(db, nil)
	require.NoError(t, err)

	return db
}
