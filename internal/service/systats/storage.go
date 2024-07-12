package systats

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/Peltoche/zapette/internal/tools/ptr"
	"github.com/Peltoche/zapette/internal/tools/sqlstorage"
)

const tableName = "systats"

var errNotFound = errors.New("not found")

var allFields = []string{
	"time",
	"total_mem",
	"available_mem",
	"free_mem",
	"buffers",
	"cached",
	"s_reclaimable",
	"sh_mem",
	"total_swap",
	"free_swap",
}

// sqlStorage use to save/retrieve Users
type sqlStorage struct {
	db sqlstorage.Querier
}

// newSqlStorage instantiates a new Storage based on sql.
func newSqlStorage(db sqlstorage.Querier) *sqlStorage {
	return &sqlStorage{db}
}

// Save the given User.
func (s *sqlStorage) Save(ctx context.Context, stats *Stats) error {
	_, err := sq.
		Insert(tableName).
		Columns(allFields...).
		Values(
			ptr.To(sqlstorage.SQLTime(stats.time)),
			stats.memory.totalMem,
			stats.memory.availableMem,
			stats.memory.freeMem,
			stats.memory.buffers,
			stats.memory.cached,
			stats.memory.sReclaimable,
			stats.memory.shmem,
			stats.memory.totalSwap,
			stats.memory.freeSwap,
		).
		RunWith(s.db).
		ExecContext(ctx)
	if err != nil {
		return fmt.Errorf("sql error: %w", err)
	}

	return nil
}

func (s *sqlStorage) GetLatest(ctx context.Context) (*Stats, error) {
	res := Stats{
		memory: &Memory{},
	}

	query := sq.
		Select(allFields...).
		OrderBy("time DESC").
		Limit(1).
		From(tableName)

	var sqlTime sqlstorage.SQLTime
	err := query.
		RunWith(s.db).
		ScanContext(ctx,
			&sqlTime,
			&res.memory.totalMem,
			&res.memory.availableMem,
			&res.memory.freeMem,
			&res.memory.buffers,
			&res.memory.cached,
			&res.memory.sReclaimable,
			&res.memory.shmem,
			&res.memory.totalSwap,
			&res.memory.freeSwap,
		)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, errNotFound
	}

	res.time = sqlTime.Time()

	if err != nil {
		return nil, fmt.Errorf("sql error: %w", err)
	}

	return &res, nil
}
