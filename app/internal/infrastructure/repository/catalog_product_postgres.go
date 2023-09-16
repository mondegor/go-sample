package repository

import (
    "context"
    "go-sample/internal/entity"

    "github.com/Masterminds/squirrel"
    "github.com/mondegor/go-components/mrcom"
    mrcom_orderer "github.com/mondegor/go-components/mrcom/orderer"
    "github.com/mondegor/go-storage/mrentity"
    "github.com/mondegor/go-storage/mrpostgres"
    "github.com/mondegor/go-webcore/mrcore"
)

type (
    CatalogProduct struct {
        client *mrpostgres.ConnAdapter
        builder squirrel.StatementBuilderType
    }
)

func NewCatalogProduct(client *mrpostgres.ConnAdapter,
                       queryBuilder squirrel.StatementBuilderType) *CatalogProduct {
    return &CatalogProduct{
        client: client,
        builder: queryBuilder,
    }
}

func (re *CatalogProduct) GetMetaData(categoryId mrentity.KeyInt32) mrcom_orderer.EntityMeta {
    return mrcom_orderer.NewEntityMeta(
        "public.catalog_products",
        "product_id",
        []any{
            squirrel.Eq{"category_id": categoryId},
            squirrel.NotEq{"product_status": mrcom.ItemStatusRemoved},
        },
    )
}

func (re *CatalogProduct) LoadAll(ctx context.Context, listFilter *entity.CatalogProductListFilter, rows *[]entity.CatalogProduct) error {
    query := re.builder.
        Select(`
            product_id,
            tag_version,
            datetime_created,
            trademark_id,
            product_article,
            product_caption,
            product_price,
            product_status`).
        From("public.catalog_products").
        Where(squirrel.Eq{"category_id": listFilter.CategoryId}).
        Where(squirrel.NotEq{"product_status": mrcom.ItemStatusRemoved}).
        OrderBy("order_field ASC, product_caption ASC, product_id ASC")

    if len(listFilter.Trademarks) > 0 {
        query = query.Where(squirrel.Eq{"trademark_id": listFilter.Trademarks})
    }

    if len(listFilter.Statuses) > 0 {
        query = query.Where(squirrel.Eq{"product_status": listFilter.Statuses})
    }

    cursor, err := re.client.SqQuery(ctx, query)

    if err != nil {
        return err
    }

    defer cursor.Close()

    for cursor.Next() {
        row := entity.CatalogProduct{CategoryId: listFilter.CategoryId}

        err = cursor.Scan(
            &row.Id,
            &row.Version,
            &row.CreatedAt,
            &row.TrademarkId,
            &row.Article,
            &row.Caption,
            &row.Price,
            &row.Status,
        )

        if err != nil {
            return mrcore.FactoryErrStorageFetchDataFailed.Wrap(err)
        }

        *rows = append(*rows, row)
    }

    if err = cursor.Err(); err != nil {
        return mrcore.FactoryErrStorageFetchDataFailed.Wrap(err)
    }

    return nil
}

// LoadOne
// uses: row{Id}
// modifies: row{Version, CreatedAt, CategoryId, TrademarkId, Article, Caption, Price, Status}
func (re *CatalogProduct) LoadOne(ctx context.Context, row *entity.CatalogProduct) error {
    sql := `
        SELECT
            tag_version,
            datetime_created,
            trademark_id,
            product_article,
            product_caption,
            product_price,
            product_status
        FROM
            public.catalog_products
        WHERE product_id = $1 AND category_id = $2 AND product_status <> $3;`

    return re.client.QueryRow(
        ctx,
        sql,
        row.Id,
        row.CategoryId,
        mrcom.ItemStatusRemoved,
    ).Scan(
        &row.Version,
        &row.CreatedAt,
        &row.TrademarkId,
        &row.Article,
        &row.Caption,
        &row.Price,
        &row.Status,
    )
}

func (re *CatalogProduct) FetchIdByArticle(ctx context.Context, article string) (mrentity.KeyInt32, error) {
    sql := `
        SELECT product_id
        FROM
            public.catalog_products
        WHERE product_article = $1;`

    var id mrentity.KeyInt32

    err := re.client.QueryRow(
        ctx,
        sql,
        article,
    ).Scan(
        &id,
    )

    return id, err
}

// FetchStatus
// uses: row{Id, Version}
func (re *CatalogProduct) FetchStatus(ctx context.Context, row *entity.CatalogProduct) (mrcom.ItemStatus, error) {
    sql := `
        SELECT product_status
        FROM
            public.catalog_products
        WHERE product_id = $1 AND tag_version = $2 AND product_status <> $3;`

    var status mrcom.ItemStatus

    err := re.client.QueryRow(
        ctx,
        sql,
        row.Id,
        row.Version,
        mrcom.ItemStatusRemoved,
    ).Scan(
        &status,
    )

    return status, err
}

// Insert
// uses: row{CategoryId, TrademarkId, Article, Caption, Price, Status}
// modifies: row{Id}
func (re *CatalogProduct) Insert(ctx context.Context, row *entity.CatalogProduct) error {
    sql := `
        INSERT INTO public.catalog_products
            (category_id,
             trademark_id,
             product_article,
             product_caption,
             product_price,
             product_status)
        VALUES
            ($1, $2, $3, $4, $5, $6)
        RETURNING product_id;`

    err := re.client.QueryRow(
        ctx,
        sql,
        row.CategoryId,
        row.TrademarkId,
        row.Article,
        row.Caption,
        row.Price,
        row.Status,
    ).Scan(
        &row.Id,
    )

    return err
}

// Update
// uses: row{Id, Version, TrademarkId, Article, Caption, Price, Status}
func (re *CatalogProduct) Update(ctx context.Context, row *entity.CatalogProduct) error {
    filledFields, err := mrentity.FilledFieldsToUpdate(row)

    if err != nil {
        return err
    }

    query := re.builder.
        Update("public.catalog_products").
        Set("tag_version", squirrel.Expr("tag_version + 1")).
        SetMap(filledFields).
        Where(squirrel.Eq{"product_id": row.Id}).
        Where(squirrel.Eq{"tag_version": row.Version}).
        Where(squirrel.NotEq{"product_status": mrcom.ItemStatusRemoved})

    return re.client.SqUpdate(ctx, query)
}

// UpdateStatus
// uses: row{Id, Version, Status}
func (re *CatalogProduct) UpdateStatus(ctx context.Context, row *entity.CatalogProduct) error {
    sql := `
        UPDATE public.catalog_products
        SET
            tag_version = tag_version + 1,
            datetime_updated = NOW(),
            product_status = $4
        WHERE
            product_id = $1 AND tag_version = $2 AND product_status <> $3;`

    commandTag, err := re.client.Exec(
        ctx,
        sql,
        row.Id,
        row.Version,
        mrcom.ItemStatusRemoved,
        row.Status,
    )

    if err != nil {
        return err
    }

    if commandTag.RowsAffected() < 1 {
        return mrcore.FactoryErrStorageRowsNotAffected.New()
    }

    return nil
}

func (re *CatalogProduct) Delete(ctx context.Context, id mrentity.KeyInt32) error {
    sql := `
        UPDATE public.catalog_products
        SET
            tag_version = tag_version + 1,
            datetime_updated = NOW(),
            product_article = NULL,
            prev_field_id = NULL,
            next_field_id = NULL,
            order_field = NULL,
            product_status = $2
        WHERE
            product_id = $1 AND product_status <> $2;`

    commandTag, err := re.client.Exec(
        ctx,
        sql,
        id,
        mrcom.ItemStatusRemoved,
    )

    if err != nil {
        return err
    }

    if commandTag.RowsAffected() < 1 {
        return mrcore.FactoryErrStorageRowsNotAffected.New()
    }

    return nil
}
