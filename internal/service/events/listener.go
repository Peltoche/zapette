package events

import (
	"context"

	"github.com/Peltoche/zapette/internal/tools/sqlstorage"
)

type EventListener struct{}

func New() *EventListener {
	return &EventListener{}
}

func (l *EventListener) Name() string {
	return "events"
}

func (l *EventListener) ShouldRunHook(table string) bool {
	return table == "sysstats"
}

func (l *EventListener) RunHook(ctx context.Context, db sqlstorage.Querier, table string) error {
	// Does nothing yet
	return nil
}
