package usecase

import (
	"context"

	"github.com/mondegor/go-webcore/mrtype"
)

type (
	TrademarkServiceAPI interface {
		CheckingAvailability(ctx context.Context, id mrtype.KeyInt32) error
	}

	TrademarkStorage interface {
		IsExists(ctx context.Context, id mrtype.KeyInt32) error
	}
)
