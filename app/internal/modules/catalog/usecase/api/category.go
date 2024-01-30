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
	categoryAPIName = "Catalog.CategoryAPI"
)

type (
	Category struct {
		storage       CategoryStorage
		usecaseHelper *mrcore.UsecaseHelper
	}

	CategoryStorage interface {
		IsExists(ctx context.Context, id mrtype.KeyInt32) error
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

func (uc *Category) CheckingAvailability(ctx context.Context, id mrtype.KeyInt32) error {
	uc.debugCmd(ctx, "CheckingAvailability", mrmsg.Data{"id": id})

	if id < 1 {
		return catalog.FactoryErrCategoryNotFound.New(id)
	}

	if err := uc.storage.IsExists(ctx, id); err != nil {
		if uc.usecaseHelper.IsNotFoundError(err) {
			return catalog.FactoryErrCategoryNotFound.New(id)
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
