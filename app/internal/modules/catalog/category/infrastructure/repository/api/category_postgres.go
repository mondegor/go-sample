package repository_api

import (
	"context"
	repository_shared "go-sample/internal/modules/catalog/category/infrastructure/repository/shared"

	"github.com/google/uuid"
	"github.com/mondegor/go-storage/mrstorage"
	"github.com/mondegor/go-webcore/mrenum"
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

// FetchStatus
// result: mrenum.ItemStatus - exists, ErrStorageNoRowFound - not exists, error - query error
func (re *CategoryPostgres) FetchStatus(ctx context.Context, rowID uuid.UUID) (mrenum.ItemStatus, error) {
	return repository_shared.CategoryFetchStatusPostgres(ctx, re.client, rowID)
}
