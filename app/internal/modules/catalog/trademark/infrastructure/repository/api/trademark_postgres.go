package repository_api

import (
	"context"

	repositoryshared "go-sample/internal/modules/catalog/trademark/infrastructure/repository/shared"

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
	return repositoryshared.TrademarkFetchStatusPostgres(ctx, re.client, rowID)
}
