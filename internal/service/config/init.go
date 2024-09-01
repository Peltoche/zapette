package config

import (
	"context"
	"database/sql"

	"github.com/Peltoche/zapette/internal/tools"
	"github.com/Peltoche/zapette/internal/tools/uuid"
)

type Service interface {
	SetSysstatInputNamespace(ctx context.Context, id uuid.UUID) error
	GetSysstatInputNamespace(ctx context.Context) (*uuid.UUID, error)
}

func Init(db *sql.DB, tools tools.Tools) Service {
	storage := newSqlStorage(db)

	return newService(storage, tools)
}
