package repository

import (
    "context"

    "github.com/Masterminds/squirrel"
    "github.com/mondegor/go-components/mrcom"
    "github.com/mondegor/go-storage/mrentity"
    "github.com/mondegor/go-storage/mrpostgres"
    "github.com/mondegor/go-webcore/mrcore"
)

type (
    CatalogCategoryImage struct {
        client *mrpostgres.ConnAdapter
        builder squirrel.StatementBuilderType
    }
)

func NewCatalogCategoryImage(client *mrpostgres.ConnAdapter,
                             queryBuilder squirrel.StatementBuilderType) *CatalogCategoryImage {
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
        mrcom.ItemStatusRemoved,
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

    commandTag, err := re.client.Exec(
        ctx,
        sql,
        categoryId,
        mrcom.ItemStatusRemoved,
        imagePath,
    )

    if err != nil {
        return err
    }

    if commandTag.RowsAffected() < 1 {
        return mrcore.FactoryErrStorageRowsNotAffected.New()
    }

    return nil
}

func (re *CatalogCategoryImage) Delete(ctx context.Context, categoryId mrentity.KeyInt32) error {
    return re.Update(ctx, categoryId, "")
}
