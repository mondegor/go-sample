package repository_shared

import (
	"context"
	module "go-sample/internal/modules/catalog/category"

	"github.com/google/uuid"
	"github.com/mondegor/go-storage/mrstorage"
	"github.com/mondegor/go-webcore/mrenum"
)

// CategoryFetchStatusPostgres
// result: mrenum.ItemStatus - exists, ErrStorageNoRowFound - not exists, error - query error
func CategoryFetchStatusPostgres(ctx context.Context, conn mrstorage.DBConn, rowID uuid.UUID) (mrenum.ItemStatus, error) {
	sql := `
		SELECT
			category_status
		FROM
			` + module.DBSchema + `.categories
		WHERE
			category_id = $1 AND deleted_at IS NULL
		LIMIT 1;`

	var status mrenum.ItemStatus

	err := conn.QueryRow(
		ctx,
		sql,
		rowID,
	).Scan(
		&status,
	)

	return status, err
}
