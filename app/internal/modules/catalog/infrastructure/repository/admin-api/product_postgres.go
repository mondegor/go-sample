package repository

import (
	"context"
	module "go-sample/internal/modules/catalog"
	"go-sample/internal/modules/catalog/entity/admin-api"
	"strings"

	"github.com/mondegor/go-components/mrorderer"
	"github.com/mondegor/go-storage/mrsql"
	"github.com/mondegor/go-storage/mrstorage"
	"github.com/mondegor/go-webcore/mrenum"
	"github.com/mondegor/go-webcore/mrtype"
)

type (
	Product struct {
		client    mrstorage.DBConn
		sqlSelect mrstorage.SqlBuilderSelect
		sqlUpdate mrstorage.SqlBuilderUpdate
	}
)

func NewProduct(
	client mrstorage.DBConn,
	sqlSelect mrstorage.SqlBuilderSelect,
	sqlUpdate mrstorage.SqlBuilderUpdate,
) *Product {
	return &Product{
		client:    client,
		sqlSelect: sqlSelect,
		sqlUpdate: sqlUpdate,
	}
}

func (re *Product) GetMetaData(categoryID mrtype.KeyInt32) mrorderer.EntityMeta {
	return mrorderer.NewEntityMeta(
		module.DBSchemaProduct+".products",
		"product_id",
		re.sqlSelect.Where(func(w mrstorage.SqlBuilderWhere) mrstorage.SqlBuilderPartFunc {
			return w.JoinAnd(
				w.Equal("category_id", categoryID),
				w.NotEqual("product_status", mrenum.ItemStatusRemoved),
			)
		}),
	)
}

func (re *Product) NewFetchParams(params entity.ProductParams) mrstorage.SqlSelectParams {
	return mrstorage.SqlSelectParams{
		Where: re.sqlSelect.Where(func(w mrstorage.SqlBuilderWhere) mrstorage.SqlBuilderPartFunc {
			return w.JoinAnd(
				w.NotEqual("product_status", mrenum.ItemStatusRemoved),
				w.FilterEqualInt64("category_id", int64(params.Filter.CategoryID), 0),
				w.FilterAnyOf("trademark_id", params.Filter.Trademarks),
				w.FilterLikeFields([]string{"UPPER(product_article)", "UPPER(product_caption)"}, strings.ToUpper(params.Filter.SearchText)),
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

func (re *Product) Fetch(ctx context.Context, params mrstorage.SqlSelectParams) ([]entity.Product, error) {
	whereStr, whereArgs := params.Where.ToSql()

	sql := `
		SELECT
			product_id,
			tag_version,
			datetime_created as createdAt,
			datetime_updated as updatedAt,
			category_id,
			trademark_id,
			product_article as article,
			product_caption as caption,
			product_price as price,
			product_status
		FROM
			` + module.DBSchemaProduct + `.products
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
			&row.CreatedAt,
			&row.UpdatedAt,
			&row.CategoryID,
			&row.TrademarkID,
			&row.Article,
			&row.Caption,
			&row.Price,
			&row.Status,
		)

		if err != nil {
			return nil, err
		}

		rows = append(rows, row)
	}

	return rows, cursor.Err()
}

func (re *Product) FetchTotal(ctx context.Context, where mrstorage.SqlBuilderPart) (int64, error) {
	whereStr, whereArgs := where.ToSql()

	sql := `
		SELECT
			COUNT(*)
		FROM
			` + module.DBSchemaProduct + `.products
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

func (re *Product) LoadOne(ctx context.Context, row *entity.Product) error {
	sql := `
		SELECT
			tag_version,
			datetime_created,
			datetime_updated,
			category_id,
			trademark_id,
			product_article,
			product_caption,
			product_price,
			product_status
		FROM
			` + module.DBSchemaProduct + `.products
		WHERE
			product_id = $1 AND product_status <> $2
		LIMIT 1;`

	return re.client.QueryRow(
		ctx,
		sql,
		row.ID,
		mrenum.ItemStatusRemoved,
	).Scan(
		&row.TagVersion,
		&row.CreatedAt,
		&row.UpdatedAt,
		&row.CategoryID,
		&row.TrademarkID,
		&row.Article,
		&row.Caption,
		&row.Price,
		&row.Status,
	)
}

func (re *Product) FetchIdByArticle(ctx context.Context, article string) (mrtype.KeyInt32, error) {
	sql := `
		SELECT
			product_id
		FROM
			` + module.DBSchemaProduct + `.products
		WHERE
			product_article = $1
		LIMIT 1;`

	var id mrtype.KeyInt32

	err := re.client.QueryRow(
		ctx,
		sql,
		article,
	).Scan(
		&id,
	)

	return id, err
}

func (re *Product) FetchStatus(ctx context.Context, row *entity.Product) (mrenum.ItemStatus, error) {
	sql := `
		SELECT
			product_status
		FROM
			` + module.DBSchemaProduct + `.products
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
func (re *Product) IsExists(ctx context.Context, id mrtype.KeyInt32) error {
	sql := `
		SELECT
			1
		FROM
			` + module.DBSchemaProduct + `.products
		WHERE
			product_id = $1 AND product_status <> $2
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

func (re *Product) Insert(ctx context.Context, row *entity.Product) error {
	sql := `
		INSERT INTO ` + module.DBSchemaProduct + `.products
			(
				category_id,
				trademark_id,
				product_article,
				product_caption,
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
		row.TrademarkID,
		row.Article,
		row.Caption,
		row.Price,
		row.Status,
	).Scan(
		&row.ID,
	)

	return err
}

func (re *Product) Update(ctx context.Context, row *entity.Product) (int32, error) {
	set, err := re.sqlUpdate.SetFromEntity(row)

	if err != nil {
		return 0, err
	}

	if set.Empty() {
		return 0, nil
	}

	args := []any{
		row.ID,
		row.TagVersion,
		mrenum.ItemStatusRemoved,
	}

	setStr, setArgs := set.Param(len(args) + 1).ToSql()

	sql := `
		UPDATE
			` + module.DBSchemaProduct + `.products
		SET
			tag_version = tag_version + 1,
			datetime_updated = NOW(),
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

func (re *Product) UpdateStatus(ctx context.Context, row *entity.Product) (int32, error) {
	sql := `
		UPDATE
			` + module.DBSchemaProduct + `.products
		SET
			tag_version = tag_version + 1,
			datetime_updated = NOW(),
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

func (re *Product) Delete(ctx context.Context, id mrtype.KeyInt32) error {
	sql := `
		UPDATE
			` + module.DBSchemaProduct + `.products
		SET
			tag_version = tag_version + 1,
			datetime_updated = NOW(),
			product_article = NULL,
			prev_field_id = NULL,
			next_field_id = NULL,
			order_field = NULL,
			product_status = $2
		WHERE
			product_id = $1 AND product_status <> $2;`

	return re.client.Exec(
		ctx,
		sql,
		id,
		mrenum.ItemStatusRemoved,
	)
}
