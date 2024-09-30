package repository

import (
	"context"
	"strings"

	"github.com/google/uuid"
	"github.com/mondegor/go-storage/mrstorage"
	"github.com/mondegor/go-webcore/mrenum"

	"github.com/mondegor/go-sample/internal/catalog/category/module"
	"github.com/mondegor/go-sample/internal/catalog/category/section/pub/entity"
)

type (
	// CategoryPostgres - comment struct.
	CategoryPostgres struct {
		client    mrstorage.DBConnManager
		sqlSelect mrstorage.SQLBuilderSelect
	}
)

// NewCategoryPostgres - создаёт объект CategoryPostgres.
func NewCategoryPostgres(client mrstorage.DBConnManager, sqlSelect mrstorage.SQLBuilderSelect) *CategoryPostgres {
	return &CategoryPostgres{
		client:    client,
		sqlSelect: sqlSelect,
	}
}

// NewSelectParams - comment method.
func (re *CategoryPostgres) NewSelectParams(params entity.CategoryParams) mrstorage.SQLSelectParams {
	return mrstorage.SQLSelectParams{
		Where: re.sqlSelect.Where(func(w mrstorage.SQLBuilderWhere) mrstorage.SQLBuilderPartFunc {
			return w.JoinAnd(
				w.Expr("deleted_at IS NULL"),
				w.Equal("category_status", mrenum.ItemStatusEnabled),
				w.FilterLike("UPPER(category_caption)", strings.ToUpper(params.Filter.SearchText)),
			)
		}),
	}
}

// Fetch - comment method.
func (re *CategoryPostgres) Fetch(ctx context.Context, params mrstorage.SQLSelectParams) ([]entity.Category, error) {
	whereStr, whereArgs := params.Where.ToSQL()

	sql := `
		SELECT
			category_id,
			category_caption,
			COALESCE(image_meta ->> 'path', '') as image_url
		FROM
			` + module.DBSchema + `.` + module.DBTableNameCategories + `
		WHERE
			` + whereStr + `
		ORDER BY
			category_caption ASC, category_id ASC;`

	cursor, err := re.client.Conn(ctx).Query(
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

// FetchTotal - comment method.
func (re *CategoryPostgres) FetchTotal(ctx context.Context, where mrstorage.SQLBuilderPart) (int64, error) {
	whereStr, whereArgs := where.ToSQL()

	sql := `
		SELECT
			COUNT(*)
		FROM
			` + module.DBSchema + `.` + module.DBTableNameCategories + `
		WHERE
			` + whereStr + `;`

	var totalRow int64

	err := re.client.Conn(ctx).QueryRow(
		ctx,
		sql,
		whereArgs...,
	).Scan(
		&totalRow,
	)

	return totalRow, err
}

// FetchOne - comment method.
func (re *CategoryPostgres) FetchOne(ctx context.Context, rowID uuid.UUID) (entity.Category, error) {
	sql := `
		SELECT
			category_caption,
			COALESCE(image_meta ->> 'path', '') as image_url
		FROM
			` + module.DBSchema + `.` + module.DBTableNameCategories + `
		WHERE
			category_id = $1 AND category_status = $2 AND deleted_at IS NULL
		LIMIT 1;`

	row := entity.Category{ID: rowID}

	err := re.client.Conn(ctx).QueryRow(
		ctx,
		sql,
		rowID,
		mrenum.ItemStatusEnabled,
	).Scan(
		&row.Caption,
		&row.ImageURL,
	)

	return row, err
}
