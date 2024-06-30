package password

import (
	"context"

	"github.com/Peltoche/zapette/internal/tools/secret"
)

type Password interface {
	Encrypt(ctx context.Context, password secret.Text) (secret.Text, error)
	Compare(ctx context.Context, hash, password secret.Text) (bool, error)
}
