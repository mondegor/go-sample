package repository

import (
    "context"
    "go-sample/internal/entity/admin-panel"

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
                w.NotEqual("category_status", mrenum.ItemStatusRemoved),
                w.FilterLike("category_caption", params.Filter.SearchText),
                w.FilterAnyOf("category_status", params.Filter.Statuses),
            )
        }),
        OrderBy: re.sqlSelect.OrderBy(func (s mrstorage.SqlBuilderOrderBy) mrstorage.SqlBuilderPartFunc {
            return s.Join(
                s.Field(s.DbName(params.Sorter.FieldName), params.Sorter.Direction),
                s.Field("category_id", mrentity.SortDirectionASC),
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
            tag_version,
            datetime_created,
            category_caption,
            image_path,
            category_status
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
            &row.Version,
            &row.CreatedAt,
            &row.Caption,
            &row.ImagePath,
            &row.Status,
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
// modifies: row{Version, CreatedAt, Caption, ImagePath, Status}
func (re *CatalogCategory) LoadOne(ctx context.Context, row *entity.CatalogCategory) error {
    sql := `
        SELECT
            tag_version,
            datetime_created,
            category_caption,
            image_path,
            category_status
        FROM
            public.catalog_categories
        WHERE category_id = $1 AND category_status <> $2;`

    return re.client.QueryRow(
        ctx,
        sql,
        row.Id,
        mrenum.ItemStatusRemoved,
    ).Scan(
        &row.Version,
        &row.CreatedAt,
        &row.Caption,
        &row.ImagePath,
        &row.Status,
    )
}

// FetchStatus
// uses: row{Id, Version}
func (re *CatalogCategory) FetchStatus(ctx context.Context, row *entity.CatalogCategory) (mrenum.ItemStatus, error) {
    sql := `
        SELECT category_status
        FROM
            public.catalog_categories
        WHERE category_id = $1 AND tag_version = $2 AND category_status <> $3;`

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

// IsExists
// result: nil - exists, ErrStorageNoRowFound - not exists, error - query error
func (re *CatalogCategory) IsExists(ctx context.Context, id mrentity.KeyInt32) error {
    sql := `
        SELECT 1
        FROM
            public.catalog_categories
        WHERE category_id = $1 AND category_status <> $2;`

    return re.client.QueryRow(
        ctx,
        sql,
        id,
        mrenum.ItemStatusRemoved,
    ).Scan(
        &id,
    )
}

// Insert
// uses: row{Caption, Status}
// modifies: row{Id}
func (re *CatalogCategory) Insert(ctx context.Context, row *entity.CatalogCategory) error {
    sql := `
        INSERT INTO public.catalog_categories
            (category_caption,
             category_status)
        VALUES
            ($1, $2)
        RETURNING category_id;`

    err := re.client.QueryRow(
        ctx,
        sql,
        row.Caption,
        row.Status,
    ).Scan(
        &row.Id,
    )

    return err
}

// Update
// uses: row{Id, Version, Caption, Status}
func (re *CatalogCategory) Update(ctx context.Context, row *entity.CatalogCategory) error {
    sql := `
        UPDATE public.catalog_categories
        SET
            tag_version = tag_version + 1,
            datetime_updated = NOW(),
            category_caption = $4
        WHERE category_id = $1 AND tag_version = $2 AND category_status <> $3;`

    return re.client.Exec(
        ctx,
        sql,
        row.Id,
        row.Version,
        mrenum.ItemStatusRemoved,
        row.Caption,
    )
}

// UpdateStatus
// uses: row{Id, Version, Status}
func (re *CatalogCategory) UpdateStatus(ctx context.Context, row *entity.CatalogCategory) error {
    sql := `
        UPDATE public.catalog_categories
        SET
            tag_version = tag_version + 1,
            datetime_updated = NOW(),
            category_status = $4
        WHERE
            category_id = $1 AND tag_version = $2 AND category_status <> $3;`

    return re.client.Exec(
        ctx,
        sql,
        row.Id,
        row.Version,
        mrenum.ItemStatusRemoved,
        row.Status,
    )
}

func (re *CatalogCategory) Delete(ctx context.Context, id mrentity.KeyInt32) error {
    sql := `
        UPDATE public.catalog_categories
        SET
            tag_version = tag_version + 1,
            datetime_updated = NOW(),
            category_status = $2
        WHERE
            category_id = $1 AND category_status <> $2;`

    return re.client.Exec(
        ctx,
        sql,
        id,
        mrenum.ItemStatusRemoved,
    )
}
