package sysinfos

import (
	"context"
	"fmt"

	"github.com/Peltoche/zapette/internal/tools"
	"github.com/spf13/afero"
)

type Service interface {
	GetInfos(ctx context.Context) (*Infos, error)
}

func Init(fs afero.Fs, tools tools.Tools) (Service, error) {
	svc := newService(fs, tools)

	err := svc.fetch(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to fetch the sysinfos data: %w", err)
	}

	return svc, nil
}
