package sqlstorage

import (
	"database/sql"
	"net/url"

	"github.com/Peltoche/zapette/internal/tools"
)

type Config struct {
	Path string `json:"path"`
}

// RowScanner is the interface that wraps the Scan method.
//
// Scan behaves like database/sql.Row.Scan.
type RowScanner interface {
	Scan(...interface{}) error
}

func NewSQliteClient(cfg *Config, hookList *SQLChangeHookList, tools tools.Tools) (*sql.DB, error) {
	connectionUrlParams := make(url.Values)
	connectionUrlParams.Add("_txlock", "immediate")
	connectionUrlParams.Add("_journal_mode", "WAL")
	connectionUrlParams.Add("_busy_timeout", "5000")
	connectionUrlParams.Add("_synchronous", "NORMAL")
	connectionUrlParams.Add("_cache_size", "1000000000")
	connectionUrlParams.Add("_foreign_keys", "true")

	dsn := "file:" + cfg.Path + "?" + connectionUrlParams.Encode()

	db := newDBWithHooks(dsn, hookList, tools)

	db.SetMaxOpenConns(1)

	err := db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}
