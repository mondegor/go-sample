package repository

import (
	"context"
	"strings"

	"github.com/google/uuid"
	"github.com/mondegor/go-storage/mrpostgres/db"
	"github.com/mondegor/go-storage/mrsql"
	"github.com/mondegor/go-storage/mrstorage"
	"github.com/mondegor/go-webcore/mrenum"
	"github.com/mondegor/go-webcore/mrtype"

	"github.com/mondegor/go-sample/internal/catalog/product/module"
	"github.com/mondegor/go-sample/internal/catalog/product/section/adm/entity"
)

type (
	// ProductPostgres - comment struct.
	ProductPostgres struct {
		client          mrstorage.DBConnManager
		sqlBuilder      mrstorage.SQLBuilder
		repoIDByArticle db.FieldFetcher[string, uint64]
		repoStatus      db.FieldWithVersionUpdater[uint64, uint32, mrenum.ItemStatus]
		repoSoftDeleter db.RowSoftDeleter[uint64]
		repoTotalRows   db.TotalRowsFetcher[uint64]
	}
)

// NewProductPostgres - создаёт объект ProductPostgres.
func NewProductPostgres(client mrstorage.DBConnManager, sqlBuilder mrstorage.SQLBuilder) *ProductPostgres {
	return &ProductPostgres{
		client:     client,
		sqlBuilder: sqlBuilder,
		repoIDByArticle: db.NewFieldFetcher[string, uint64](
			client,
			module.DBTableNameProducts,
			"product_article",
			"product_id",
			module.DBFieldDeletedAt,
		),
		repoStatus: db.NewFieldWithVersionUpdater[uint64, uint32, mrenum.ItemStatus](
			client,
			module.DBTableNameProducts,
			"product_id",
			module.DBFieldTagVersion,
			"product_status",
			module.DBFieldDeletedAt,
		),
		repoSoftDeleter: db.NewRowSoftDeleter[uint64](
			client,
			module.DBTableNameProducts,
			"product_id",
			module.DBFieldTagVersion,
			module.DBFieldDeletedAt,
		),
		repoTotalRows: db.NewTotalRowsFetcher[uint64](
			client,
			module.DBTableNameProducts,
		),
	}
}

// NewCondition - comment method.
func (re *ProductPostgres) NewCondition(categoryID uuid.UUID) mrstorage.SQLPartFunc {
	return re.sqlBuilder.Condition().HelpFunc(
		func(c mrstorage.SQLConditionHelper) mrstorage.SQLPartFunc {
			return c.JoinAnd(
				c.Equal("category_id", categoryID),
				c.Expr("deleted_at IS NULL"),
			)
		},
	)
}

// FetchWithTotal - comment method.
func (re *ProductPostgres) FetchWithTotal(ctx context.Context, params entity.ProductParams) (rows []entity.Product, countRows uint64, err error) {
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
func (re *ProductPostgres) fetch(
	ctx context.Context,
	condition mrstorage.SQLPart,
	orderBy mrstorage.SQLPart,
	limit mrstorage.SQLPart,
	maxRows uint64,
) ([]entity.Product, error) {
	whereStr, whereArgs := condition.ToSQL()

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
			` + module.DBTableNameProducts + `
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

	rows := make([]entity.Product, 0, maxRows)

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

// FetchOne - comment method.
func (re *ProductPostgres) FetchOne(ctx context.Context, rowID uint64) (entity.Product, error) {
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
			` + module.DBTableNameProducts + `
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
func (re *ProductPostgres) FetchIDByArticle(ctx context.Context, article string) (rowID uint64, err error) {
	return re.repoIDByArticle.Fetch(ctx, article)
}

// FetchStatus - comment method.
// result: mrenum.ItemStatus - exists, ErrStorageNoRowFound - not exists, error - query error.
func (re *ProductPostgres) FetchStatus(ctx context.Context, rowID uint64) (mrenum.ItemStatus, error) {
	return re.repoStatus.Fetch(ctx, rowID)
}

// Insert - comment method.
func (re *ProductPostgres) Insert(ctx context.Context, row entity.Product) (rowID uint64, err error) {
	sql := `
		INSERT INTO ` + module.DBTableNameProducts + `
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

	err = re.client.Conn(ctx).QueryRow(
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
func (re *ProductPostgres) Update(ctx context.Context, row entity.Product) (tagVersion uint32, err error) {
	set, err := re.sqlBuilder.Set().BuildEntity(row)
	if err != nil || set.Empty() {
		return 0, err
	}

	args := []any{
		row.ID,
		row.TagVersion,
	}

	setStr, setArgs := set.WithStartArg(len(args) + 1).ToSQL()

	sql := `
		UPDATE
			` + module.DBTableNameProducts + `
		SET
			tag_version = tag_version + 1,
			updated_at = NOW(),
			` + setStr + `
		WHERE
			product_id = $1 AND tag_version = $2 AND deleted_at IS NULL
		RETURNING
			tag_version;`

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
func (re *ProductPostgres) UpdateStatus(ctx context.Context, row entity.Product) (tagVersion uint32, err error) {
	return re.repoStatus.Update(ctx, row.ID, row.TagVersion, row.Status)
}

// Delete - comment method.
func (re *ProductPostgres) Delete(ctx context.Context, rowID uint64) error {
	return re.repoSoftDeleter.Delete(ctx, rowID)
}

func (re *ProductPostgres) fetchCondition(filter entity.ProductListFilter) mrstorage.SQLPartFunc {
	return re.sqlBuilder.Condition().HelpFunc(
		func(c mrstorage.SQLConditionHelper) mrstorage.SQLPartFunc {
			return c.JoinAnd(
				c.Expr("deleted_at IS NULL"),
				c.FilterEqual("category_id", filter.CategoryID),
				c.FilterLikeFields([]string{"UPPER(product_article)", "UPPER(product_caption)"}, strings.ToUpper(filter.SearchText)),
				c.FilterAnyOf("trademark_id", filter.TrademarkIDs),
				c.FilterRangeInt64("product_price", filter.Price, 0),
				c.FilterAnyOf("product_status", filter.Statuses),
			)
		},
	)
}

func (re *ProductPostgres) fetchOrderBy(sorter mrtype.SortParams) mrstorage.SQLPartFunc {
	return re.sqlBuilder.OrderBy().HelpFunc(
		func(o mrstorage.SQLOrderByHelper) mrstorage.SQLPartFunc {
			return o.JoinComma(
				o.Field(sorter.FieldName, sorter.Direction),
				o.Field("product_id", mrenum.SortDirectionASC),
			)
		},
	)
}
