package repository_shared

import (
	"context"
	module "go-sample/internal/modules/catalog/category"

	"github.com/mondegor/go-storage/mrstorage"
	"github.com/mondegor/go-webcore/mrenum"
	"github.com/mondegor/go-webcore/mrtype"
)

// CategoryIsExistsPostgres
// result: nil - exists, ErrStorageNoRowFound - not exists, error - query error
func CategoryIsExistsPostgres(ctx context.Context, conn mrstorage.DBConn, id mrtype.KeyInt32) error {
	sql := `
		SELECT
			1
		FROM
			` + module.DBSchema + `.categories
		WHERE
			category_id = $1 AND category_status <> $2
		LIMIT 1;`

	return conn.QueryRow(
		ctx,
		sql,
		id,
		mrenum.ItemStatusRemoved,
	).Scan(
		&id,
	)
}
