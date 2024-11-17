package repository

import (
	"context"
	"strings"

	"github.com/google/uuid"
	"github.com/mondegor/go-storage/mrpostgres/db"
	"github.com/mondegor/go-storage/mrstorage"
	"github.com/mondegor/go-webcore/mrenum"

	"github.com/mondegor/go-sample/internal/catalog/category/module"
	"github.com/mondegor/go-sample/internal/catalog/category/section/pub/entity"
)

type (
	// CategoryPostgres - comment struct.
	CategoryPostgres struct {
		client        mrstorage.DBConnManager
		sqlBuilder    mrstorage.SQLBuilder
		repoTotalRows db.TotalRowsFetcher[uint64]
	}
)

// NewCategoryPostgres - создаёт объект CategoryPostgres.
func NewCategoryPostgres(client mrstorage.DBConnManager, sqlBuilder mrstorage.SQLBuilder) *CategoryPostgres {
	return &CategoryPostgres{
		client:     client,
		sqlBuilder: sqlBuilder,
		repoTotalRows: db.NewTotalRowsFetcher[uint64](
			client,
			module.DBTableNameCategories,
		),
	}
}

// FetchWithTotal - comment method.
func (re *CategoryPostgres) FetchWithTotal(ctx context.Context, params entity.CategoryParams) (rows []entity.Category, countRows uint64, err error) {
	condition := re.sqlBuilder.Condition().Build(re.fetchCondition(params.Filter))

	total, err := re.repoTotalRows.Fetch(ctx, condition)
	if err != nil || total == 0 {
		return nil, 0, err
	}

	rows, err = re.fetch(ctx, condition, total)
	if err != nil {
		return nil, 0, err
	}

	return rows, total, nil
}

// Fetch - comment method.
func (re *CategoryPostgres) fetch(
	ctx context.Context,
	condition mrstorage.SQLPart,
	maxRows uint64,
) ([]entity.Category, error) {
	whereStr, whereArgs := condition.ToSQL()

	sql := `
		SELECT
			category_id,
			category_caption,
			COALESCE(image_meta ->> 'path', '') as image_url
		FROM
			` + module.DBTableNameCategories + `
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

	rows := make([]entity.Category, 0, maxRows)

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

// FetchOne - comment method.
func (re *CategoryPostgres) FetchOne(ctx context.Context, rowID uuid.UUID) (entity.Category, error) {
	sql := `
		SELECT
			category_caption,
			COALESCE(image_meta ->> 'path', '') as image_url
		FROM
			` + module.DBTableNameCategories + `
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

func (re *CategoryPostgres) fetchCondition(filter entity.CategoryListFilter) mrstorage.SQLPartFunc {
	return re.sqlBuilder.Condition().HelpFunc(
		func(c mrstorage.SQLConditionHelper) mrstorage.SQLPartFunc {
			return c.JoinAnd(
				c.Expr("deleted_at IS NULL"),
				c.Equal("category_status", mrenum.ItemStatusEnabled),
				c.FilterLike("UPPER(category_caption)", strings.ToUpper(filter.SearchText)),
			)
		},
	)
}
