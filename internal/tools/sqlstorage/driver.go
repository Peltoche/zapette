package sqlstorage

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"log/slog"

	"github.com/Peltoche/zapette/internal/tools"
	"github.com/mattn/go-sqlite3"
)

func newDBWithHooks(dsn string, hookList *SQLChangeHookList, tools tools.Tools) *sql.DB {
	driver := sqlite3.SQLiteDriver{
		ConnectHook: func(conn *sqlite3.SQLiteConn) error {
			if hookList == nil {
				return nil
			}

			conn.RegisterUpdateHook(func(op int, dbName string, tableName string, rowID int64) {
				for _, hook := range hookList.GetHooks() {
					err := hook.RunHook(tableName)
					if err != nil {
						tools.Logger().Error("RunHook error", slog.String("hook", hook.Name()), slog.String("error", err.Error()))
						return
					}

					tools.Logger().Debug("RunHook success", slog.String("hook", hook.Name()), slog.String("table", tableName))
				}
			})

			return nil
		},
	}

	return sql.OpenDB(dsnConnector{dsn: dsn, driver: &driver})
}

type dsnConnector struct {
	dsn    string
	driver driver.Driver
}

func (t dsnConnector) Connect(_ context.Context) (driver.Conn, error) {
	return t.driver.Open(t.dsn)
}

func (t dsnConnector) Driver() driver.Driver {
	return t.driver
}
