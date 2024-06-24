package shared

import (
	"context"

	"go-sample/internal/modules/catalog/category/module"

	"github.com/google/uuid"
	"github.com/mondegor/go-storage/mrstorage"
	"github.com/mondegor/go-webcore/mrenum"
)

// CategoryFetchStatusPostgres - comment func.
// result: mrenum.ItemStatus - exists, ErrStorageNoRowFound - not exists, error - query error.
func CategoryFetchStatusPostgres(ctx context.Context, client mrstorage.DBConnManager, rowID uuid.UUID) (mrenum.ItemStatus, error) {
	sql := `
		SELECT
			category_status
		FROM
			` + module.DBSchema + `.` + module.DBTableNameCategories + `
		WHERE
			category_id = $1 AND deleted_at IS NULL
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
