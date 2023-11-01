package repository

import (
    "context"
    "go-sample/internal/entity/public"

    "github.com/mondegor/go-storage/mrentity"
    "github.com/mondegor/go-storage/mrstorage"
    "github.com/mondegor/go-webcore/mrenum"
)

type (
    CatalogCategory struct {
        client mrstorage.DbConn
        sqlSelect mrstorage.SqlBuilderSelect
    }
)

func NewCatalogCategory(
    client mrstorage.DbConn,
    sqlSelect mrstorage.SqlBuilderSelect,
) *CatalogCategory {
    return &CatalogCategory{
        client: client,
        sqlSelect: sqlSelect,
    }
}

func (re *CatalogCategory) NewFetchParams(params entity.CatalogCategoryParams) mrstorage.SqlSelectParams {
    return mrstorage.SqlSelectParams{
        Where: re.sqlSelect.Where(func (w mrstorage.SqlBuilderWhere) mrstorage.SqlBuilderPartFunc {
            return w.JoinAnd(
                w.Equal("category_status", mrenum.ItemStatusEnabled),
                w.FilterLike("category_caption", params.Filter.SearchText),
            )
        }),
        OrderBy: re.sqlSelect.OrderBy(func (s mrstorage.SqlBuilderOrderBy) mrstorage.SqlBuilderPartFunc {
            return s.Join(
                s.Field("category_caption", mrentity.SortDirectionASC),
            )
        }),
        Pager: re.sqlSelect.Pager(func (p mrstorage.SqlBuilderPager) mrstorage.SqlBuilderPartFunc {
            return p.OffsetLimit(params.Pager.Index, params.Pager.Size)
        }),
    }
}

func (re *CatalogCategory) Fetch(ctx context.Context, params mrstorage.SqlSelectParams) ([]entity.CatalogCategory, error) {
    whereStr, whereArgs := params.Where.ToSql()

    sql := `
        SELECT
            category_id,
            category_caption,
            image_path
        FROM
            public.catalog_categories
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

    rows := make([]entity.CatalogCategory, 0)

    for cursor.Next() {
        var row entity.CatalogCategory

        err = cursor.Scan(
            &row.Id,
            &row.Caption,
            &row.ImagePath,
        )

        rows = append(rows, row)
    }

    return rows, nil
}

func (re *CatalogCategory) FetchTotal(ctx context.Context, where mrstorage.SqlBuilderPart) (int64, error) {
    whereStr, whereArgs := where.ToSql()

    sql := `
        SELECT
            COUNT(*)
        FROM
            public.catalog_categories
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
// modifies: row{Caption, ImagePath}
func (re *CatalogCategory) LoadOne(ctx context.Context, row *entity.CatalogCategory) error {
    sql := `
        SELECT
            category_caption,
            image_path
        FROM
            public.catalog_categories
        WHERE category_id = $1 AND category_status = $2;`

    return re.client.QueryRow(
        ctx,
        sql,
        row.Id,
        mrenum.ItemStatusEnabled,
    ).Scan(
        &row.Caption,
        &row.ImagePath,
    )
}
