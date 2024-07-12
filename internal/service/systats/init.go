package systats

import (
	"context"

	"github.com/spf13/afero"
)

type Service interface {
	FetchMeminfos(ctx context.Context) (*Memory, error)
}

func Init(fs afero.Fs) Service {
	return newService(fs)
}
