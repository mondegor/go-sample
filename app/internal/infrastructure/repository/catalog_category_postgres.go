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

type CatalogCategory struct {
    client *mrpostgres.ConnAdapter
    builder squirrel.StatementBuilderType
}

func NewCatalogCategory(client *mrpostgres.ConnAdapter,
                        queryBuilder squirrel.StatementBuilderType) *CatalogCategory {
    return &CatalogCategory{
        client: client,
        builder: queryBuilder,
    }
}

func (re *CatalogCategory) LoadAll(ctx context.Context, listFilter *entity.CatalogCategoryListFilter, rows *[]entity.CatalogCategory) error {
    query := re.builder.
        Select(`
            category_id,
            tag_version,
            datetime_created,
            category_caption,
            category_status`).
        From("public.catalog_categories").
        Where(squirrel.NotEq{"category_status": mrcom.ItemStatusRemoved}).
        OrderBy("category_caption ASC, category_id ASC")

    if len(listFilter.Statuses) > 0 {
        query = query.Where(squirrel.Eq{"category_status": listFilter.Statuses})
    }

    cursor, err := re.client.SqQuery(ctx, query)

    if err != nil {
        return err
    }

    for cursor.Next() {
        var row entity.CatalogCategory

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
func (re *CatalogCategory) LoadOne(ctx context.Context, row *entity.CatalogCategory) error {
    sql := `
        SELECT
            tag_version,
            datetime_created,
            category_caption,
            category_status
        FROM
            public.catalog_categories
        WHERE category_id = $1 AND category_status <> $2;`

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
func (re *CatalogCategory) FetchStatus(ctx context.Context, row *entity.CatalogCategory) (mrcom.ItemStatus, error) {
    sql := `
        SELECT category_status
        FROM
            public.catalog_categories
        WHERE category_id = $1 AND tag_version = $2 AND category_status <> $3;`

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
        mrcom.ItemStatusRemoved,
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
            category_caption = $4
        WHERE category_id = $1 AND tag_version = $2 AND category_status <> $3;`

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
func (re *CatalogCategory) UpdateStatus(ctx context.Context, row *entity.CatalogCategory) error {
    sql := `
        UPDATE public.catalog_categories
        SET
            tag_version = tag_version + 1,
            category_status = $4
        WHERE
            category_id = $1 AND tag_version = $2 AND category_status <> $3;`

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

func (re *CatalogCategory) Delete(ctx context.Context, id mrentity.KeyInt32) error {
    sql := `
        UPDATE public.catalog_categories
        SET
            tag_version = tag_version + 1,
            category_status = $2
        WHERE
            category_id = $1 AND category_status <> $2;`

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
