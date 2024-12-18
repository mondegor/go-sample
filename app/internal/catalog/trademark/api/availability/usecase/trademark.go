package usecase

import (
	"context"

	"github.com/mondegor/go-sysmess/mrmsg"
	"github.com/mondegor/go-webcore/mrcore"
	"github.com/mondegor/go-webcore/mrenum"
	"github.com/mondegor/go-webcore/mrlog"

	"github.com/mondegor/go-sample/internal/catalog/trademark/api/availability"
	"github.com/mondegor/go-sample/pkg/catalog/api"
)

type (
	// Trademark - comment struct.
	Trademark struct {
		storage      availability.TrademarkStorage
		errorWrapper mrcore.UseCaseErrorWrapper
	}
)

// NewTrademark - создаёт объект Trademark.
func NewTrademark(storage availability.TrademarkStorage, errorWrapper mrcore.UseCaseErrorWrapper) *Trademark {
	return &Trademark{
		storage:      storage,
		errorWrapper: errorWrapper,
	}
}

// CheckingAvailability - comment method.
func (uc *Trademark) CheckingAvailability(ctx context.Context, itemID uint64) error {
	uc.debugCmd(ctx, "CheckingAvailability", mrmsg.Data{"id": itemID})

	if itemID == 0 {
		return api.ErrTrademarkRequired.New()
	}

	if status, err := uc.storage.FetchStatus(ctx, itemID); err != nil {
		if uc.errorWrapper.IsNotFoundError(err) {
			return api.ErrTrademarkNotFound.New(itemID)
		}

		return uc.errorWrapper.WrapErrorFailed(err, api.TrademarkAvailabilityName)
	} else if status != mrenum.ItemStatusEnabled {
		return api.ErrTrademarkNotAvailable.New(itemID)
	}

	return nil
}

func (uc *Trademark) debugCmd(ctx context.Context, command string, data mrmsg.Data) {
	mrlog.Ctx(ctx).
		Debug().
		Str("storage", api.TrademarkAvailabilityName).
		Str("cmd", command).
		Any("data", data).
		Send()
}
