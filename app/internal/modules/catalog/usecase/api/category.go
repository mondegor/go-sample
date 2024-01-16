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
	Category struct {
		storage       CategoryStorage
		serviceHelper *mrtool.ServiceHelper
	}

	CategoryStorage interface {
		IsExists(ctx context.Context, id mrtype.KeyInt32) error
	}
)

func NewCategory(
	storage CategoryStorage,
	serviceHelper *mrtool.ServiceHelper,
) *Category {
	return &Category{
		storage:       storage,
		serviceHelper: serviceHelper,
	}
}

func (uc *Category) CheckingAvailability(ctx context.Context, id mrtype.KeyInt32) error {
	uc.debugCmd(ctx, "CheckingAvailability", mrmsg.Data{"id": id})

	if id < 1 {
		return catalog.FactoryErrCategoryNotFound.New(id)
	}

	if err := uc.storage.IsExists(ctx, id); err != nil {
		if uc.serviceHelper.IsNotFoundError(err) {
			return catalog.FactoryErrCategoryNotFound.New(id)
		}

		return uc.serviceHelper.WrapErrorFailed(err, "Catalog.CategoryAPI")
	}

	return nil
}

func (uc *Category) debugCmd(ctx context.Context, command string, data mrmsg.Data) {
	mrctx.Logger(ctx).Debug(
		"Catalog.CategoryAPI: cmd=%s, data=%s",
		command,
		data,
	)
}
