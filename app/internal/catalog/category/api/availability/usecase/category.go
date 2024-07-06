package usecase

import (
	"context"

	"github.com/mondegor/go-sample/internal/catalog/category/api/availability"
	"github.com/mondegor/go-sample/pkg/catalog/api"

	"github.com/google/uuid"
	"github.com/mondegor/go-sysmess/mrmsg"
	"github.com/mondegor/go-webcore/mrcore"
	"github.com/mondegor/go-webcore/mrenum"
	"github.com/mondegor/go-webcore/mrlog"
)

type (
	// Category - comment struct.
	Category struct {
		storage      availability.CategoryStorage
		errorWrapper mrcore.UsecaseErrorWrapper
	}
)

// NewCategory - создаёт объект mrcore.UsecaseErrorWrapper.
func NewCategory(storage availability.CategoryStorage, errorWrapper mrcore.UsecaseErrorWrapper) *Category {
	return &Category{
		storage:      storage,
		errorWrapper: errorWrapper,
	}
}

// CheckingAvailability - comment method.
func (uc *Category) CheckingAvailability(ctx context.Context, itemID uuid.UUID) error {
	uc.debugCmd(ctx, "CheckingAvailability", mrmsg.Data{"id": itemID})

	if itemID == uuid.Nil {
		return api.ErrCategoryRequired.New()
	}

	if status, err := uc.storage.FetchStatus(ctx, itemID); err != nil {
		if uc.errorWrapper.IsNotFoundError(err) {
			return api.ErrCategoryNotFound.New(itemID)
		}

		return uc.errorWrapper.WrapErrorFailed(err, api.CategoryAvailabilityName)
	} else if status != mrenum.ItemStatusEnabled {
		return api.ErrCategoryNotAvailable.New(itemID)
	}

	return nil
}

func (uc *Category) debugCmd(ctx context.Context, command string, data mrmsg.Data) {
	mrlog.Ctx(ctx).
		Debug().
		Str("storage", api.CategoryAvailabilityName).
		Str("cmd", command).
		Any("data", data).
		Send()
}
