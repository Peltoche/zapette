package config

import (
	"context"
	"testing"

	"github.com/Peltoche/zapette/internal/tools"
	"github.com/Peltoche/zapette/internal/tools/sqlstorage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConfig(t *testing.T) {
	ctx := context.Background()

	tools := tools.NewToolboxForTest(t)
	db := sqlstorage.NewTestStorage(t)
	store := newSqlStorage(db)
	svc := newService(store, tools)

	someID := tools.UUID().New()

	t.Run("SetSysstatsInputNamespace success", func(t *testing.T) {
		err := svc.SetSysstatInputNamespace(ctx, someID)
		require.NoError(t, err)
	})

	t.Run("GetSysstatsInputNamespace success", func(t *testing.T) {
		res, err := svc.GetSysstatInputNamespace(ctx)
		require.NoError(t, err)

		assert.Equal(t, &someID, res)
	})
}
