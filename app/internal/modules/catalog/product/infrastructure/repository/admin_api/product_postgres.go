package repository

import (
	"context"
	"strings"

	entity "go-sample/internal/modules/catalog/product/entity/admin_api"
	"go-sample/internal/modules/catalog/product/module"

	"github.com/google/uuid"
	"github.com/mondegor/go-storage/mrsql"
	"github.com/mondegor/go-storage/mrstorage"
	"github.com/mondegor/go-webcore/mrenum"
	"github.com/mondegor/go-webcore/mrtype"
)

type (
	// ProductPostgres - comment struct.
	ProductPostgres struct {
		client    mrstorage.DBConnManager
		sqlSelect mrstorage.SQLBuilderSelect
		sqlUpdate mrstorage.SQLBuilderUpdate
	}
)

// NewProductPostgres - comment func.
func NewProductPostgres(client mrstorage.DBConnManager, sqlSelect mrstorage.SQLBuilderSelect, sqlUpdate mrstorage.SQLBuilderUpdate) *ProductPostgres {
	return &ProductPostgres{
		client:    client,
		sqlSelect: sqlSelect,
		sqlUpdate: sqlUpdate,
	}
}

// NewOrderMeta - comment method.
func (re *ProductPostgres) NewOrderMeta(categoryID uuid.UUID) mrstorage.MetaGetter {
	return mrsql.NewEntityMeta(
		module.DBSchema+"."+module.DBTableNameProducts,
		"product_id",
		re.sqlSelect.Where(func(w mrstorage.SQLBuilderWhere) mrstorage.SQLBuilderPartFunc {
			return w.JoinAnd(
				w.Equal("category_id", categoryID),
				w.Expr("deleted_at IS NULL"),
			)
		}),
	)
}

// NewSelectParams - comment method.
func (re *ProductPostgres) NewSelectParams(params entity.ProductParams) mrstorage.SQLSelectParams {
	return mrstorage.SQLSelectParams{
		Where: re.sqlSelect.Where(func(w mrstorage.SQLBuilderWhere) mrstorage.SQLBuilderPartFunc {
			return w.JoinAnd(
				w.Expr("deleted_at IS NULL"),
				w.FilterEqualUUID("category_id", params.Filter.CategoryID),
				w.FilterLikeFields([]string{"UPPER(product_article)", "UPPER(product_caption)"}, strings.ToUpper(params.Filter.SearchText)),
				w.FilterAnyOf("trademark_id", params.Filter.TrademarkIDs),
				w.FilterRangeInt64("product_price", params.Filter.Price, 0),
				w.FilterAnyOf("product_status", params.Filter.Statuses),
			)
		}),
		OrderBy: re.sqlSelect.OrderBy(func(s mrstorage.SQLBuilderOrderBy) mrstorage.SQLBuilderPartFunc {
			return s.Join(
				s.Field(params.Sorter.FieldName, params.Sorter.Direction),
				s.Field("product_id", mrenum.SortDirectionASC),
			)
		}),
		Limit: re.sqlSelect.Limit(func(p mrstorage.SQLBuilderLimit) mrstorage.SQLBuilderPartFunc {
			return p.OffsetLimit(params.Pager.Index, params.Pager.Size)
		}),
	}
}

// Fetch - comment method.
func (re *ProductPostgres) Fetch(ctx context.Context, params mrstorage.SQLSelectParams) ([]entity.Product, error) {
	whereStr, whereArgs := params.Where.ToSQL()

	sql := `
		SELECT
			product_id,
			tag_version,
			category_id,
			product_article as article,
			product_caption as caption,
			trademark_id,
			product_price as price,
			product_status,
			created_at as createdAt,
			updated_at as updatedAt
		FROM
			` + module.DBSchema + `.` + module.DBTableNameProducts + `
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

	rows := make([]entity.Product, 0)

	for cursor.Next() {
		var row entity.Product

		err = cursor.Scan(
			&row.ID,
			&row.TagVersion,
			&row.CategoryID,
			&row.Article,
			&row.Caption,
			&row.TrademarkID,
			&row.Price,
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
func (re *ProductPostgres) FetchTotal(ctx context.Context, where mrstorage.SQLBuilderPart) (int64, error) {
	whereStr, whereArgs := where.ToSQL()

	sql := `
		SELECT
			COUNT(*)
		FROM
			` + module.DBSchema + `.` + module.DBTableNameProducts + `
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
func (re *ProductPostgres) FetchOne(ctx context.Context, rowID mrtype.KeyInt32) (entity.Product, error) {
	sql := `
		SELECT
			tag_version,
			category_id,
			product_article,
			product_caption,
			trademark_id,
			product_price,
			product_status,
			created_at,
			updated_at
		FROM
			` + module.DBSchema + `.` + module.DBTableNameProducts + `
		WHERE
			product_id = $1 AND deleted_at IS NULL
		LIMIT 1;`

	row := entity.Product{ID: rowID}

	err := re.client.Conn(ctx).QueryRow(
		ctx,
		sql,
		rowID,
	).Scan(
		&row.TagVersion,
		&row.CategoryID,
		&row.Article,
		&row.Caption,
		&row.TrademarkID,
		&row.Price,
		&row.Status,
		&row.CreatedAt,
		&row.UpdatedAt,
	)

	return row, err
}

// FetchIDByArticle - comment method.
func (re *ProductPostgres) FetchIDByArticle(ctx context.Context, article string) (mrtype.KeyInt32, error) {
	sql := `
		SELECT
			product_id
		FROM
			` + module.DBSchema + `.` + module.DBTableNameProducts + `
		WHERE
			product_article = $1 AND deleted_at IS NULL
		LIMIT 1;`

	var rowID mrtype.KeyInt32

	err := re.client.Conn(ctx).QueryRow(
		ctx,
		sql,
		article,
	).Scan(
		&rowID,
	)

	return rowID, err
}

// FetchStatus - comment method.
// result: mrenum.ItemStatus - exists, ErrStorageNoRowFound - not exists, error - query error.
func (re *ProductPostgres) FetchStatus(ctx context.Context, rowID mrtype.KeyInt32) (mrenum.ItemStatus, error) {
	sql := `
		SELECT
			product_status
		FROM
			` + module.DBSchema + `.` + module.DBTableNameProducts + `
		WHERE
			product_id = $1 AND deleted_at IS NULL
		LIMIT 1;`

	var status mrenum.ItemStatus

	err := re.client.Conn(ctx).QueryRow(
		ctx,
		sql,
		rowID,
	).Scan(
		&status,
	)

	return status, err
}

// Insert - comment method.
func (re *ProductPostgres) Insert(ctx context.Context, row entity.Product) (mrtype.KeyInt32, error) {
	sql := `
		INSERT INTO ` + module.DBSchema + `.` + module.DBTableNameProducts + `
			(
				category_id,
				product_article,
				product_caption,
				trademark_id,
				product_price,
				product_status
			)
		VALUES
			($1, $2, $3, $4, $5, $6)
		RETURNING
			product_id;`

	err := re.client.Conn(ctx).QueryRow(
		ctx,
		sql,
		row.CategoryID,
		row.Article,
		row.Caption,
		row.TrademarkID,
		row.Price,
		row.Status,
	).Scan(
		&row.ID,
	)

	return row.ID, err
}

// Update - comment method.
func (re *ProductPostgres) Update(ctx context.Context, row entity.Product) (int32, error) {
	set, err := re.sqlUpdate.SetFromEntity(row)

	if err != nil || set.Empty() {
		return 0, err
	}

	args := []any{
		row.ID,
		row.TagVersion,
	}

	setStr, setArgs := set.WithParam(len(args) + 1).ToSQL()

	sql := `
		UPDATE
			` + module.DBSchema + `.` + module.DBTableNameProducts + `
		SET
			tag_version = tag_version + 1,
			updated_at = NOW(),
			` + setStr + `
		WHERE
			product_id = $1 AND tag_version = $2 AND deleted_at IS NULL
		RETURNING
			tag_version;`

	var tagVersion int32

	err = re.client.Conn(ctx).QueryRow(
		ctx,
		sql,
		mrsql.MergeArgs(args, setArgs)...,
	).Scan(
		&tagVersion,
	)

	return tagVersion, err
}

// UpdateStatus - comment method.
func (re *ProductPostgres) UpdateStatus(ctx context.Context, row entity.Product) (int32, error) {
	sql := `
		UPDATE
			` + module.DBSchema + `.` + module.DBTableNameProducts + `
		SET
			tag_version = tag_version + 1,
			updated_at = NOW(),
			product_status = $3
		WHERE
			product_id = $1 AND tag_version = $2 AND deleted_at IS NULL
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
func (re *ProductPostgres) Delete(ctx context.Context, rowID mrtype.KeyInt32) error {
	sql := `
		UPDATE
			` + module.DBSchema + `.` + module.DBTableNameProducts + `
		SET
			tag_version = tag_version + 1,
			deleted_at = NOW()
		WHERE
			product_id = $1 AND deleted_at IS NULL;`

	return re.client.Conn(ctx).Exec(
		ctx,
		sql,
		rowID,
	)
}
