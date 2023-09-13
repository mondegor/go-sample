package repository

import (
    "context"
    "go-sample/internal/entity"

    "github.com/Masterminds/squirrel"
    "github.com/mondegor/go-components/mrcom"
    "github.com/mondegor/go-storage/mrentity"
    "github.com/mondegor/go-storage/mrpostgres"
    "github.com/mondegor/go-webcore/mrcore"
)

type CatalogTrademark struct {
    client *mrpostgres.ConnAdapter
    builder squirrel.StatementBuilderType
}

func NewCatalogTrademark(client *mrpostgres.ConnAdapter,
                         queryBuilder squirrel.StatementBuilderType) *CatalogTrademark {
    return &CatalogTrademark{
        client: client,
        builder: queryBuilder,
    }
}

func (re *CatalogTrademark) LoadAll(ctx context.Context, listFilter *entity.CatalogTrademarkListFilter, rows *[]entity.CatalogTrademark) error {
    query := re.builder.
        Select(`
            trademark_id,
            tag_version,
            datetime_created,
            trademark_caption,
            trademark_status`).
        From("public.catalog_trademarks").
        Where(squirrel.NotEq{"trademark_status": mrcom.ItemStatusRemoved}).
        OrderBy("trademark_caption ASC, trademark_id ASC")

    if len(listFilter.Statuses) > 0 {
        query = query.Where(squirrel.Eq{"trademark_status": listFilter.Statuses})
    }

    cursor, err := re.client.SqQuery(ctx, query)

    if err != nil {
        return err
    }

    for cursor.Next() {
        var row entity.CatalogTrademark

        err = cursor.Scan(
            &row.Id,
            &row.Version,
            &row.CreatedAt,
            &row.Caption,
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
        mrcom.ItemStatusRemoved,
    ).Scan(
        &row.Version,
        &row.CreatedAt,
        &row.Caption,
        &row.Status,
    )
}

// FetchStatus
// uses: row{Id, Version}
func (re *CatalogTrademark) FetchStatus(ctx context.Context, row *entity.CatalogTrademark) (mrcom.ItemStatus, error) {
    sql := `
        SELECT trademark_status
        FROM
            public.catalog_trademarks
        WHERE trademark_id = $1 AND tag_version = $2 AND trademark_status <> $3;`

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
        mrcom.ItemStatusRemoved,
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
            trademark_caption = $4
        WHERE trademark_id = $1 AND tag_version = $2 AND trademark_status <> $3;`

    commandTag, err := re.client.Exec(
        ctx,
        sql,
        row.Id,
        row.Version,
        mrcom.ItemStatusRemoved,
        row.Caption,
    )

    if err != nil {
        return err
    }

    if commandTag.RowsAffected() < 1 {
        return mrcore.FactoryErrStorageRowsNotAffected.New()
    }

    return nil
}

// UpdateStatus
// uses: row{Id, Version, Status}
func (re *CatalogTrademark) UpdateStatus(ctx context.Context, row *entity.CatalogTrademark) error {
    sql := `
        UPDATE public.catalog_trademarks
        SET
            tag_version = tag_version + 1,
            trademark_status = $4
        WHERE
            trademark_id = $1 AND tag_version = $2 AND trademark_status <> $3;`

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

func (re *CatalogTrademark) Delete(ctx context.Context, id mrentity.KeyInt32) error {
    sql := `
        UPDATE public.catalog_trademarks
        SET
            tag_version = tag_version + 1,
            trademark_status = $2
        WHERE
            trademark_id = $1 AND trademark_status <> $2;`

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
