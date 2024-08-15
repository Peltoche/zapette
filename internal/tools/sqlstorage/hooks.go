package sqlstorage

import (
	"context"
	"log/slog"

	"github.com/Peltoche/zapette/internal/tools"
)

type SQLChangeHook interface {
	SQLHookName() string
	RunSQLHook(ctx context.Context, table string) error
}

type SQLChangeHookList struct {
	inner  []SQLChangeHook
	logger *slog.Logger
}

func NewSQLChangeHookList(tools tools.Tools) *SQLChangeHookList {
	return &SQLChangeHookList{
		inner:  []SQLChangeHook{},
		logger: tools.Logger(),
	}
}

func (l *SQLChangeHookList) GetHooks() []SQLChangeHook {
	return l.inner
}

func (l *SQLChangeHookList) AddHooks(hooks ...SQLChangeHook) {
	l.inner = append(l.inner, hooks...)

	for _, hook := range hooks {
		l.logger.Debug("register SQL hook", slog.String("hook", hook.SQLHookName()))
	}
}
