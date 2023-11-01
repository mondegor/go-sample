package repository

import (
    "context"

    "github.com/mondegor/go-storage/mrentity"
    "github.com/mondegor/go-storage/mrstorage"
    "github.com/mondegor/go-webcore/mrenum"
)

type (
    CatalogCategoryImage struct {
        client mrstorage.DbConn
    }
)

func NewCatalogCategoryImage(
    client mrstorage.DbConn,
) *CatalogCategoryImage {
    return &CatalogCategoryImage{
        client: client,
    }
}

func (re *CatalogCategoryImage) FetchOne(ctx context.Context, categoryId mrentity.KeyInt32) (string, error) {
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
        mrenum.ItemStatusRemoved,
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
        mrenum.ItemStatusRemoved,
        imagePath,
    )
}

func (re *CatalogCategoryImage) Delete(ctx context.Context, categoryId mrentity.KeyInt32) error {
    return re.Update(ctx, categoryId, "")
}
