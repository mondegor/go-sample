package repository

import (
	"context"
	"strings"

	"github.com/google/uuid"
	"github.com/mondegor/go-storage/mrpostgres/db"
	"github.com/mondegor/go-storage/mrstorage"
	"github.com/mondegor/go-webcore/mrenum"
	"github.com/mondegor/go-webcore/mrtype"

	"github.com/mondegor/go-sample/internal/catalog/category/module"
	"github.com/mondegor/go-sample/internal/catalog/category/section/adm/entity"
)

type (
	// CategoryPostgres - comment struct.
	CategoryPostgres struct {
		client          mrstorage.DBConnManager
		sqlBuilder      mrstorage.SQLBuilder
		repoStatus      db.FieldWithVersionUpdater[uuid.UUID, uint32, mrenum.ItemStatus]
		repoSoftDeleter db.RowSoftDeleter[uuid.UUID]
		repoTotalRows   db.TotalRowsFetcher[uint64]
	}
)

// NewCategoryPostgres - создаёт объект CategoryPostgres.
func NewCategoryPostgres(client mrstorage.DBConnManager, sqlBuilder mrstorage.SQLBuilder) *CategoryPostgres {
	return &CategoryPostgres{
		client:     client,
		sqlBuilder: sqlBuilder,
		repoStatus: db.NewFieldWithVersionUpdater[uuid.UUID, uint32, mrenum.ItemStatus](
			client,
			module.DBTableNameCategories,
			"category_id",
			module.DBFieldTagVersion,
			"category_status",
			module.DBFieldDeletedAt,
		),
		repoSoftDeleter: db.NewRowSoftDeleter[uuid.UUID](
			client,
			module.DBTableNameCategories,
			"category_id",
			module.DBFieldTagVersion,
			module.DBFieldDeletedAt,
		),
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

	if params.Pager.Size > total {
		params.Pager.Size = total
	}

	orderBy := re.sqlBuilder.OrderBy().Build(re.fetchOrderBy(params.Sorter))
	limit := re.sqlBuilder.Limit().Build(params.Pager.Index, params.Pager.Size)

	rows, err = re.fetch(ctx, condition, orderBy, limit, params.Pager.Size)
	if err != nil {
		return nil, 0, err
	}

	return rows, total, nil
}

// Fetch - comment method.
func (re *CategoryPostgres) fetch(
	ctx context.Context,
	condition mrstorage.SQLPart,
	orderBy mrstorage.SQLPart,
	limit mrstorage.SQLPart,
	maxRows uint64,
) ([]entity.Category, error) {
	whereStr, whereArgs := condition.ToSQL()

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
			` + module.DBTableNameCategories + `
		WHERE
			` + whereStr + `
		ORDER BY
			` + orderBy.String() + limit.String() + `;`

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
			` + module.DBTableNameCategories + `
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
	return re.repoStatus.Fetch(ctx, rowID)
}

// Insert - comment method.
func (re *CategoryPostgres) Insert(ctx context.Context, row entity.Category) (rowID uuid.UUID, err error) {
	sql := `
		INSERT INTO ` + module.DBTableNameCategories + `
			(
				category_id,
				category_caption,
				category_status
			)
		VALUES
			(gen_random_uuid(), $1, $2)
		RETURNING
			category_id;`

	err = re.client.Conn(ctx).QueryRow(
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
func (re *CategoryPostgres) Update(ctx context.Context, row entity.Category) (tagVersion uint32, err error) {
	sql := `
		UPDATE
			` + module.DBTableNameCategories + `
		SET
			tag_version = tag_version + 1,
			updated_at = NOW(),
			category_caption = $3
		WHERE
			category_id = $1 AND tag_version = $2 AND deleted_at IS NULL
		RETURNING
			tag_version;`

	err = re.client.Conn(ctx).QueryRow(
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
func (re *CategoryPostgres) UpdateStatus(ctx context.Context, row entity.Category) (tagVersion uint32, err error) {
	return re.repoStatus.Update(ctx, row.ID, row.TagVersion, row.Status)
}

// Delete - comment method.
func (re *CategoryPostgres) Delete(ctx context.Context, rowID uuid.UUID) error {
	return re.repoSoftDeleter.Delete(ctx, rowID)
}

func (re *CategoryPostgres) fetchCondition(filter entity.CategoryListFilter) mrstorage.SQLPartFunc {
	return re.sqlBuilder.Condition().HelpFunc(
		func(c mrstorage.SQLConditionHelper) mrstorage.SQLPartFunc {
			return c.JoinAnd(
				c.Expr("deleted_at IS NULL"),
				c.FilterLike("UPPER(category_caption)", strings.ToUpper(filter.SearchText)),
				c.FilterAnyOf("category_status", filter.Statuses),
			)
		},
	)
}

func (re *CategoryPostgres) fetchOrderBy(sorter mrtype.SortParams) mrstorage.SQLPartFunc {
	return re.sqlBuilder.OrderBy().HelpFunc(
		func(o mrstorage.SQLOrderByHelper) mrstorage.SQLPartFunc {
			return o.JoinComma(
				o.Field(sorter.FieldName, sorter.Direction),
				o.Field("category_id", mrenum.SortDirectionASC),
			)
		},
	)
}
