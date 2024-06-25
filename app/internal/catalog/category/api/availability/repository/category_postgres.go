package repository

import (
	"context"

	"github.com/mondegor/go-sample/internal/catalog/category/shared/repository"

	"github.com/google/uuid"
	"github.com/mondegor/go-storage/mrstorage"
	"github.com/mondegor/go-webcore/mrenum"
)

type (
	// CategoryPostgres - comment struct.
	CategoryPostgres struct {
		client mrstorage.DBConnManager
	}
)

// NewCategoryPostgres - comment func.
func NewCategoryPostgres(client mrstorage.DBConnManager) *CategoryPostgres {
	return &CategoryPostgres{
		client: client,
	}
}

// FetchStatus - comment method.
// result: mrenum.ItemStatus - exists, ErrStorageNoRowFound - not exists, error - query error.
func (re *CategoryPostgres) FetchStatus(ctx context.Context, rowID uuid.UUID) (mrenum.ItemStatus, error) {
	return repository.CategoryFetchStatusPostgres(ctx, re.client, rowID)
}
