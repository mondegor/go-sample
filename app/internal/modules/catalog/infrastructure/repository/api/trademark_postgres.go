package repository_api

import (
	"context"
	repository_shared "go-sample/internal/modules/catalog/infrastructure/repository/shared"

	"github.com/mondegor/go-storage/mrstorage"
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

// IsExists
// result: nil - exists, ErrStorageNoRowFound - not exists, error - query error
func (re *TrademarkPostgres) IsExists(ctx context.Context, id mrtype.KeyInt32) error {
	return repository_shared.TrademarkIsExistsPostgres(ctx, re.client, id)
}
