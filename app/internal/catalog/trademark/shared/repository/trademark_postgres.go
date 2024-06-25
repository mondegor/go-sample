package repository

import (
	"context"

	"github.com/mondegor/go-sample/internal/catalog/trademark/module"

	"github.com/mondegor/go-storage/mrstorage"
	"github.com/mondegor/go-webcore/mrenum"
	"github.com/mondegor/go-webcore/mrtype"
)

// TrademarkFetchStatusPostgres - comment func.
// result: mrenum.ItemStatus - exists, ErrStorageNoRowFound - not exists, error - query error.
func TrademarkFetchStatusPostgres(ctx context.Context, client mrstorage.DBConnManager, rowID mrtype.KeyInt32) (mrenum.ItemStatus, error) {
	sql := `
		SELECT
			trademark_status
		FROM
			` + module.DBSchema + `.` + module.DBTableNameTrademarks + `
		WHERE
			trademark_id = $1 AND deleted_at IS NULL
		LIMIT 1;`

	var status mrenum.ItemStatus

	err := client.Conn(ctx).QueryRow(
		ctx,
		sql,
		rowID,
	).Scan(
		&status,
	)

	return status, err
}
