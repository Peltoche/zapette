package config

import (
	"context"
	"errors"
	"fmt"

	"github.com/Peltoche/zapette/internal/tools"
	"github.com/Peltoche/zapette/internal/tools/errs"
	"github.com/Peltoche/zapette/internal/tools/uuid"
)

type storage interface {
	Save(ctx context.Context, key ConfigKey, value string) error
	Get(ctx context.Context, key ConfigKey) (string, error)
}

type service struct {
	storage storage
	uuid    uuid.Service
}

func newService(storage storage, tools tools.Tools) *service {
	return &service{
		storage: storage,
		uuid:    tools.UUID(),
	}
}

func (s *service) SetSysstatInputNamespace(ctx context.Context, id uuid.UUID) error {
	err := s.storage.Save(ctx, sysstatsInputNamespace, string(id))
	if err != nil {
		return fmt.Errorf("failed to Save: %w", err)
	}

	return nil
}

func (s *service) GetSysstatInputNamespace(ctx context.Context) (*uuid.UUID, error) {
	idStr, err := s.storage.Get(ctx, sysstatsInputNamespace)
	if errors.Is(err, errNotfound) {
		return nil, errs.ErrNotFound
	}

	if err != nil {
		return nil, fmt.Errorf("failed to Get: %w", err)
	}

	res, err := s.uuid.Parse(idStr)
	if err != nil {
		return nil, fmt.Errorf("invalid id format: %w", err)
	}

	return &res, nil
}
