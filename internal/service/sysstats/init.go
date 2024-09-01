package sysstats

import (
	"context"
	"database/sql"

	"github.com/Peltoche/zapette/internal/tools"
	"github.com/Peltoche/zapette/internal/tools/sqlstorage"
	"github.com/spf13/afero"
	"go.uber.org/fx"
)

type Result struct {
	fx.Out
	Service Service
	Watcher sqlstorage.SQLChangeHook `group:"hooks"`
	Cron    *SystatsCron
}

type Service interface {
	GetLatest(ctx context.Context) (*Stats, error)
	GetStatsForGraph(ctx context.Context, graph *Graph) ([]Stats, error)
	Watch(ctx context.Context) chan struct{}
	fetchAndRegister(ctx context.Context) (*Stats, error)
}

func Init(db *sql.DB, fs afero.Fs, tools tools.Tools) Result {
	storage := newSqlStorage(db)

	svc := newService(storage, fs, tools)

	return Result{
		Service: svc,
		Watcher: svc,
		Cron:    newSystatCron(svc, tools),
	}
}
