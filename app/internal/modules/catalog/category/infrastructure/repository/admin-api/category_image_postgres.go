package repository

import (
	"context"
	module "go-sample/internal/modules/catalog/category"

	"github.com/google/uuid"
	"github.com/mondegor/go-storage/mrentity"
	"github.com/mondegor/go-storage/mrstorage"
	"github.com/mondegor/go-webcore/mrenum"
)

type (
	CategoryImagePostgres struct {
		client mrstorage.DBConn
	}
)

func NewCategoryImagePostgres(
	client mrstorage.DBConn,
) *CategoryImagePostgres {
	return &CategoryImagePostgres{
		client: client,
	}
}

func (re *CategoryImagePostgres) FetchMeta(ctx context.Context, categoryID uuid.UUID) (mrentity.ImageMeta, error) {
	sql := `
		SELECT
			image_meta
		FROM
			` + module.DBSchema + `.categories
		WHERE
			category_id = $1 AND category_status <> $2
		LIMIT 1;`

	var imageMeta mrentity.ImageMeta

	err := re.client.QueryRow(
		ctx,
		sql,
		categoryID,
		mrenum.ItemStatusRemoved,
	).Scan(
		&imageMeta,
	)

	return imageMeta, err
}

func (re *CategoryImagePostgres) UpdateMeta(ctx context.Context, categoryID uuid.UUID, meta mrentity.ImageMeta) error {
	sql := `
		UPDATE
			` + module.DBSchema + `.categories
		SET
			updated_at = NOW(),
			image_meta = $3
		WHERE
			category_id = $1 AND category_status <> $2;`

	return re.client.Exec(
		ctx,
		sql,
		categoryID,
		mrenum.ItemStatusRemoved,
		meta,
	)
}

func (re *CategoryImagePostgres) DeleteMeta(ctx context.Context, categoryID uuid.UUID) error {
	return re.UpdateMeta(ctx, categoryID, mrentity.ImageMeta{})
}
