package sysstats

import (
	"context"
	"testing"
	"time"

	"github.com/Peltoche/zapette/internal/tools/clock"
	"github.com/Peltoche/zapette/internal/tools/sqlstorage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSystatSqlStorage(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	db := sqlstorage.NewTestStorage(t)
	store := newSqlStorage(db)
	clock := clock.NewDefault()

	time1 := clock.Now()
	time2 := clock.Now().Add(time.Hour)

	stats := NewFakeStats(t).
		WithTime(time1).
		Build()

	stats2 := NewFakeStats(t).
		WithTime(time2).
		Build()

	t.Run("Save succes", func(t *testing.T) {
		err := store.Save(ctx, MinGraph, stats)
		require.NoError(t, err)
	})

	t.Run("Save succes 2", func(t *testing.T) {
		err := store.Save(ctx, MinGraph, stats2)
		require.NoError(t, err)
	})

	t.Run("GetLatest succes", func(t *testing.T) {
		res, err := store.GetLatest(ctx)
		require.NoError(t, err)

		assert.EqualValues(t, stats2, res)
	})

	t.Run("GetRange succes", func(t *testing.T) {
		res, err := store.GetRange(ctx, MinGraph, time1, time2)
		require.NoError(t, err)

		// Do not include "stats" as the first value must be stricly superior and not equal
		assert.EqualValues(t, []Stats{*stats2}, res)
	})

	t.Run("GetRange with an empty namespace", func(t *testing.T) {
		res, err := store.GetRange(ctx, Unknown, time1, time2)
		require.NoError(t, err)

		assert.Empty(t, res)
	})

	t.Run("GetRange succes 2", func(t *testing.T) {
		res, err := store.GetRange(ctx, MinGraph, time1.Add(-time.Second), time2)
		require.NoError(t, err)

		assert.EqualValues(t, []Stats{*stats, *stats2}, res)
	})
}
