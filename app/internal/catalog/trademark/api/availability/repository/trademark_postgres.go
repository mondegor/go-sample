package repository

import (
	"context"

	"github.com/mondegor/go-sample/internal/catalog/trademark/shared/repository"

	"github.com/mondegor/go-storage/mrstorage"
	"github.com/mondegor/go-webcore/mrenum"
	"github.com/mondegor/go-webcore/mrtype"
)

type (
	// TrademarkPostgres - comment struct.
	TrademarkPostgres struct {
		client mrstorage.DBConnManager
	}
)

// NewTrademarkPostgres - comment func.
func NewTrademarkPostgres(client mrstorage.DBConnManager) *TrademarkPostgres {
	return &TrademarkPostgres{
		client: client,
	}
}

// FetchStatus - comment method.
// result: mrenum.ItemStatus - exists, ErrStorageNoRowFound - not exists, error - query error.
func (re *TrademarkPostgres) FetchStatus(ctx context.Context, rowID mrtype.KeyInt32) (mrenum.ItemStatus, error) {
	return repository.TrademarkFetchStatusPostgres(ctx, re.client, rowID)
}
