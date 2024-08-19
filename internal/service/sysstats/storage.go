package sysstats

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
)

const tableName = "sysstats"

var errNotFound = errors.New("not found")

var allFields = []string{"time", "namespace", "content"}

// sqlStorage use to save/retrieve Users
type sqlStorage struct {
	db *sql.DB
}

// newSqlStorage instantiates a new Storage based on sql.
func newSqlStorage(db *sql.DB) *sqlStorage {
	return &sqlStorage{db}
}

// Save the given User.
func (s *sqlStorage) Save(ctx context.Context, ns Namespace, stats *Stats) error {
	rawStats, _ := stats.MarshalBinary()

	_, err := sq.
		Insert(tableName).
		Columns(allFields...).
		Values(
			stats.time.Unix(),
			ns,
			rawStats,
		).
		RunWith(s.db).
		ExecContext(ctx)
	if err != nil {
		return fmt.Errorf("sql error: %w", err)
	}

	return nil
}

func (s *sqlStorage) GetRange(ctx context.Context, ns Namespace, start time.Time, end time.Time) ([]Stats, error) {
	rows, err := sq.
		Select(allFields...).
		Where(sq.And{sq.Eq{"namespace": ns}, sq.GtOrEq{"time": start.Unix()}, sq.LtOrEq{"time": end.Unix()}}).
		OrderBy("time ASC").
		From(tableName).
		RunWith(s.db).
		QueryContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to query the db: %w", err)
	}

	return s.scanRows(rows)
}

func (s *sqlStorage) GetLatest(ctx context.Context) (*Stats, error) {
	rawContent := []byte{}
	var unixTime int64
	var ns Namespace

	err := sq.
		Select(allFields...).
		OrderBy("time DESC").
		Limit(1).
		From(tableName).
		RunWith(s.db).
		ScanContext(ctx,
			&unixTime,
			&ns,
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

func (s *sqlStorage) scanRows(rows *sql.Rows) ([]Stats, error) {
	stats := []Stats{}
	var ns Namespace

	for rows.Next() {
		var unixTime int64
		rawContent := []byte{}

		err := rows.Scan(
			&unixTime,
			&ns,
			&rawContent,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan the result: %w", err)
		}

		var res Stats
		err = res.UnmarshalBinary(rawContent)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal the stats: %w", err)
		}

		stats = append(stats, res)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("scan error: %w", err)
	}

	return stats, nil
}
