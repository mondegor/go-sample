package repository

import (
	"context"
	"strings"

	"github.com/mondegor/go-storage/mrpostgres/db"
	"github.com/mondegor/go-storage/mrstorage"
	"github.com/mondegor/go-webcore/mrenum"
	"github.com/mondegor/go-webcore/mrtype"

	"github.com/mondegor/go-sample/internal/catalog/trademark/module"
	"github.com/mondegor/go-sample/internal/catalog/trademark/section/adm/entity"
)

type (
	// TrademarkPostgres - comment struct.
	TrademarkPostgres struct {
		client          mrstorage.DBConnManager
		sqlBuilder      mrstorage.SQLBuilder
		repoStatus      db.FieldWithVersionUpdater[uint64, uint32, mrenum.ItemStatus]
		repoSoftDeleter db.RowSoftDeleter[uint64]
		repoTotalRows   db.TotalRowsFetcher[uint64]
	}
)

// NewTrademarkPostgres - создаёт объект TrademarkPostgres.
func NewTrademarkPostgres(client mrstorage.DBConnManager, sqlBuilder mrstorage.SQLBuilder) *TrademarkPostgres {
	return &TrademarkPostgres{
		client:     client,
		sqlBuilder: sqlBuilder,
		repoStatus: db.NewFieldWithVersionUpdater[uint64, uint32, mrenum.ItemStatus](
			client,
			module.DBTableNameTrademarks,
			"trademark_id",
			module.DBFieldTagVersion,
			"trademark_status",
			module.DBFieldDeletedAt,
		),
		repoSoftDeleter: db.NewRowSoftDeleter[uint64](
			client,
			module.DBTableNameTrademarks,
			"trademark_id",
			module.DBFieldTagVersion,
			module.DBFieldDeletedAt,
		),
		repoTotalRows: db.NewTotalRowsFetcher[uint64](
			client,
			module.DBTableNameTrademarks,
		),
	}
}

// FetchWithTotal - comment method.
func (re *TrademarkPostgres) FetchWithTotal(ctx context.Context, params entity.TrademarkParams) (rows []entity.Trademark, countRows uint64, err error) {
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
func (re *TrademarkPostgres) fetch(
	ctx context.Context,
	condition mrstorage.SQLPart,
	orderBy mrstorage.SQLPart,
	limit mrstorage.SQLPart,
	maxRows uint64,
) ([]entity.Trademark, error) {
	whereStr, whereArgs := condition.ToSQL()

	sql := `
		SELECT
			trademark_id,
			tag_version,
			trademark_caption as caption,
			trademark_status,
			created_at as createdAt,
			updated_at as updatedAt
		FROM
			` + module.DBTableNameTrademarks + `
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

	rows := make([]entity.Trademark, 0, maxRows)

	for cursor.Next() {
		var row entity.Trademark

		err = cursor.Scan(
			&row.ID,
			&row.TagVersion,
			&row.Caption,
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
func (re *TrademarkPostgres) FetchOne(ctx context.Context, rowID uint64) (entity.Trademark, error) {
	sql := `
		SELECT
			tag_version,
			trademark_caption,
			trademark_status,
			created_at,
			updated_at
		FROM
			` + module.DBTableNameTrademarks + `
		WHERE
			trademark_id = $1 AND deleted_at IS NULL
		LIMIT 1;`

	row := entity.Trademark{ID: rowID}

	err := re.client.Conn(ctx).QueryRow(
		ctx,
		sql,
		rowID,
	).Scan(
		&row.TagVersion,
		&row.Caption,
		&row.Status,
		&row.CreatedAt,
		&row.UpdatedAt,
	)

	return row, err
}

// FetchStatus - comment method.
// result: mrenum.ItemStatus - exists, ErrStorageNoRowFound - not exists, error - query error.
func (re *TrademarkPostgres) FetchStatus(ctx context.Context, rowID uint64) (mrenum.ItemStatus, error) {
	return re.repoStatus.Fetch(ctx, rowID)
}

// Insert - comment method.
func (re *TrademarkPostgres) Insert(ctx context.Context, row entity.Trademark) (rowID uint64, err error) {
	sql := `
		INSERT INTO ` + module.DBTableNameTrademarks + `
			(
				trademark_caption,
				trademark_status
			)
		VALUES
			($1, $2)
		RETURNING
			trademark_id;`

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
func (re *TrademarkPostgres) Update(ctx context.Context, row entity.Trademark) (tagVersion uint32, err error) {
	sql := `
		UPDATE
			` + module.DBTableNameTrademarks + `
		SET
			tag_version = tag_version + 1,
			updated_at = NOW(),
			trademark_caption = $3
		WHERE
			trademark_id = $1 AND tag_version = $2 AND deleted_at IS NULL
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
func (re *TrademarkPostgres) UpdateStatus(ctx context.Context, row entity.Trademark) (tagVersion uint32, err error) {
	return re.repoStatus.Update(ctx, row.ID, row.TagVersion, row.Status)
}

// Delete - comment method.
func (re *TrademarkPostgres) Delete(ctx context.Context, rowID uint64) error {
	return re.repoSoftDeleter.Delete(ctx, rowID)
}

func (re *TrademarkPostgres) fetchCondition(filter entity.TrademarkListFilter) mrstorage.SQLPartFunc {
	return re.sqlBuilder.Condition().HelpFunc(
		func(c mrstorage.SQLConditionHelper) mrstorage.SQLPartFunc {
			return c.JoinAnd(
				c.Expr("deleted_at IS NULL"),
				c.FilterLike("UPPER(trademark_caption)", strings.ToUpper(filter.SearchText)),
				c.FilterAnyOf("trademark_status", filter.Statuses),
			)
		},
	)
}

func (re *TrademarkPostgres) fetchOrderBy(sorter mrtype.SortParams) mrstorage.SQLPartFunc {
	return re.sqlBuilder.OrderBy().HelpFunc(
		func(o mrstorage.SQLOrderByHelper) mrstorage.SQLPartFunc {
			return o.JoinComma(
				o.Field(sorter.FieldName, sorter.Direction),
				o.Field("trademark_id", mrenum.SortDirectionASC),
			)
		},
	)
}
