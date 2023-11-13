package repository

import (
    "context"
    "go-sample/internal/modules/catalog/entity/admin-api"
    "strings"

    "github.com/mondegor/go-storage/mrstorage"
    "github.com/mondegor/go-webcore/mrenum"
    "github.com/mondegor/go-webcore/mrtype"
)

type (
    Trademark struct {
        client mrstorage.DBConn
        sqlSelect mrstorage.SqlBuilderSelect
    }
)

func NewTrademark(
    client mrstorage.DBConn,
    sqlSelect mrstorage.SqlBuilderSelect,
) *Trademark {
    return &Trademark{
        client: client,
        sqlSelect: sqlSelect,
    }
}

func (re *Trademark) NewFetchParams(params entity.TrademarkParams) mrstorage.SqlSelectParams {
    return mrstorage.SqlSelectParams{
        Where: re.sqlSelect.Where(func (w mrstorage.SqlBuilderWhere) mrstorage.SqlBuilderPartFunc {
            return w.JoinAnd(
                w.NotEqual("trademark_status", mrenum.ItemStatusRemoved),
                w.FilterLike("UPPER(trademark_caption)", strings.ToUpper(params.Filter.SearchText)),
                w.FilterAnyOf("trademark_status", params.Filter.Statuses),
            )
        }),
        OrderBy: re.sqlSelect.OrderBy(func (s mrstorage.SqlBuilderOrderBy) mrstorage.SqlBuilderPartFunc {
            return s.Join(
                s.Field(params.Sorter.FieldName, params.Sorter.Direction),
                s.Field("trademark_id", mrenum.SortDirectionASC),
            )
        }),
        Pager: re.sqlSelect.Pager(func (p mrstorage.SqlBuilderPager) mrstorage.SqlBuilderPartFunc {
            return p.OffsetLimit(params.Pager.Index, params.Pager.Size)
        }),
    }
}

func (re *Trademark) Fetch(ctx context.Context, params mrstorage.SqlSelectParams) ([]entity.Trademark, error) {
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
        WHERE
            ` + whereStr + `
        ORDER BY
            ` + params.OrderBy.String() + params.Pager.String() + `;`

    cursor, err := re.client.Query(
        ctx,
        sql,
        whereArgs...
    )

    if err != nil {
        return nil, err
    }

    defer cursor.Close()

    rows := make([]entity.Trademark, 0)

    for cursor.Next() {
        var row entity.Trademark

        err = cursor.Scan(
            &row.ID,
            &row.TagVersion,
            &row.CreatedAt,
            &row.Caption,
            &row.Status,
        )

        if err != nil {
            return nil, err
        }

        rows = append(rows, row)
    }

    return rows, nil
}

func (re *Trademark) FetchTotal(ctx context.Context, where mrstorage.SqlBuilderPart) (int64, error) {
    whereStr, whereArgs := where.ToSql()

    sql := `
        SELECT
            COUNT(*)
        FROM
            public.catalog_trademarks
        WHERE
            ` + whereStr + `;`

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

func (re *Trademark) LoadOne(ctx context.Context, row *entity.Trademark) error {
    sql := `
        SELECT
            tag_version,
            datetime_created,
            trademark_caption,
            trademark_status
        FROM
            public.catalog_trademarks
        WHERE
            trademark_id = $1 AND trademark_status <> $2
        LIMIT 1;`

    return re.client.QueryRow(
        ctx,
        sql,
        row.ID,
        mrenum.ItemStatusRemoved,
    ).Scan(
        &row.TagVersion,
        &row.CreatedAt,
        &row.Caption,
        &row.Status,
    )
}

func (re *Trademark) FetchStatus(ctx context.Context, row *entity.Trademark) (mrenum.ItemStatus, error) {
    sql := `
        SELECT
            trademark_status
        FROM
            public.catalog_trademarks
        WHERE
            trademark_id = $1 AND tag_version = $2 AND trademark_status <> $3
        LIMIT 1;`

    var status mrenum.ItemStatus

    err := re.client.QueryRow(
        ctx,
        sql,
        row.ID,
        row.TagVersion,
        mrenum.ItemStatusRemoved,
    ).Scan(
        &status,
    )

    return status, err
}

// IsExists
// result: nil - exists, ErrStorageNoRowFound - not exists, error - query error
func (re *Trademark) IsExists(ctx context.Context, id mrtype.KeyInt32) error {
    sql := `
        SELECT
            1
        FROM
            public.catalog_trademarks
        WHERE
            trademark_id = $1 AND trademark_status <> $2
        LIMIT 1;`

    return re.client.QueryRow(
        ctx,
        sql,
        id,
        mrenum.ItemStatusRemoved,
    ).Scan(
        &id,
    )
}

func (re *Trademark) Insert(ctx context.Context, row *entity.Trademark) error {
    sql := `
        INSERT INTO public.catalog_trademarks
            (
                trademark_caption,
                trademark_status
            )
        VALUES
            ($1, $2)
        RETURNING
            trademark_id;`

    err := re.client.QueryRow(
        ctx,
        sql,
        row.Caption,
        row.Status,
    ).Scan(
        &row.ID,
    )

    return err
}

func (re *Trademark) Update(ctx context.Context, row *entity.Trademark) error {
    sql := `
        UPDATE
            public.catalog_trademarks
        SET
            tag_version = tag_version + 1,
            datetime_updated = NOW(),
            trademark_caption = $4
        WHERE
            trademark_id = $1 AND tag_version = $2 AND trademark_status <> $3;`

    return re.client.Exec(
        ctx,
        sql,
        row.ID,
        row.TagVersion,
        mrenum.ItemStatusRemoved,
        row.Caption,
    )
}

func (re *Trademark) UpdateStatus(ctx context.Context, row *entity.Trademark) error {
    sql := `
        UPDATE
            public.catalog_trademarks
        SET
            tag_version = tag_version + 1,
            datetime_updated = NOW(),
            trademark_status = $4
        WHERE
            trademark_id = $1 AND tag_version = $2 AND trademark_status <> $3;`

    return re.client.Exec(
        ctx,
        sql,
        row.ID,
        row.TagVersion,
        mrenum.ItemStatusRemoved,
        row.Status,
    )
}

func (re *Trademark) Delete(ctx context.Context, id mrtype.KeyInt32) error {
    sql := `
        UPDATE
            public.catalog_trademarks
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
