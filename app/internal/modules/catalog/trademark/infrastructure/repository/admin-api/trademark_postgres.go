package repository

import (
	"context"
	module "go-sample/internal/modules/catalog/trademark"
	"go-sample/internal/modules/catalog/trademark/entity/admin-api"
	repository_shared "go-sample/internal/modules/catalog/trademark/infrastructure/repository/shared"
	"strings"

	"github.com/mondegor/go-storage/mrstorage"
	"github.com/mondegor/go-webcore/mrenum"
	"github.com/mondegor/go-webcore/mrtype"
)

type (
	TrademarkPostgres struct {
		client    mrstorage.DBConn
		sqlSelect mrstorage.SqlBuilderSelect
	}
)

func NewTrademarkPostgres(
	client mrstorage.DBConn,
	sqlSelect mrstorage.SqlBuilderSelect,
) *TrademarkPostgres {
	return &TrademarkPostgres{
		client:    client,
		sqlSelect: sqlSelect,
	}
}

func (re *TrademarkPostgres) NewSelectParams(params entity.TrademarkParams) mrstorage.SqlSelectParams {
	return mrstorage.SqlSelectParams{
		Where: re.sqlSelect.Where(func(w mrstorage.SqlBuilderWhere) mrstorage.SqlBuilderPartFunc {
			return w.JoinAnd(
				w.NotEqual("trademark_status", mrenum.ItemStatusRemoved),
				w.FilterLike("UPPER(trademark_caption)", strings.ToUpper(params.Filter.SearchText)),
				w.FilterAnyOf("trademark_status", params.Filter.Statuses),
			)
		}),
		OrderBy: re.sqlSelect.OrderBy(func(s mrstorage.SqlBuilderOrderBy) mrstorage.SqlBuilderPartFunc {
			return s.Join(
				s.Field(params.Sorter.FieldName, params.Sorter.Direction),
				s.Field("trademark_id", mrenum.SortDirectionASC),
			)
		}),
		Pager: re.sqlSelect.Pager(func(p mrstorage.SqlBuilderPager) mrstorage.SqlBuilderPartFunc {
			return p.OffsetLimit(params.Pager.Index, params.Pager.Size)
		}),
	}
}

func (re *TrademarkPostgres) Fetch(ctx context.Context, params mrstorage.SqlSelectParams) ([]entity.Trademark, error) {
	whereStr, whereArgs := params.Where.ToSql()

	sql := `
		SELECT
			trademark_id,
			tag_version,
			trademark_caption as caption,
			trademark_status,
			created_at as createdAt,
			updated_at as updatedAt
		FROM
			` + module.DBSchema + `.trademarks
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

func (re *TrademarkPostgres) FetchTotal(ctx context.Context, where mrstorage.SqlBuilderPart) (int64, error) {
	whereStr, whereArgs := where.ToSql()

	sql := `
		SELECT
			COUNT(*)
		FROM
			` + module.DBSchema + `.trademarks
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

func (re *TrademarkPostgres) FetchOne(ctx context.Context, rowID mrtype.KeyInt32) (entity.Trademark, error) {
	sql := `
		SELECT
			tag_version,
			trademark_caption,
			trademark_status,
			created_at,
			updated_at
		FROM
			` + module.DBSchema + `.trademarks
		WHERE
			trademark_id = $1 AND trademark_status <> $2
		LIMIT 1;`

	row := entity.Trademark{ID: rowID}

	err := re.client.QueryRow(
		ctx,
		sql,
		rowID,
		mrenum.ItemStatusRemoved,
	).Scan(
		&row.TagVersion,
		&row.Caption,
		&row.Status,
		&row.CreatedAt,
		&row.UpdatedAt,
	)

	return row, err
}

func (re *TrademarkPostgres) FetchStatus(ctx context.Context, row entity.Trademark) (mrenum.ItemStatus, error) {
	sql := `
		SELECT
			trademark_status
		FROM
			` + module.DBSchema + `.trademarks
		WHERE
			trademark_id = $1 AND trademark_status <> $2
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
func (re *TrademarkPostgres) IsExists(ctx context.Context, rowID mrtype.KeyInt32) error {
	return repository_shared.TrademarkIsExistsPostgres(ctx, re.client, rowID)
}

func (re *TrademarkPostgres) Insert(ctx context.Context, row entity.Trademark) (mrtype.KeyInt32, error) {
	sql := `
		INSERT INTO ` + module.DBSchema + `.trademarks
			(
				trademark_caption,
				trademark_status
			)
		VALUES
			($1, $2)
		RETURNING
			trademark_id;`

	err := re.client.QueryRow(
		ctx,
		sql,
		row.Caption,
		row.Status,
	).Scan(
		&row.ID,
	)

	return row.ID, err
}

func (re *TrademarkPostgres) Update(ctx context.Context, row entity.Trademark) (int32, error) {
	sql := `
		UPDATE
			` + module.DBSchema + `.trademarks
		SET
			tag_version = tag_version + 1,
			updated_at = NOW(),
			trademark_caption = $4
		WHERE
			trademark_id = $1 AND tag_version = $2 AND trademark_status <> $3
		RETURNING
			tag_version;`

	var tagVersion int32

	err := re.client.QueryRow(
		ctx,
		sql,
		row.ID,
		row.TagVersion,
		mrenum.ItemStatusRemoved,
		row.Caption,
	).Scan(
		&tagVersion,
	)

	return tagVersion, err
}

func (re *TrademarkPostgres) UpdateStatus(ctx context.Context, row entity.Trademark) (int32, error) {
	sql := `
		UPDATE
			` + module.DBSchema + `.trademarks
		SET
			tag_version = tag_version + 1,
			updated_at = NOW(),
			trademark_status = $4
		WHERE
			trademark_id = $1 AND tag_version = $2 AND trademark_status <> $3
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

func (re *TrademarkPostgres) Delete(ctx context.Context, rowID mrtype.KeyInt32) error {
	sql := `
		UPDATE
			` + module.DBSchema + `.trademarks
		SET
			tag_version = tag_version + 1,
			updated_at = NOW(),
			trademark_status = $2
		WHERE
			trademark_id = $1 AND trademark_status <> $2;`

	return re.client.Exec(
		ctx,
		sql,
		rowID,
		mrenum.ItemStatusRemoved,
	)
}
