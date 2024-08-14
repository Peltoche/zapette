package sqlstorage

import (
	"database/sql"
	"fmt"

	"github.com/Peltoche/zapette/internal/tools"
	"go.uber.org/fx"
)

type Result struct {
	fx.Out
	DB *sql.DB
}

func Init(cfg Config, tools tools.Tools) (Result, error) {
	db, err := NewSQliteClient(&cfg)
	if err != nil {
		return Result{}, fmt.Errorf("sqlite error: %w", err)
	}

	return Result{
		DB: db,
	}, nil
}
