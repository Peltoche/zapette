package sqlstorage

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/Peltoche/zapette/internal/tools"
	"go.uber.org/fx"
)

type Result struct {
	fx.Out
	Querier    Querier
	DB         *sql.DB
	Transactor Transactor
}

// Transactor runs logic inside a single database transaction
type Transactor interface {
	// WithinTransaction runs a function within a database transaction.
	//
	// Transaction is propagated in the context,
	// so it is important to propagate it to underlying repositories.
	// Function commits if error is nil, and rollbacks if not.
	// It returns the same error.
	WithinTransaction(context.Context, func(ctx context.Context) error) error
}

func Init(hooks []SQLChangeHook, cfg Config, tools tools.Tools) (Result, error) {
	fmt.Printf("run init \n#####\n########\n")
	db, err := NewSQliteClient(&cfg, hooks, tools.Logger())
	if err != nil {
		return Result{}, fmt.Errorf("sqlite error: %w", err)
	}

	return Result{
		DB:         db,
		Querier:    db,
		Transactor: NewTransacGenerator(db, tools),
	}, nil
}
