package usecase_api

import (
	"context"
	"go-sample/pkg/modules/catalog"

	"github.com/google/uuid"
	"github.com/mondegor/go-sysmess/mrmsg"
	"github.com/mondegor/go-webcore/mrcore"
	"github.com/mondegor/go-webcore/mrlog"
)

const (
	categoryAPIName = "Catalog.CategoryAPI"
)

type (
	Category struct {
		storage       CategoryStorage
		usecaseHelper *mrcore.UsecaseHelper
	}
)

func NewCategory(
	storage CategoryStorage,
	usecaseHelper *mrcore.UsecaseHelper,
) *Category {
	return &Category{
		storage:       storage,
		usecaseHelper: usecaseHelper,
	}
}

func (uc *Category) CheckingAvailability(ctx context.Context, itemID uuid.UUID) error {
	uc.debugCmd(ctx, "CheckingAvailability", mrmsg.Data{"id": itemID})

	if itemID == uuid.Nil {
		return catalog.FactoryErrCategoryRequired.New()
	}

	if err := uc.storage.IsExists(ctx, itemID); err != nil {
		if uc.usecaseHelper.IsNotFoundError(err) {
			return catalog.FactoryErrCategoryNotFound.New(itemID)
		}

		return uc.usecaseHelper.WrapErrorFailed(err, categoryAPIName)
	}

	return nil
}

func (uc *Category) debugCmd(ctx context.Context, command string, data mrmsg.Data) {
	mrlog.Ctx(ctx).
		Debug().
		Str("storage", categoryAPIName).
		Str("cmd", command).
		Any("data", data).
		Send()
}
