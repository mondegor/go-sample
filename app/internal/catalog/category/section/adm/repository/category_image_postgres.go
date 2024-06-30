package repository

import (
	"context"

	"github.com/mondegor/go-sample/internal/catalog/category/module"

	"github.com/google/uuid"
	"github.com/mondegor/go-storage/mrentity"
	"github.com/mondegor/go-storage/mrstorage"
)

type (
	// CategoryImagePostgres - comment struct.
	CategoryImagePostgres struct {
		client mrstorage.DBConnManager
	}
)

// NewCategoryImagePostgres - создаёт объект CategoryImagePostgres.
func NewCategoryImagePostgres(client mrstorage.DBConnManager) *CategoryImagePostgres {
	return &CategoryImagePostgres{
		client: client,
	}
}

// FetchMeta - comment method.
func (re *CategoryImagePostgres) FetchMeta(ctx context.Context, categoryID uuid.UUID) (mrentity.ImageMeta, error) {
	sql := `
		SELECT
			image_meta
		FROM
			` + module.DBSchema + `.` + module.DBTableNameCategories + `
		WHERE
			category_id = $1 AND deleted_at IS NULL
		LIMIT 1;`

	var imageMeta mrentity.ImageMeta

	err := re.client.Conn(ctx).QueryRow(
		ctx,
		sql,
		categoryID,
	).Scan(
		&imageMeta,
	)

	return imageMeta, err
}

// UpdateMeta - comment method.
func (re *CategoryImagePostgres) UpdateMeta(ctx context.Context, categoryID uuid.UUID, meta mrentity.ImageMeta) error {
	sql := `
		UPDATE
			` + module.DBSchema + `.` + module.DBTableNameCategories + `
		SET
			updated_at = NOW(),
			image_meta = $2
		WHERE
			category_id = $1 AND deleted_at IS NULL;`

	return re.client.Conn(ctx).Exec(
		ctx,
		sql,
		categoryID,
		meta,
	)
}

// DeleteMeta - comment method.
func (re *CategoryImagePostgres) DeleteMeta(ctx context.Context, categoryID uuid.UUID) error {
	return re.UpdateMeta(ctx, categoryID, mrentity.ImageMeta{})
}
