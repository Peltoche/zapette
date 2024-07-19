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

var allFields = []string{"time", "content"}

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
	rawStats, _ := stats.MarshalBinary()

	_, err := sq.
		Insert(tableName).
		Columns(allFields...).
		Values(
			ptr.To(sqlstorage.SQLTime(stats.time)),
			rawStats,
		).
		RunWith(s.db).
		ExecContext(ctx)
	if err != nil {
		return fmt.Errorf("sql error: %w", err)
	}

	return nil
}

func (s *sqlStorage) GetLatest(ctx context.Context) (*Stats, error) {
	rawContent := []byte{}

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
			&rawContent,
		)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, errNotFound
	}

	var res Stats

	err = res.UnmarshalBinary(rawContent)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal the stats: %w", err)
	}

	return &res, nil
}
