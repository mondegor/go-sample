package repository_api

import (
	"context"
	repository_shared "go-sample/internal/modules/catalog/trademark/infrastructure/repository/shared"

	"github.com/mondegor/go-storage/mrstorage"
	"github.com/mondegor/go-webcore/mrenum"
	"github.com/mondegor/go-webcore/mrtype"
)

type (
	TrademarkPostgres struct {
		client mrstorage.DBConn
	}
)

func NewTrademarkPostgres(
	client mrstorage.DBConn,
) *TrademarkPostgres {
	return &TrademarkPostgres{
		client: client,
	}
}

// FetchStatus
// result: mrenum.ItemStatus - exists, ErrStorageNoRowFound - not exists, error - query error
func (re *TrademarkPostgres) FetchStatus(ctx context.Context, rowID mrtype.KeyInt32) (mrenum.ItemStatus, error) {
	return repository_shared.TrademarkFetchStatusPostgres(ctx, re.client, rowID)
}
