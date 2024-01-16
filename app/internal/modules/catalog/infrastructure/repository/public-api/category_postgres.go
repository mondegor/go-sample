package repository

import (
	"context"
	module "go-sample/internal/modules/catalog"
	"go-sample/internal/modules/catalog/entity/public-api"
	"strings"

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
				w.Equal("category_status", mrenum.ItemStatusEnabled),
				w.FilterLike("UPPER(category_caption)", strings.ToUpper(params.Filter.SearchText)),
			)
		}),
	}
}

func (re *CategoryPostgres) Fetch(ctx context.Context, params mrstorage.SqlSelectParams) ([]entity.Category, error) {
	whereStr, whereArgs := params.Where.ToSql()

	sql := `
		SELECT
			category_id,
			category_caption,
			COALESCE(image_meta ->> 'path', '') as image_url
		FROM
			` + module.UnitCategoryDBSchema + `.categories
		WHERE
			` + whereStr + `
		ORDER BY
			category_caption ASC, category_id ASC;`

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
			&row.Caption,
			&row.ImageURL,
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
			` + module.UnitCategoryDBSchema + `.categories
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

func (re *CategoryPostgres) LoadOne(ctx context.Context, row *entity.Category) error {
	sql := `
		SELECT
			category_caption,
			COALESCE(image_meta ->> 'path', '') as image_url
		FROM
			` + module.UnitCategoryDBSchema + `.categories
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
		&row.ImageURL,
	)
}
