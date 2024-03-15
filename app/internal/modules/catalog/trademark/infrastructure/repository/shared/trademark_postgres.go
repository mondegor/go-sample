package repository_shared

import (
	"context"
	module "go-sample/internal/modules/catalog/trademark"

	"github.com/mondegor/go-storage/mrstorage"
	"github.com/mondegor/go-webcore/mrenum"
	"github.com/mondegor/go-webcore/mrtype"
)

// TrademarkIsExistsPostgres
// result: nil - exists, ErrStorageNoRowFound - not exists, error - query error
func TrademarkIsExistsPostgres(ctx context.Context, conn mrstorage.DBConn, rowID mrtype.KeyInt32) error {
	sql := `
		SELECT
			1
		FROM
			` + module.DBSchema + `.trademarks
		WHERE
			trademark_id = $1 AND trademark_status <> $2
		LIMIT 1;`

	return conn.QueryRow(
		ctx,
		sql,
		rowID,
		mrenum.ItemStatusRemoved,
	).Scan(
		&rowID,
	)
}
