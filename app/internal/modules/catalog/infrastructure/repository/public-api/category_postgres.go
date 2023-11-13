package repository

import (
	"context"
	"go-sample/internal/modules/catalog/entity/public-api"
	"strings"

	"github.com/mondegor/go-storage/mrstorage"
	"github.com/mondegor/go-webcore/mrenum"
)

type (
	Category struct {
		client mrstorage.DBConn
		sqlSelect mrstorage.SqlBuilderSelect
	}
)

func NewCategory(
	client mrstorage.DBConn,
	sqlSelect mrstorage.SqlBuilderSelect,
) *Category {
	return &Category{
		client: client,
		sqlSelect: sqlSelect,
	}
}

func (re *Category) NewFetchParams(params entity.CategoryParams) mrstorage.SqlSelectParams {
	return mrstorage.SqlSelectParams{
		Where: re.sqlSelect.Where(func (w mrstorage.SqlBuilderWhere) mrstorage.SqlBuilderPartFunc {
			return w.JoinAnd(
				w.Equal("category_status", mrenum.ItemStatusEnabled),
				w.FilterLike("UPPER(category_caption)", strings.ToUpper(params.Filter.SearchText)),
			)
		}),
	}
}

func (re *Category) Fetch(ctx context.Context, params mrstorage.SqlSelectParams) ([]entity.Category, error) {
	whereStr, whereArgs := params.Where.ToSql()

	sql := `
		SELECT
			category_id,
			category_caption,
			image_path
		FROM
			public.catalog_categories
		WHERE
			` + whereStr + `
		ORDER BY
			category_caption ASC, category_id ASC;`

	cursor, err := re.client.Query(
		ctx,
		sql,
		whereArgs...
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
			&row.Caption,
			&row.ImagePath,
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
			public.catalog_categories
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
			category_caption,
			image_path
		FROM
			public.catalog_categories
		WHERE
			category_id = $1 AND category_status = $2
		LIMIT 1;`

	return re.client.QueryRow(
		ctx,
		sql,
		row.ID,
		mrenum.ItemStatusEnabled,
	).Scan(
		&row.Caption,
		&row.ImagePath,
	)
}
