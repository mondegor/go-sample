package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/mondegor/go-storage/mrpostgres/db"
	"github.com/mondegor/go-storage/mrstorage"
	"github.com/mondegor/go-webcore/mrenum"

	"github.com/mondegor/go-sample/internal/catalog/category/module"
)

type (
	// CategoryPostgres - comment struct.
	CategoryPostgres struct {
		repoStatus db.FieldFetcher[uuid.UUID, mrenum.ItemStatus]
	}
)

// NewCategoryPostgres - создаёт объект CategoryPostgres.
func NewCategoryPostgres(client mrstorage.DBConnManager) *CategoryPostgres {
	return &CategoryPostgres{
		repoStatus: db.NewFieldFetcher[uuid.UUID, mrenum.ItemStatus](
			client,
			module.DBTableNameCategories,
			"category_id",
			"category_status",
			module.DBFieldDeletedAt,
		),
	}
}

// FetchStatus - comment method.
// result: mrenum.ItemStatus - exists, ErrStorageNoRowFound - not exists, error - query error.
func (re *CategoryPostgres) FetchStatus(ctx context.Context, rowID uuid.UUID) (mrenum.ItemStatus, error) {
	return re.repoStatus.Fetch(ctx, rowID)
}
