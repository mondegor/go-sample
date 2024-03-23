package repository_shared

import (
	"context"
	module "go-sample/internal/modules/catalog/trademark"

	"github.com/mondegor/go-storage/mrstorage"
	"github.com/mondegor/go-webcore/mrenum"
	"github.com/mondegor/go-webcore/mrtype"
)

// TrademarkFetchStatusPostgres
// result: mrenum.ItemStatus - exists, ErrStorageNoRowFound - not exists, error - query error
func TrademarkFetchStatusPostgres(ctx context.Context, conn mrstorage.DBConn, rowID mrtype.KeyInt32) (mrenum.ItemStatus, error) {
	sql := `
		SELECT
			trademark_status
		FROM
			` + module.DBSchema + `.trademarks
		WHERE
			trademark_id = $1 AND trademark_status <> $2
		LIMIT 1;`

	var status mrenum.ItemStatus

	err := conn.QueryRow(
		ctx,
		sql,
		rowID,
		mrenum.ItemStatusRemoved,
	).Scan(
		&status,
	)

	return status, err
}
