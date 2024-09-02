package timeseries

import (
	"context"
	"database/sql"
	"time"

	"github.com/Peltoche/zapette/internal/tools"
	"github.com/Peltoche/zapette/internal/tools/uuid"
)

type Service interface {
	CreateTimeSerie(ctx context.Context, cmd *CreateCmd) (*Timeserie, error)
	GetAllTimeseries(ctx context.Context) ([]Timeserie, error)
	GetTimeserieByID(ctx context.Context, id uuid.UUID) (*Timeserie, error)

	SaveLatestData(ctx context.Context, ts *Timeserie, at time.Time, data []byte) error
	GetLatestData(ctx context.Context, ts *Timeserie) (*TimeData, error)
}

func Init(tools tools.Tools, db *sql.DB) Service {
	store := newSqlStorage(db)

	return newService(tools, store)
}
