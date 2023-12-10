package repository

import (
	"context"
	module "go-sample/internal/modules/catalog"

	"github.com/mondegor/go-storage/mrstorage"
	"github.com/mondegor/go-webcore/mrenum"
	"github.com/mondegor/go-webcore/mrtype"
)

type (
	CategoryImage struct {
		client mrstorage.DBConn
	}
)

func NewCategoryImage(
	client mrstorage.DBConn,
) *CategoryImage {
	return &CategoryImage{
		client: client,
	}
}

func (re *CategoryImage) FetchPath(ctx context.Context, categoryID mrtype.KeyInt32) (string, error) {
	sql := `
		SELECT
			image_path
		FROM
			` + module.DBSchemaCategory + `.categories
		WHERE
			category_id = $1 AND category_status <> $2
		LIMIT 1;`

	var imagePath string

	err := re.client.QueryRow(
		ctx,
		sql,
		categoryID,
		mrenum.ItemStatusRemoved,
	).Scan(
		&imagePath,
	)

	return imagePath, err
}

func (re *CategoryImage) Update(ctx context.Context, categoryID mrtype.KeyInt32, imagePath string) error {
	sql := `
		UPDATE
			` + module.DBSchemaCategory + `.categories
		SET
			datetime_updated = NOW(),
			image_path = $3
		WHERE
			category_id = $1 AND category_status <> $2;`

	return re.client.Exec(
		ctx,
		sql,
		categoryID,
		mrenum.ItemStatusRemoved,
		imagePath,
	)
}

func (re *CategoryImage) Delete(ctx context.Context, categoryID mrtype.KeyInt32) error {
	return re.Update(ctx, categoryID, "")
}
