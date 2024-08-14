package users

import (
	"context"
	"database/sql"

	"github.com/Peltoche/zapette/internal/tools"
	"github.com/Peltoche/zapette/internal/tools/secret"
	"github.com/Peltoche/zapette/internal/tools/sqlstorage"
	"github.com/Peltoche/zapette/internal/tools/uuid"
)

type Service interface {
	Create(ctx context.Context, user *CreateCmd) (*User, error)
	Bootstrap(ctx context.Context, cmd *BootstrapCmd) (*User, error)
	GetByID(ctx context.Context, userID uuid.UUID) (*User, error)
	Authenticate(ctx context.Context, username string, password secret.Text) (*User, error)
	GetAll(ctx context.Context, paginateCmd *sqlstorage.PaginateCmd) ([]User, error)
	AddToDeletion(ctx context.Context, userID uuid.UUID) error
	HardDelete(ctx context.Context, userID uuid.UUID) error
	GetAllWithStatus(ctx context.Context, status Status, cmd *sqlstorage.PaginateCmd) ([]User, error)
	MarkInitAsFinished(ctx context.Context, userID uuid.UUID) (*User, error)
	UpdateUserPassword(ctx context.Context, cmd *UpdatePasswordCmd) error
}

func Init(
	tools tools.Tools,
	db *sql.DB,
) Service {
	store := newSqlStorage(db)

	return newService(tools, store)
}
