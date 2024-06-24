package repository

import (
	"context"
	"strings"

	entity "go-sample/internal/modules/catalog/trademark/entity/admin_api"
	repositoryshared "go-sample/internal/modules/catalog/trademark/infrastructure/repository/shared"
	"go-sample/internal/modules/catalog/trademark/module"

	"github.com/mondegor/go-storage/mrstorage"
	"github.com/mondegor/go-webcore/mrenum"
	"github.com/mondegor/go-webcore/mrtype"
)

type (
	// TrademarkPostgres - comment struct.
	TrademarkPostgres struct {
		client    mrstorage.DBConnManager
		sqlSelect mrstorage.SQLBuilderSelect
	}
)

// NewTrademarkPostgres - comment func.
func NewTrademarkPostgres(client mrstorage.DBConnManager, sqlSelect mrstorage.SQLBuilderSelect) *TrademarkPostgres {
	return &TrademarkPostgres{
		client:    client,
		sqlSelect: sqlSelect,
	}
}

// NewSelectParams - comment method.
func (re *TrademarkPostgres) NewSelectParams(params entity.TrademarkParams) mrstorage.SQLSelectParams {
	return mrstorage.SQLSelectParams{
		Where: re.sqlSelect.Where(func(w mrstorage.SQLBuilderWhere) mrstorage.SQLBuilderPartFunc {
			return w.JoinAnd(
				w.Expr("deleted_at IS NULL"),
				w.FilterLike("UPPER(trademark_caption)", strings.ToUpper(params.Filter.SearchText)),
				w.FilterAnyOf("trademark_status", params.Filter.Statuses),
			)
		}),
		OrderBy: re.sqlSelect.OrderBy(func(s mrstorage.SQLBuilderOrderBy) mrstorage.SQLBuilderPartFunc {
			return s.Join(
				s.Field(params.Sorter.FieldName, params.Sorter.Direction),
				s.Field("trademark_id", mrenum.SortDirectionASC),
			)
		}),
		Limit: re.sqlSelect.Limit(func(p mrstorage.SQLBuilderLimit) mrstorage.SQLBuilderPartFunc {
			return p.OffsetLimit(params.Pager.Index, params.Pager.Size)
		}),
	}
}

// Fetch - comment method.
func (re *TrademarkPostgres) Fetch(ctx context.Context, params mrstorage.SQLSelectParams) ([]entity.Trademark, error) {
	whereStr, whereArgs := params.Where.ToSQL()

	sql := `
		SELECT
			trademark_id,
			tag_version,
			trademark_caption as caption,
			trademark_status,
			created_at as createdAt,
			updated_at as updatedAt
		FROM
			` + module.DBSchema + `.` + module.DBTableNameTrademarks + `
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

	rows := make([]entity.Trademark, 0)

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

// FetchTotal - comment method.
func (re *TrademarkPostgres) FetchTotal(ctx context.Context, where mrstorage.SQLBuilderPart) (int64, error) {
	whereStr, whereArgs := where.ToSQL()

	sql := `
		SELECT
			COUNT(*)
		FROM
			` + module.DBSchema + `.` + module.DBTableNameTrademarks + `
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
func (re *TrademarkPostgres) FetchOne(ctx context.Context, rowID mrtype.KeyInt32) (entity.Trademark, error) {
	sql := `
		SELECT
			tag_version,
			trademark_caption,
			trademark_status,
			created_at,
			updated_at
		FROM
			` + module.DBSchema + `.` + module.DBTableNameTrademarks + `
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
func (re *TrademarkPostgres) FetchStatus(ctx context.Context, rowID mrtype.KeyInt32) (mrenum.ItemStatus, error) {
	return repositoryshared.TrademarkFetchStatusPostgres(ctx, re.client, rowID)
}

// Insert - comment method.
func (re *TrademarkPostgres) Insert(ctx context.Context, row entity.Trademark) (mrtype.KeyInt32, error) {
	sql := `
		INSERT INTO ` + module.DBSchema + `.` + module.DBTableNameTrademarks + `
			(
				trademark_caption,
				trademark_status
			)
		VALUES
			($1, $2)
		RETURNING
			trademark_id;`

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
func (re *TrademarkPostgres) Update(ctx context.Context, row entity.Trademark) (int32, error) {
	sql := `
		UPDATE
			` + module.DBSchema + `.` + module.DBTableNameTrademarks + `
		SET
			tag_version = tag_version + 1,
			updated_at = NOW(),
			trademark_caption = $3
		WHERE
			trademark_id = $1 AND tag_version = $2 AND deleted_at IS NULL
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
func (re *TrademarkPostgres) UpdateStatus(ctx context.Context, row entity.Trademark) (int32, error) {
	sql := `
		UPDATE
			` + module.DBSchema + `.` + module.DBTableNameTrademarks + `
		SET
			tag_version = tag_version + 1,
			updated_at = NOW(),
			trademark_status = $3
		WHERE
			trademark_id = $1 AND tag_version = $2 AND deleted_at IS NULL
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
func (re *TrademarkPostgres) Delete(ctx context.Context, rowID mrtype.KeyInt32) error {
	sql := `
		UPDATE
			` + module.DBSchema + `.` + module.DBTableNameTrademarks + `
		SET
			tag_version = tag_version + 1,
			deleted_at = NOW()
		WHERE
			trademark_id = $1 AND deleted_at IS NULL;`

	return re.client.Conn(ctx).Exec(
		ctx,
		sql,
		rowID,
	)
}
