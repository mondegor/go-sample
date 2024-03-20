package repository

import (
	"context"
	module "go-sample/internal/modules/catalog/product"
	"go-sample/internal/modules/catalog/product/entity/admin-api"
	"strings"

	"github.com/google/uuid"
	"github.com/mondegor/go-components/mrorderer"
	"github.com/mondegor/go-storage/mrsql"
	"github.com/mondegor/go-storage/mrstorage"
	"github.com/mondegor/go-webcore/mrenum"
	"github.com/mondegor/go-webcore/mrtype"
)

type (
	ProductPostgres struct {
		client    mrstorage.DBConn
		sqlSelect mrstorage.SqlBuilderSelect
		sqlUpdate mrstorage.SqlBuilderUpdate
	}
)

func NewProductPostgres(
	client mrstorage.DBConn,
	sqlSelect mrstorage.SqlBuilderSelect,
	sqlUpdate mrstorage.SqlBuilderUpdate,
) *ProductPostgres {
	return &ProductPostgres{
		client:    client,
		sqlSelect: sqlSelect,
		sqlUpdate: sqlUpdate,
	}
}

func (re *ProductPostgres) NewOrderMeta(categoryID uuid.UUID) mrorderer.EntityMeta {
	return mrorderer.NewEntityMeta(
		module.DBSchema+".products",
		"product_id",
		re.sqlSelect.Where(func(w mrstorage.SqlBuilderWhere) mrstorage.SqlBuilderPartFunc {
			return w.JoinAnd(
				w.Equal("category_id", categoryID),
				w.NotEqual("product_status", mrenum.ItemStatusRemoved),
			)
		}),
	)
}

func (re *ProductPostgres) NewSelectParams(params entity.ProductParams) mrstorage.SqlSelectParams {
	return mrstorage.SqlSelectParams{
		Where: re.sqlSelect.Where(func(w mrstorage.SqlBuilderWhere) mrstorage.SqlBuilderPartFunc {
			return w.JoinAnd(
				w.NotEqual("product_status", mrenum.ItemStatusRemoved),
				w.FilterEqualUUID("category_id", params.Filter.CategoryID),
				w.FilterLikeFields([]string{"UPPER(product_article)", "UPPER(product_caption)"}, strings.ToUpper(params.Filter.SearchText)),
				w.FilterAnyOf("trademark_id", params.Filter.TrademarkIDs),
				w.FilterRangeInt64("product_price", params.Filter.Price, 0),
				w.FilterAnyOf("product_status", params.Filter.Statuses),
			)
		}),
		OrderBy: re.sqlSelect.OrderBy(func(s mrstorage.SqlBuilderOrderBy) mrstorage.SqlBuilderPartFunc {
			return s.Join(
				s.Field(params.Sorter.FieldName, params.Sorter.Direction),
				s.Field("product_id", mrenum.SortDirectionASC),
			)
		}),
		Pager: re.sqlSelect.Pager(func(p mrstorage.SqlBuilderPager) mrstorage.SqlBuilderPartFunc {
			return p.OffsetLimit(params.Pager.Index, params.Pager.Size)
		}),
	}
}

func (re *ProductPostgres) Fetch(ctx context.Context, params mrstorage.SqlSelectParams) ([]entity.Product, error) {
	whereStr, whereArgs := params.Where.ToSql()

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
			` + module.DBSchema + `.products
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

func (re *ProductPostgres) FetchTotal(ctx context.Context, where mrstorage.SqlBuilderPart) (int64, error) {
	whereStr, whereArgs := where.ToSql()

	sql := `
		SELECT
			COUNT(*)
		FROM
			` + module.DBSchema + `.products
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
			` + module.DBSchema + `.products
		WHERE
			product_id = $1 AND product_status <> $2
		LIMIT 1;`

	row := entity.Product{ID: rowID}

	err := re.client.QueryRow(
		ctx,
		sql,
		rowID,
		mrenum.ItemStatusRemoved,
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

func (re *ProductPostgres) FetchIdByArticle(ctx context.Context, article string) (mrtype.KeyInt32, error) {
	sql := `
		SELECT
			product_id
		FROM
			` + module.DBSchema + `.products
		WHERE
			product_article = $1
		LIMIT 1;`

	var rowID mrtype.KeyInt32

	err := re.client.QueryRow(
		ctx,
		sql,
		article,
	).Scan(
		&rowID,
	)

	return rowID, err
}

func (re *ProductPostgres) FetchStatus(ctx context.Context, row entity.Product) (mrenum.ItemStatus, error) {
	sql := `
		SELECT
			product_status
		FROM
			` + module.DBSchema + `.products
		WHERE
			product_id = $1 AND product_status <> $2
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
func (re *ProductPostgres) IsExists(ctx context.Context, rowID mrtype.KeyInt32) error {
	sql := `
		SELECT
			product_id
		FROM
			` + module.DBSchema + `.products
		WHERE
			product_id = $1 AND product_status <> $2
		LIMIT 1;`

	return re.client.QueryRow(
		ctx,
		sql,
		rowID,
		mrenum.ItemStatusRemoved,
	).Scan(
		&rowID,
	)
}

func (re *ProductPostgres) Insert(ctx context.Context, row entity.Product) (mrtype.KeyInt32, error) {
	sql := `
		INSERT INTO ` + module.DBSchema + `.products
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

	err := re.client.QueryRow(
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

func (re *ProductPostgres) Update(ctx context.Context, row entity.Product) (int32, error) {
	set, err := re.sqlUpdate.SetFromEntity(row)

	if err != nil || set.Empty() {
		return 0, err
	}

	args := []any{
		row.ID,
		row.TagVersion,
		mrenum.ItemStatusRemoved,
	}

	setStr, setArgs := set.Param(len(args) + 1).ToSql()

	sql := `
		UPDATE
			` + module.DBSchema + `.products
		SET
			tag_version = tag_version + 1,
			updated_at = NOW(),
			` + setStr + `
		WHERE
			product_id = $1 AND tag_version = $2 AND product_status <> $3
		RETURNING
			tag_version;`

	var tagVersion int32

	err = re.client.QueryRow(
		ctx,
		sql,
		mrsql.MergeArgs(args, setArgs)...,
	).Scan(
		&tagVersion,
	)

	return tagVersion, err
}

func (re *ProductPostgres) UpdateStatus(ctx context.Context, row entity.Product) (int32, error) {
	sql := `
		UPDATE
			` + module.DBSchema + `.products
		SET
			tag_version = tag_version + 1,
			updated_at = NOW(),
			product_status = $4
		WHERE
			product_id = $1 AND tag_version = $2 AND product_status <> $3
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

func (re *ProductPostgres) Delete(ctx context.Context, rowID mrtype.KeyInt32) error {
	sql := `
		UPDATE
			` + module.DBSchema + `.products
		SET
			tag_version = tag_version + 1,
			updated_at = NOW(),
			product_article = NULL,
			order_index = NULL,
			product_status = $2
		WHERE
			product_id = $1 AND product_status <> $2;`

	return re.client.Exec(
		ctx,
		sql,
		rowID,
		mrenum.ItemStatusRemoved,
	)
}
