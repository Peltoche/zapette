package sqlstorage

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"net/url"

	"github.com/mattn/go-sqlite3"
)

const driverName = "sqlite3-with-hooks"

type Config struct {
	Path string `json:"path"`
}

// RowScanner is the interface that wraps the Scan method.
//
// Scan behaves like database/sql.Row.Scan.
type RowScanner interface {
	Scan(...interface{}) error
}

type Querier interface {
	Exec(query string, args ...any) (sql.Result, error)
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	Query(query string, args ...any) (*sql.Rows, error)
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	QueryRow(query string, args ...any) *sql.Row
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
}

type SQLChangeHook interface {
	Name() string
	ShouldRunHook(table string) bool
	RunHook(ctx context.Context, db Querier, table string) error
}

type Client struct {
	db *sql.DB
}

func NewSQLQuerier(db *sql.DB) Querier {
	return &Client{db}
}

func NewSQliteClient(cfg *Config, hooks []SQLChangeHook, logger *slog.Logger) (*sql.DB, error) {
	var db *sql.DB
	var err error

	connectionUrlParams := make(url.Values)
	connectionUrlParams.Add("_txlock", "immediate")
	connectionUrlParams.Add("_journal_mode", "WAL")
	connectionUrlParams.Add("_busy_timeout", "5000")
	connectionUrlParams.Add("_synchronous", "NORMAL")
	connectionUrlParams.Add("_cache_size", "1000000000")
	connectionUrlParams.Add("_foreign_keys", "true")

	dsn := "file:" + cfg.Path + "?" + connectionUrlParams.Encode()

	// The driver need a custom name due to custom registration with the hook.
	driverName := fmt.Sprintf("sqlite3-hooks-%d", len(sql.Drivers()))

	sql.Register(driverName, &sqlite3.SQLiteDriver{
		ConnectHook: func(conn *sqlite3.SQLiteConn) error {
			conn.RegisterUpdateHook(func(op int, dbName string, tableName string, rowID int64) {
				for _, hook := range hooks {
					if !hook.ShouldRunHook(tableName) {
						continue
					}

					err := hook.RunHook(nil, db, tableName)
					if err != nil {
						logger.Error("RunHook error", slog.String("hook", hook.Name()), slog.String("error", err.Error()))
						return
					}

					logger.Debug("RunHook success", slog.String("hook", hook.Name()), slog.String("table", tableName))
				}
			})

			return nil
		},
	})

	for _, hook := range hooks {
		logger.Debug("register SQL hook", slog.String("hook", hook.Name()))
	}

	logger.Info("start listening sql hooks", slog.Int("hooks", len(hooks)))

	db, err = sql.Open(driverName, dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open %q: %w", dsn, err)
	}

	db.SetMaxOpenConns(1)

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}

func (c *Client) Exec(query string, args ...any) (sql.Result, error) {
	return c.db.Exec(query, args...)
}

func (c *Client) ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error) {
	return c.db.ExecContext(ctx, query, args...)
}

func (c *Client) Query(query string, args ...any) (*sql.Rows, error) {
	return c.db.Query(query, args...)
}

func (c *Client) QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error) {
	return c.db.QueryContext(ctx, query, args...)
}

func (c *Client) QueryRow(query string, args ...any) *sql.Row {
	return c.db.QueryRow(query, args...)
}

func (c *Client) QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row {
	return c.db.QueryRowContext(ctx, query, args...)
}
