package repository

import (
    "context"

    "github.com/Masterminds/squirrel"
    mrcom_status "github.com/mondegor/go-components/mrcom/status"
    "github.com/mondegor/go-storage/mrentity"
    "github.com/mondegor/go-storage/mrstorage"
)

type (
    CatalogCategoryImage struct {
        client mrstorage.DbConn
        builder squirrel.StatementBuilderType
    }
)

func NewCatalogCategoryImage(
    client mrstorage.DbConn,
    queryBuilder squirrel.StatementBuilderType,
) *CatalogCategoryImage {
    return &CatalogCategoryImage{
        client: client,
        builder: queryBuilder,
    }
}

func (re *CatalogCategoryImage) Fetch(ctx context.Context, categoryId mrentity.KeyInt32) (string, error) {
    sql := `
        SELECT image_path
        FROM
            public.catalog_categories
        WHERE category_id = $1 AND category_status <> $2;`

    var imagePath string

    err := re.client.QueryRow(
        ctx,
        sql,
        categoryId,
        mrcom_status.ItemStatusRemoved,
    ).Scan(
        &imagePath,
    )

    return imagePath, err
}

func (re *CatalogCategoryImage) Update(ctx context.Context, categoryId mrentity.KeyInt32, imagePath string) error {
    sql := `
        UPDATE public.catalog_categories
        SET
            datetime_updated = NOW(),
            image_path = $3
        WHERE category_id = $1 AND category_status <> $2;`

    return re.client.Exec(
        ctx,
        sql,
        categoryId,
        mrcom_status.ItemStatusRemoved,
        imagePath,
    )
}

func (re *CatalogCategoryImage) Delete(ctx context.Context, categoryId mrentity.KeyInt32) error {
    return re.Update(ctx, categoryId, "")
}
