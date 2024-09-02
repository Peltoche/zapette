package timeseries

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/Peltoche/zapette/internal/tools/uuid"
)

const (
	dataTable   = "timeseries_data"
	configTable = "timeseries_config"
)

var errNotFound = errors.New("not found")

var (
	dataFields   = []string{"ts_id", "time_unix_sec", "content"}
	configFields = []string{"ts_id", "graph_span_s", "tick_span_s"}
)

type sqlStorage struct {
	db *sql.DB
}

func newSqlStorage(db *sql.DB) *sqlStorage {
	return &sqlStorage{db}
}

func (s *sqlStorage) saveData(ctx context.Context, tsID uuid.UUID, t time.Time, data []byte) error {
	_, err := sq.
		Insert(dataTable).
		Columns(dataFields...).
		Values(
			string(tsID),
			t.Unix(),
			data,
		).
		RunWith(s.db).
		ExecContext(ctx)
	if err != nil {
		return fmt.Errorf("sql error: %w", err)
	}

	return nil
}

func (s *sqlStorage) getLatestData(ctx context.Context, tsID uuid.UUID) (*TimeData, error) {
	var res TimeData
	var unixTime int64

	err := sq.
		Select(dataFields...).
		Where(sq.Eq{"ts_id": tsID}).
		OrderBy("time DESC").
		Limit(1).
		From(dataTable).
		RunWith(s.db).
		ScanContext(ctx,
			&res.tsID,
			&unixTime,
			&res.data,
		)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, errNotFound
	}

	if err != nil {
		return nil, fmt.Errorf("sql error: %w", err)
	}

	res.at = time.UnixMilli(unixTime)

	return &res, nil
}

func (s *sqlStorage) saveConfig(ctx context.Context, ts *Timeserie) error {
	_, err := sq.
		Insert(dataTable).
		Columns(configFields...).
		Values(
			string(ts.tsID),
			int(ts.graphSpan.Seconds()),
			int(ts.tickSpan.Seconds()),
		).
		RunWith(s.db).
		ExecContext(ctx)
	if err != nil {
		return fmt.Errorf("sql error: %w", err)
	}

	return nil
}

func (s *sqlStorage) getConfigByID(ctx context.Context, tsID uuid.UUID) (*Timeserie, error) {
	res := Timeserie{}

	err := sq.
		Select(configFields...).
		From(configTable).
		Where(sq.Eq{"ts_id": string(tsID)}).
		RunWith(s.db).
		ScanContext(ctx, &res.tsID, &res.graphSpan, &res.tickSpan)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, errNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to query the db: %w", err)
	}

	return &res, nil
}

func (s *sqlStorage) getAllConfigs(ctx context.Context) ([]Timeserie, error) {
	rows, err := sq.
		Select(configFields...).
		From(configTable).
		RunWith(s.db).
		QueryContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to query the db: %w", err)
	}

	return s.scanConfigRows(rows)
}

func (s *sqlStorage) scanConfigRows(rows *sql.Rows) ([]Timeserie, error) {
	timeseries := []Timeserie{}

	for rows.Next() {
		ts := Timeserie{}

		err := rows.Scan(
			&ts.tsID,
			&ts.graphSpan,
			&ts.tickSpan,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan the result: %w", err)
		}

		timeseries = append(timeseries, ts)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("scan error: %w", err)
	}

	return timeseries, nil
}
