package repository

import (
    "context"
    "go-sample/internal/entity/admin-panel"

    "github.com/mondegor/go-storage/mrentity"
    "github.com/mondegor/go-storage/mrstorage"
    "github.com/mondegor/go-webcore/mrenum"
)

type (
    CatalogTrademark struct {
        client mrstorage.DbConn
        sqlSelect mrstorage.SqlBuilderSelect
    }
)

func NewCatalogTrademark(
    client mrstorage.DbConn,
    sqlSelect mrstorage.SqlBuilderSelect,
) *CatalogTrademark {
    return &CatalogTrademark{
        client: client,
        sqlSelect: sqlSelect,
    }
}

func (re *CatalogTrademark) NewFetchParams(params entity.CatalogTrademarkParams) mrstorage.SqlSelectParams {
    return mrstorage.SqlSelectParams{
        Where: re.sqlSelect.Where(func (w mrstorage.SqlBuilderWhere) mrstorage.SqlBuilderPartFunc {
            return w.JoinAnd(
                w.NotEqual("trademark_status", mrenum.ItemStatusRemoved),
                w.FilterLike("trademark_caption", params.Filter.SearchText),
                w.FilterAnyOf("trademark_status", params.Filter.Statuses),
            )
        }),
        OrderBy: re.sqlSelect.OrderBy(func (s mrstorage.SqlBuilderOrderBy) mrstorage.SqlBuilderPartFunc {
            return s.Join(
                s.Field(s.DbName(params.Sorter.FieldName), params.Sorter.Direction),
                s.Field("trademark_id", mrentity.SortDirectionASC),
            )
        }),
        Pager: re.sqlSelect.Pager(func (p mrstorage.SqlBuilderPager) mrstorage.SqlBuilderPartFunc {
            return p.OffsetLimit(params.Pager.Index, params.Pager.Size)
        }),
    }
}

func (re *CatalogTrademark) Fetch(ctx context.Context, params mrstorage.SqlSelectParams) ([]entity.CatalogTrademark, error) {
    whereStr, whereArgs := params.Where.ToSql()

    sql := `
        SELECT
            trademark_id,
            tag_version,
            datetime_created,
            trademark_caption,
            trademark_status
        FROM
            public.catalog_trademarks
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

    rows := make([]entity.CatalogTrademark, 0)

    for cursor.Next() {
        var row entity.CatalogTrademark

        err = cursor.Scan(
            &row.Id,
            &row.Version,
            &row.CreatedAt,
            &row.Caption,
            &row.Status,
        )

        rows = append(rows, row)
    }

    return rows, nil
}

func (re *CatalogTrademark) FetchTotal(ctx context.Context, where mrstorage.SqlBuilderPart) (int64, error) {
    whereStr, whereArgs := where.ToSql()

    sql := `
        SELECT
            COUNT(*)
        FROM
            public.catalog_trademarks
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
// modifies: row{Version, CreatedAt, Caption, Status}
func (re *CatalogTrademark) LoadOne(ctx context.Context, row *entity.CatalogTrademark) error {
    sql := `
        SELECT
            tag_version,
            datetime_created,
            trademark_caption,
            trademark_status
        FROM
            public.catalog_trademarks
        WHERE trademark_id = $1 AND trademark_status <> $2;`

    return re.client.QueryRow(
        ctx,
        sql,
        row.Id,
        mrenum.ItemStatusRemoved,
    ).Scan(
        &row.Version,
        &row.CreatedAt,
        &row.Caption,
        &row.Status,
    )
}

// FetchStatus
// uses: row{Id, Version}
func (re *CatalogTrademark) FetchStatus(ctx context.Context, row *entity.CatalogTrademark) (mrenum.ItemStatus, error) {
    sql := `
        SELECT trademark_status
        FROM
            public.catalog_trademarks
        WHERE trademark_id = $1 AND tag_version = $2 AND trademark_status <> $3;`

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
func (re *CatalogTrademark) IsExists(ctx context.Context, id mrentity.KeyInt32) error {
    sql := `
        SELECT 1
        FROM
            public.catalog_trademarks
        WHERE trademark_id = $1 AND trademark_status <> $2;`

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
func (re *CatalogTrademark) Insert(ctx context.Context, row *entity.CatalogTrademark) error {
    sql := `
        INSERT INTO public.catalog_trademarks
            (trademark_caption,
             trademark_status)
        VALUES
            ($1, $2)
        RETURNING trademark_id;`

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
func (re *CatalogTrademark) Update(ctx context.Context, row *entity.CatalogTrademark) error {
    sql := `
        UPDATE public.catalog_trademarks
        SET
            tag_version = tag_version + 1,
            datetime_updated = NOW(),
            trademark_caption = $4
        WHERE trademark_id = $1 AND tag_version = $2 AND trademark_status <> $3;`

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
func (re *CatalogTrademark) UpdateStatus(ctx context.Context, row *entity.CatalogTrademark) error {
    sql := `
        UPDATE public.catalog_trademarks
        SET
            tag_version = tag_version + 1,
            datetime_updated = NOW(),
            trademark_status = $4
        WHERE
            trademark_id = $1 AND tag_version = $2 AND trademark_status <> $3;`

    return re.client.Exec(
        ctx,
        sql,
        row.Id,
        row.Version,
        mrenum.ItemStatusRemoved,
        row.Status,
    )
}

func (re *CatalogTrademark) Delete(ctx context.Context, id mrentity.KeyInt32) error {
    sql := `
        UPDATE public.catalog_trademarks
        SET
            tag_version = tag_version + 1,
            datetime_updated = NOW(),
            trademark_status = $2
        WHERE
            trademark_id = $1 AND trademark_status <> $2;`

    return re.client.Exec(
        ctx,
        sql,
        id,
        mrenum.ItemStatusRemoved,
    )
}
