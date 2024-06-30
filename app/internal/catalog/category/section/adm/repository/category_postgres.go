package repository

import (
	"context"
	"strings"

	"github.com/mondegor/go-sample/internal/catalog/category/module"
	"github.com/mondegor/go-sample/internal/catalog/category/section/adm/entity"
	"github.com/mondegor/go-sample/internal/catalog/category/shared/repository"

	"github.com/google/uuid"
	"github.com/mondegor/go-storage/mrstorage"
	"github.com/mondegor/go-webcore/mrenum"
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
				w.FilterLike("UPPER(category_caption)", strings.ToUpper(params.Filter.SearchText)),
				w.FilterAnyOf("category_status", params.Filter.Statuses),
			)
		}),
		OrderBy: re.sqlSelect.OrderBy(func(s mrstorage.SQLBuilderOrderBy) mrstorage.SQLBuilderPartFunc {
			return s.Join(
				s.Field(params.Sorter.FieldName, params.Sorter.Direction),
				s.Field("category_id", mrenum.SortDirectionASC),
			)
		}),
		Limit: re.sqlSelect.Limit(func(p mrstorage.SQLBuilderLimit) mrstorage.SQLBuilderPartFunc {
			return p.OffsetLimit(params.Pager.Index, params.Pager.Size)
		}),
	}
}

// Fetch - comment method.
func (re *CategoryPostgres) Fetch(ctx context.Context, params mrstorage.SQLSelectParams) ([]entity.Category, error) {
	whereStr, whereArgs := params.Where.ToSQL()

	sql := `
		SELECT
			category_id,
			tag_version,
			category_caption as caption,
			image_meta,
			category_status,
			created_at as createdAt,
			updated_at as updatedAt
		FROM
			` + module.DBSchema + `.` + module.DBTableNameCategories + `
		WHERE
			` + whereStr + `
		ORDER BY
			` + params.OrderBy.String() + params.Limit.String() + `;`

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
			&row.TagVersion,
			&row.Caption,
			&row.ImageMeta,
			&row.Status,
			&row.CreatedAt,
			&row.UpdatedAt,
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
			tag_version,
			category_caption,
			image_meta,
			category_status,
			created_at,
			updated_at
		FROM
			` + module.DBSchema + `.` + module.DBTableNameCategories + `
		WHERE
			category_id = $1 AND deleted_at IS NULL
		LIMIT 1;`

	row := entity.Category{ID: rowID}

	err := re.client.Conn(ctx).QueryRow(
		ctx,
		sql,
		rowID,
	).Scan(
		&row.TagVersion,
		&row.Caption,
		&row.ImageMeta,
		&row.Status,
		&row.CreatedAt,
		&row.UpdatedAt,
	)

	return row, err
}

// FetchStatus - comment method.
// result: mrenum.ItemStatus - exists, ErrStorageNoRowFound - not exists, error - query error.
func (re *CategoryPostgres) FetchStatus(ctx context.Context, rowID uuid.UUID) (mrenum.ItemStatus, error) {
	return repository.CategoryFetchStatusPostgres(ctx, re.client, rowID)
}

// Insert - comment method.
func (re *CategoryPostgres) Insert(ctx context.Context, row entity.Category) (uuid.UUID, error) {
	sql := `
		INSERT INTO ` + module.DBSchema + `.` + module.DBTableNameCategories + `
			(
				category_id,
				category_caption,
				category_status
			)
		VALUES
			(gen_random_uuid(), $1, $2)
		RETURNING
			category_id;`

	err := re.client.Conn(ctx).QueryRow(
		ctx,
		sql,
		row.Caption,
		row.Status,
	).Scan(
		&row.ID,
	)

	return row.ID, err
}

// Update - comment method.
func (re *CategoryPostgres) Update(ctx context.Context, row entity.Category) (int32, error) {
	sql := `
		UPDATE
			` + module.DBSchema + `.` + module.DBTableNameCategories + `
		SET
			tag_version = tag_version + 1,
			updated_at = NOW(),
			category_caption = $3
		WHERE
			category_id = $1 AND tag_version = $2 AND deleted_at IS NULL
		RETURNING
			tag_version;`

	var tagVersion int32

	err := re.client.Conn(ctx).QueryRow(
		ctx,
		sql,
		row.ID,
		row.TagVersion,
		row.Caption,
	).Scan(
		&tagVersion,
	)

	return tagVersion, err
}

// UpdateStatus - comment method.
func (re *CategoryPostgres) UpdateStatus(ctx context.Context, row entity.Category) (int32, error) {
	sql := `
		UPDATE
			` + module.DBSchema + `.` + module.DBTableNameCategories + `
		SET
			tag_version = tag_version + 1,
			updated_at = NOW(),
			category_status = $3
		WHERE
			category_id = $1 AND tag_version = $2 AND deleted_at IS NULL
		RETURNING
			tag_version;`

	var tagVersion int32

	err := re.client.Conn(ctx).QueryRow(
		ctx,
		sql,
		row.ID,
		row.TagVersion,
		row.Status,
	).Scan(
		&tagVersion,
	)

	return tagVersion, err
}

// Delete - comment method.
func (re *CategoryPostgres) Delete(ctx context.Context, rowID uuid.UUID) error {
	sql := `
		UPDATE
			` + module.DBSchema + `.` + module.DBTableNameCategories + `
		SET
			tag_version = tag_version + 1,
			deleted_at = NOW()
		WHERE
			category_id = $1 AND deleted_at IS NULL;`

	return re.client.Conn(ctx).Exec(
		ctx,
		sql,
		rowID,
	)
}
