package repository

import (
	"context"
	module "go-sample/internal/modules/catalog/category"
	"go-sample/internal/modules/catalog/category/entity/admin-api"
	repository_shared "go-sample/internal/modules/catalog/category/infrastructure/repository/shared"
	"strings"

	"github.com/google/uuid"
	"github.com/mondegor/go-storage/mrstorage"
	"github.com/mondegor/go-webcore/mrenum"
)

type (
	CategoryPostgres struct {
		client    mrstorage.DBConn
		sqlSelect mrstorage.SqlBuilderSelect
	}
)

func NewCategoryPostgres(
	client mrstorage.DBConn,
	sqlSelect mrstorage.SqlBuilderSelect,
) *CategoryPostgres {
	return &CategoryPostgres{
		client:    client,
		sqlSelect: sqlSelect,
	}
}

func (re *CategoryPostgres) NewFetchParams(params entity.CategoryParams) mrstorage.SqlSelectParams {
	return mrstorage.SqlSelectParams{
		Where: re.sqlSelect.Where(func(w mrstorage.SqlBuilderWhere) mrstorage.SqlBuilderPartFunc {
			return w.JoinAnd(
				w.NotEqual("category_status", mrenum.ItemStatusRemoved),
				w.FilterLike("UPPER(category_caption)", strings.ToUpper(params.Filter.SearchText)),
				w.FilterAnyOf("category_status", params.Filter.Statuses),
			)
		}),
		OrderBy: re.sqlSelect.OrderBy(func(s mrstorage.SqlBuilderOrderBy) mrstorage.SqlBuilderPartFunc {
			return s.Join(
				s.Field(params.Sorter.FieldName, params.Sorter.Direction),
				s.Field("category_id", mrenum.SortDirectionASC),
			)
		}),
		Pager: re.sqlSelect.Pager(func(p mrstorage.SqlBuilderPager) mrstorage.SqlBuilderPartFunc {
			return p.OffsetLimit(params.Pager.Index, params.Pager.Size)
		}),
	}
}

func (re *CategoryPostgres) Fetch(ctx context.Context, params mrstorage.SqlSelectParams) ([]entity.Category, error) {
	whereStr, whereArgs := params.Where.ToSql()

	sql := `
		SELECT
			category_id,
			tag_version,
			created_at as createdAt,
			updated_at as updatedAt,
			category_caption as caption,
			image_meta,
			category_status
		FROM
			` + module.DBSchema + `.categories
		WHERE
			` + whereStr + `
		ORDER BY
			` + params.OrderBy.String() + params.Pager.String() + `;`

	cursor, err := re.client.Query(
		ctx,
		sql,
		whereArgs...,
	)

	if err != nil {
		return nil, err
	}

	defer cursor.Close()

	rows := make([]entity.Category, 0)

	for cursor.Next() {
		var row entity.Category

		err = cursor.Scan(
			&row.ID,
			&row.TagVersion,
			&row.CreatedAt,
			&row.UpdatedAt,
			&row.Caption,
			&row.ImageMeta,
			&row.Status,
		)

		if err != nil {
			return nil, err
		}

		rows = append(rows, row)
	}

	return rows, cursor.Err()
}

func (re *CategoryPostgres) FetchTotal(ctx context.Context, where mrstorage.SqlBuilderPart) (int64, error) {
	whereStr, whereArgs := where.ToSql()

	sql := `
		SELECT
			COUNT(*)
		FROM
			` + module.DBSchema + `.categories
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

func (re *CategoryPostgres) FetchOne(ctx context.Context, rowID uuid.UUID) (entity.Category, error) {
	sql := `
		SELECT
			tag_version,
			created_at,
			updated_at,
			category_caption,
			image_meta,
			category_status
		FROM
			` + module.DBSchema + `.categories
		WHERE
			category_id = $1 AND category_status <> $2
		LIMIT 1;`

	row := entity.Category{ID: rowID}

	err := re.client.QueryRow(
		ctx,
		sql,
		rowID,
		mrenum.ItemStatusRemoved,
	).Scan(
		&row.TagVersion,
		&row.CreatedAt,
		&row.UpdatedAt,
		&row.Caption,
		&row.ImageMeta,
		&row.Status,
	)

	return row, err
}

func (re *CategoryPostgres) FetchStatus(ctx context.Context, row entity.Category) (mrenum.ItemStatus, error) {
	sql := `
		SELECT
			category_status
		FROM
			` + module.DBSchema + `.categories
		WHERE
			category_id = $1 AND category_status <> $2
		LIMIT 1;`

	var status mrenum.ItemStatus

	err := re.client.QueryRow(
		ctx,
		sql,
		row.ID,
		mrenum.ItemStatusRemoved,
	).Scan(
		&status,
	)

	return status, err
}

// IsExists
// result: nil - exists, ErrStorageNoRowFound - not exists, error - query error
func (re *CategoryPostgres) IsExists(ctx context.Context, rowID uuid.UUID) error {
	return repository_shared.CategoryIsExistsPostgres(ctx, re.client, rowID)
}

func (re *CategoryPostgres) Insert(ctx context.Context, row entity.Category) (uuid.UUID, error) {
	sql := `
		INSERT INTO ` + module.DBSchema + `.categories
			(
				category_id,
				category_caption,
				category_status
			)
		VALUES
			(gen_random_uuid(), $1, $2)
		RETURNING
			category_id;`

	err := re.client.QueryRow(
		ctx,
		sql,
		row.Caption,
		row.Status,
	).Scan(
		&row.ID,
	)

	return row.ID, err
}

func (re *CategoryPostgres) Update(ctx context.Context, row entity.Category) (int32, error) {
	sql := `
		UPDATE
			` + module.DBSchema + `.categories
		SET
			tag_version = tag_version + 1,
			updated_at = NOW(),
			category_caption = $4
		WHERE
			category_id = $1 AND tag_version = $2 AND category_status <> $3
		RETURNING
			tag_version;`

	var tagVersion int32

	err := re.client.QueryRow(
		ctx,
		sql,
		row.ID,
		row.TagVersion,
		mrenum.ItemStatusRemoved,
		row.Caption,
	).Scan(
		&tagVersion,
	)

	return tagVersion, err
}

func (re *CategoryPostgres) UpdateStatus(ctx context.Context, row entity.Category) (int32, error) {
	sql := `
		UPDATE
			` + module.DBSchema + `.categories
		SET
			tag_version = tag_version + 1,
			updated_at = NOW(),
			category_status = $4
		WHERE
			category_id = $1 AND tag_version = $2 AND category_status <> $3
		RETURNING
			tag_version;`

	var tagVersion int32

	err := re.client.QueryRow(
		ctx,
		sql,
		row.ID,
		row.TagVersion,
		mrenum.ItemStatusRemoved,
		row.Status,
	).Scan(
		&tagVersion,
	)

	return tagVersion, err
}

func (re *CategoryPostgres) Delete(ctx context.Context, rowID uuid.UUID) error {
	sql := `
		UPDATE
			` + module.DBSchema + `.categories
		SET
			tag_version = tag_version + 1,
			updated_at = NOW(),
			category_status = $2
		WHERE
			category_id = $1 AND category_status <> $2;`

	return re.client.Exec(
		ctx,
		sql,
		rowID,
		mrenum.ItemStatusRemoved,
	)
}
