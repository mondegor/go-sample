package repository

import (
    "context"
    "go-sample/internal/entity/admin-panel"

    "github.com/mondegor/go-components/mrorderer"
    "github.com/mondegor/go-storage/mrentity"
    "github.com/mondegor/go-storage/mrsql"
    "github.com/mondegor/go-storage/mrstorage"
    "github.com/mondegor/go-webcore/mrenum"
)

type (
    CatalogProduct struct {
        client mrstorage.DbConn
        sqlSelect mrstorage.SqlBuilderSelect
        sqlUpdate mrstorage.SqlBuilderUpdate
    }
)

func NewCatalogProduct(
    client mrstorage.DbConn,
    sqlSelect mrstorage.SqlBuilderSelect,
    sqlUpdate mrstorage.SqlBuilderUpdate,
) *CatalogProduct {
    return &CatalogProduct{
        client: client,
        sqlSelect: sqlSelect,
        sqlUpdate: sqlUpdate,
    }
}

func (re *CatalogProduct) GetMetaData(categoryId mrentity.KeyInt32) mrorderer.EntityMeta {
    return mrorderer.NewEntityMeta(
        "public.catalog_products",
        "product_id",
        re.sqlSelect.Where(func (w mrstorage.SqlBuilderWhere) mrstorage.SqlBuilderPartFunc {
            return w.JoinAnd(
                w.Equal("category_id", categoryId),
                w.NotEqual("product_status", mrenum.ItemStatusRemoved),
            )
        }),
    )
}

func (re *CatalogProduct) NewFetchParams(params entity.CatalogProductParams) mrstorage.SqlSelectParams {
    return mrstorage.SqlSelectParams{
        Where: re.sqlSelect.Where(func (w mrstorage.SqlBuilderWhere) mrstorage.SqlBuilderPartFunc {
            return w.JoinAnd(
                w.NotEqual("product_status", mrenum.ItemStatusRemoved),
                w.FilterEqualInt64("category_id", int64(params.Filter.CategoryId), 0),
                w.FilterAnyOf("trademark_id", params.Filter.Trademarks),
                w.FilterLikeFields([]string{"product_article", "product_caption"}, params.Filter.SearchText),
                w.FilterRangeInt64("product_price", params.Filter.Price, 0),
                w.FilterAnyOf("product_status", params.Filter.Statuses),
            )
        }),
        OrderBy: re.sqlSelect.OrderBy(func (s mrstorage.SqlBuilderOrderBy) mrstorage.SqlBuilderPartFunc {
            return s.Join(
                s.Field(s.DbName(params.Sorter.FieldName), params.Sorter.Direction),
                s.Field("product_id", mrentity.SortDirectionASC),
            )
        }),
        Pager: re.sqlSelect.Pager(func (p mrstorage.SqlBuilderPager) mrstorage.SqlBuilderPartFunc {
            return p.OffsetLimit(params.Pager.Index, params.Pager.Size)
        }),
    }
}

func (re *CatalogProduct) Fetch(ctx context.Context, params mrstorage.SqlSelectParams) ([]entity.CatalogProduct, error) {
    whereStr, whereArgs := params.Where.ToSql()

    sql := `
        SELECT
            product_id,
            tag_version,
            datetime_created,
            category_id,
            trademark_id,
            product_article,
            product_caption,
            product_price,
            product_status
        FROM
            public.catalog_products
        WHERE ` + whereStr + `
        ORDER BY ` + params.OrderBy.String() + params.Pager.String() + `;`

    cursor, err := re.client.Query(
        ctx,
        sql,
        whereArgs...
    )

    if err != nil {
        return nil, err
    }

    defer cursor.Close()

    rows := make([]entity.CatalogProduct, 0)

    for cursor.Next() {
        var row entity.CatalogProduct

        err = cursor.Scan(
            &row.Id,
            &row.Version,
            &row.CreatedAt,
            &row.CategoryId,
            &row.TrademarkId,
            &row.Article,
            &row.Caption,
            &row.Price,
            &row.Status,
        )

        rows = append(rows, row)
    }

    return rows, nil
}

func (re *CatalogProduct) FetchTotal(ctx context.Context, where mrstorage.SqlBuilderPart) (int64, error) {
    whereStr, whereArgs := where.ToSql()

    sql := `
        SELECT
            COUNT(*)
        FROM
            public.catalog_products
        WHERE ` + whereStr + `;`

    var totalRow int64

    err := re.client.QueryRow(
        ctx,
        sql,
        whereArgs...,
    ).Scan(
        &totalRow,
    )

    return totalRow, err
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
        WHERE product_id = $1 AND product_status <> $2;`

    return re.client.QueryRow(
        ctx,
        sql,
        row.Id,
        mrenum.ItemStatusRemoved,
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
func (re *CatalogProduct) FetchStatus(ctx context.Context, row *entity.CatalogProduct) (mrenum.ItemStatus, error) {
    sql := `
        SELECT product_status
        FROM
            public.catalog_products
        WHERE product_id = $1 AND tag_version = $2 AND product_status <> $3;`

    var status mrenum.ItemStatus

    err := re.client.QueryRow(
        ctx,
        sql,
        row.Id,
        row.Version,
        mrenum.ItemStatusRemoved,
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
    set, err := re.sqlUpdate.SetFromEntity(row)

    if err != nil {
        return err
    }

    if set.Empty() {
        return nil
    }

    args := []any{
        row.Id,
        row.Version,
        mrenum.ItemStatusRemoved,
    }

    setStr, setArgs := set.Param(len(args) + 1).ToSql()

    sql := `
        UPDATE public.catalog_products
        SET
            tag_version = tag_version + 1,
            ` + setStr + `
        WHERE
            product_id = $1 AND tag_version = $2 AND product_status <> $3;`

    return re.client.Exec(
        ctx,
        sql,
        mrsql.MergeArgs(args, setArgs)...
    )
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

    return re.client.Exec(
        ctx,
        sql,
        row.Id,
        row.Version,
        mrenum.ItemStatusRemoved,
        row.Status,
    )
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

    return re.client.Exec(
        ctx,
        sql,
        id,
        mrenum.ItemStatusRemoved,
    )
}
