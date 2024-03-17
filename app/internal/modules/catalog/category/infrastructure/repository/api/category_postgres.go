package repository_api

import (
	"context"
	repository_shared "go-sample/internal/modules/catalog/category/infrastructure/repository/shared"

	"github.com/google/uuid"
	"github.com/mondegor/go-storage/mrstorage"
)

type (
	CategoryPostgres struct {
		client mrstorage.DBConn
	}
)

func NewCategoryPostgres(
	client mrstorage.DBConn,
) *CategoryPostgres {
	return &CategoryPostgres{
		client: client,
	}
}

// IsExists
// result: nil - exists, ErrStorageNoRowFound - not exists, error - query error
func (re *CategoryPostgres) IsExists(ctx context.Context, rowID uuid.UUID) error {
	return repository_shared.CategoryIsExistsPostgres(ctx, re.client, rowID)
}
