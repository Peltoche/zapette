package sqlstorage

import (
	"testing"

	"github.com/Peltoche/zapette/internal/migrations"
	"github.com/Peltoche/zapette/internal/tools"
	"github.com/stretchr/testify/require"
)

func NewTestStorage(t *testing.T) Querier {
	cfg := Config{Path: ":memory:"}
	hooks := []SQLChangeHook{}
	tools := tools.NewToolboxForTest(t)

	db, err := NewSQliteClient(&cfg, hooks, tools.Logger())
	require.NoError(t, err)

	err = db.Ping()
	require.NoError(t, err)

	err = migrations.Run(db, nil)
	require.NoError(t, err)

	querier := NewSQLQuerier(db)

	return querier
}
