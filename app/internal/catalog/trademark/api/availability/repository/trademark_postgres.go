package repository

import (
	"context"

	"github.com/mondegor/go-storage/mrpostgres/db"
	"github.com/mondegor/go-storage/mrstorage"
	"github.com/mondegor/go-webcore/mrenum"

	"github.com/mondegor/go-sample/internal/catalog/trademark/module"
)

type (
	// TrademarkPostgres - comment struct.
	TrademarkPostgres struct {
		repoStatus db.FieldFetcher[uint64, mrenum.ItemStatus]
	}
)

// NewTrademarkPostgres - создаёт объект TrademarkPostgres.
func NewTrademarkPostgres(client mrstorage.DBConnManager) *TrademarkPostgres {
	return &TrademarkPostgres{
		repoStatus: db.NewFieldFetcher[uint64, mrenum.ItemStatus](
			client,
			module.DBTableNameTrademarks,
			"trademark_id",
			"trademark_status",
			module.DBFieldDeletedAt,
		),
	}
}

// FetchStatus - comment method.
// result: mrenum.ItemStatus - exists, ErrStorageNoRowFound - not exists, error - query error.
func (re *TrademarkPostgres) FetchStatus(ctx context.Context, rowID uint64) (mrenum.ItemStatus, error) {
	return re.repoStatus.Fetch(ctx, rowID)
}
