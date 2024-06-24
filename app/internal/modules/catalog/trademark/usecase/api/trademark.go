package usecase_api

import (
	"context"

	"go-sample/pkg/modules/catalog/api"

	"github.com/mondegor/go-sysmess/mrmsg"
	"github.com/mondegor/go-webcore/mrcore"
	"github.com/mondegor/go-webcore/mrenum"
	"github.com/mondegor/go-webcore/mrlog"
	"github.com/mondegor/go-webcore/mrtype"
)

const (
	trademarkAPIName = "Catalog.TrademarkAPI"
)

type (
	// Trademark - comment struct.
	Trademark struct {
		storage      TrademarkStorage
		errorWrapper mrcore.UsecaseErrorWrapper
	}
)

// NewTrademark - comment func.
func NewTrademark(storage TrademarkStorage, errorWrapper mrcore.UsecaseErrorWrapper) *Trademark {
	return &Trademark{
		storage:      storage,
		errorWrapper: errorWrapper,
	}
}

// CheckingAvailability - comment method.
func (uc *Trademark) CheckingAvailability(ctx context.Context, itemID mrtype.KeyInt32) error {
	uc.debugCmd(ctx, "CheckingAvailability", mrmsg.Data{"id": itemID})

	if itemID < 1 {
		return api.ErrTrademarkRequired.New()
	}

	if status, err := uc.storage.FetchStatus(ctx, itemID); err != nil {
		if uc.errorWrapper.IsNotFoundError(err) {
			return api.ErrTrademarkNotFound.New(itemID)
		}

		return uc.errorWrapper.WrapErrorFailed(err, trademarkAPIName)
	} else if status != mrenum.ItemStatusEnabled {
		return api.ErrTrademarkNotAvailable.New(itemID)
	}

	return nil
}

func (uc *Trademark) debugCmd(ctx context.Context, command string, data mrmsg.Data) {
	mrlog.Ctx(ctx).
		Debug().
		Str("storage", trademarkAPIName).
		Str("cmd", command).
		Any("data", data).
		Send()
}
