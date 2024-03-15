package usecase_api

import (
	"context"
	"go-sample/pkg/modules/catalog"

	"github.com/mondegor/go-sysmess/mrmsg"
	"github.com/mondegor/go-webcore/mrcore"
	"github.com/mondegor/go-webcore/mrlog"
	"github.com/mondegor/go-webcore/mrtype"
)

const (
	trademarkAPIName = "Catalog.TrademarkAPI"
)

type (
	Trademark struct {
		storage       TrademarkStorage
		usecaseHelper *mrcore.UsecaseHelper
	}

	TrademarkStorage interface {
		IsExists(ctx context.Context, rowID mrtype.KeyInt32) error
	}
)

func NewTrademark(
	storage TrademarkStorage,
	usecaseHelper *mrcore.UsecaseHelper,
) *Trademark {
	return &Trademark{
		storage:       storage,
		usecaseHelper: usecaseHelper,
	}
}

func (uc *Trademark) CheckingAvailability(ctx context.Context, itemID mrtype.KeyInt32) error {
	uc.debugCmd(ctx, "CheckingAvailability", mrmsg.Data{"id": itemID})

	if itemID < 1 {
		return catalog.FactoryErrTrademarkNotFound.New(itemID)
	}

	if err := uc.storage.IsExists(ctx, itemID); err != nil {
		if uc.usecaseHelper.IsNotFoundError(err) {
			return catalog.FactoryErrTrademarkNotFound.New(itemID)
		}

		return uc.usecaseHelper.WrapErrorFailed(err, trademarkAPIName)
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
