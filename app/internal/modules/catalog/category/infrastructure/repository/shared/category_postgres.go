package repository_shared

import (
	"context"
	module "go-sample/internal/modules/catalog/category"

	"github.com/google/uuid"
	"github.com/mondegor/go-storage/mrstorage"
	"github.com/mondegor/go-webcore/mrenum"
)

// CategoryIsExistsPostgres
// result: nil - exists, ErrStorageNoRowFound - not exists, error - query error
func CategoryIsExistsPostgres(ctx context.Context, conn mrstorage.DBConn, rowID uuid.UUID) error {
	sql := `
		SELECT
			category_id
		FROM
			` + module.DBSchema + `.categories
		WHERE
			category_id = $1 AND category_status <> $2
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
