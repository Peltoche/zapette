package config

import (
	"context"
	"testing"

	"github.com/Peltoche/zapette/internal/tools/sqlstorage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSQLStorage(t *testing.T) {
	ctx := context.Background()

	db := sqlstorage.NewTestStorage(t)
	store := newSqlStorage(db)

	t.Run("Save success", func(t *testing.T) {
		err := store.Save(ctx, sysstatsInputNamespace, "some-content")
		require.NoError(t, err)
	})

	t.Run("Get success", func(t *testing.T) {
		res, err := store.Get(ctx, sysstatsInputNamespace)
		require.NoError(t, err)
		assert.Equal(t, "some-content", res)
	})
}
