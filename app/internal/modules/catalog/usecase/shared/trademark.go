package usecase

import (
	"context"

	"github.com/mondegor/go-webcore/mrtool"
	"github.com/mondegor/go-webcore/mrtype"
)

type (
	TrademarkAPI struct {
		storage       TrademarkStorage
		serviceHelper *mrtool.ServiceHelper
	}
)

func NewTrademarkAPI(
	storage TrademarkStorage,
	serviceHelper *mrtool.ServiceHelper,
) *TrademarkAPI {
	return &TrademarkAPI{
		storage:       storage,
		serviceHelper: serviceHelper,
	}
}

func (uc *TrademarkAPI) CheckingAvailability(ctx context.Context, id mrtype.KeyInt32) error {
	if id < 1 {
		return FactoryErrTrademarkNotFound.New(id)
	}

	if err := uc.storage.IsExists(ctx, id); err != nil {
		if uc.serviceHelper.IsNotFound(err) {
			return FactoryErrTrademarkNotFound.New(id)
		}

		return uc.serviceHelper.WrapErrorFailed(err, "CatalogTrademarkAPI")
	}

	return nil
}
