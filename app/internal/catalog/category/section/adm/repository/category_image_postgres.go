package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/mondegor/go-storage/mrentity"
	"github.com/mondegor/go-storage/mrpostgres/db"
	"github.com/mondegor/go-storage/mrstorage"

	"github.com/mondegor/go-sample/internal/catalog/category/module"
)

type (
	// CategoryImagePostgres - comment struct.
	CategoryImagePostgres struct {
		repoMeta db.FieldUpdater[uuid.UUID, mrentity.ImageMeta]
	}
)

// NewCategoryImagePostgres - создаёт объект CategoryImagePostgres.
func NewCategoryImagePostgres(client mrstorage.DBConnManager) *CategoryImagePostgres {
	return &CategoryImagePostgres{
		repoMeta: db.NewFieldUpdater[uuid.UUID, mrentity.ImageMeta](
			client,
			module.DBTableNameCategories,
			"category_id",
			"image_meta",
			module.DBFieldDeletedAt,
		),
	}
}

// FetchMeta - comment method.
func (re *CategoryImagePostgres) FetchMeta(ctx context.Context, categoryID uuid.UUID) (mrentity.ImageMeta, error) {
	return re.repoMeta.Fetch(ctx, categoryID)
}

// UpdateMeta - comment method.
func (re *CategoryImagePostgres) UpdateMeta(ctx context.Context, categoryID uuid.UUID, meta mrentity.ImageMeta) error {
	return re.repoMeta.Update(ctx, categoryID, meta)
}

// DeleteMeta - comment method.
func (re *CategoryImagePostgres) DeleteMeta(ctx context.Context, categoryID uuid.UUID) error {
	return re.repoMeta.Update(ctx, categoryID, mrentity.ImageMeta{})
}
