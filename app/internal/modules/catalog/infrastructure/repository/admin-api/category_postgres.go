package repository

import (
	"context"
	module "go-sample/internal/modules/catalog"
	"go-sample/internal/modules/catalog/entity/admin-api"
	"strings"

	"github.com/mondegor/go-storage/mrstorage"
	"github.com/mondegor/go-webcore/mrenum"
	"github.com/mondegor/go-webcore/mrtype"
)

type (
	Category struct {
		client    mrstorage.DBConn
		sqlSelect mrstorage.SqlBuilderSelect
	}
)

func NewCategory(
	client mrstorage.DBConn,
	sqlSelect mrstorage.SqlBuilderSelect,
) *Category {
	return &Category{
		client:    client,
		sqlSelect: sqlSelect,
	}
}

func (re *Category) NewFetchParams(params entity.CategoryParams) mrstorage.SqlSelectParams {
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

func (re *Category) Fetch(ctx context.Context, params mrstorage.SqlSelectParams) ([]entity.Category, error) {
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
			` + module.DBSchema + `.catalog_categories
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
			&row.Caption,
			&row.ImagePath,
			&row.Status,
		)

		if err != nil {
			return nil, err
		}

		rows = append(rows, row)
	}

	return rows, nil
}

func (re *Category) FetchTotal(ctx context.Context, where mrstorage.SqlBuilderPart) (int64, error) {
	whereStr, whereArgs := where.ToSql()

	sql := `
		SELECT
			COUNT(*)
		FROM
			` + module.DBSchema + `.catalog_categories
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

func (re *Category) LoadOne(ctx context.Context, row *entity.Category) error {
	sql := `
		SELECT
			tag_version,
			datetime_created,
			category_caption,
			image_path,
			category_status
		FROM
			` + module.DBSchema + `.catalog_categories
		WHERE
			category_id = $1 AND category_status <> $2
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
		&row.ImagePath,
		&row.Status,
	)
}

func (re *Category) FetchStatus(ctx context.Context, row *entity.Category) (mrenum.ItemStatus, error) {
	sql := `
		SELECT
			category_status
		FROM
			` + module.DBSchema + `.catalog_categories
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
func (re *Category) IsExists(ctx context.Context, id mrtype.KeyInt32) error {
	sql := `
		SELECT
			1
		FROM
			` + module.DBSchema + `.catalog_categories
		WHERE
			category_id = $1 AND category_status <> $2
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

func (re *Category) Insert(ctx context.Context, row *entity.Category) error {
	sql := `
		INSERT INTO ` + module.DBSchema + `.catalog_categories
			(
				category_caption,
				category_status
			)
		VALUES
			($1, $2)
		RETURNING
			category_id;`

	return re.client.QueryRow(
		ctx,
		sql,
		row.Caption,
		row.Status,
	).Scan(
		&row.ID,
	)
}

func (re *Category) Update(ctx context.Context, row *entity.Category) (int32, error) {
	sql := `
		UPDATE
			` + module.DBSchema + `.catalog_categories
		SET
			tag_version = tag_version + 1,
			datetime_updated = NOW(),
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

func (re *Category) UpdateStatus(ctx context.Context, row *entity.Category) (int32, error) {
	sql := `
		UPDATE
			` + module.DBSchema + `.catalog_categories
		SET
			tag_version = tag_version + 1,
			datetime_updated = NOW(),
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

func (re *Category) Delete(ctx context.Context, id mrtype.KeyInt32) error {
	sql := `
		UPDATE
			` + module.DBSchema + `.catalog_categories
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
