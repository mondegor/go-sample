package usecase_api

import (
	"context"
	"go-sample/pkg/modules/catalog"

	"github.com/mondegor/go-sysmess/mrmsg"
	"github.com/mondegor/go-webcore/mrctx"
	"github.com/mondegor/go-webcore/mrtool"
	"github.com/mondegor/go-webcore/mrtype"
)

type (
	Trademark struct {
		storage       TrademarkStorage
		serviceHelper *mrtool.ServiceHelper
	}

	TrademarkStorage interface {
		IsExists(ctx context.Context, id mrtype.KeyInt32) error
	}
)

func NewTrademark(
	storage TrademarkStorage,
	serviceHelper *mrtool.ServiceHelper,
) *Trademark {
	return &Trademark{
		storage:       storage,
		serviceHelper: serviceHelper,
	}
}

func (uc *Trademark) CheckingAvailability(ctx context.Context, id mrtype.KeyInt32) error {
	uc.debugCmd(ctx, "CheckingAvailability", mrmsg.Data{"id": id})

	if id < 1 {
		return catalog.FactoryErrTrademarkNotFound.New(id)
	}

	if err := uc.storage.IsExists(ctx, id); err != nil {
		if uc.serviceHelper.IsNotFoundError(err) {
			return catalog.FactoryErrTrademarkNotFound.New(id)
		}

		return uc.serviceHelper.WrapErrorFailed(err, "Catalog.TrademarkAPI")
	}

	return nil
}

func (uc *Trademark) debugCmd(ctx context.Context, command string, data mrmsg.Data) {
	mrctx.Logger(ctx).Debug(
		"Catalog.TrademarkAPI: cmd=%s, data=%s",
		command,
		data,
	)
}
