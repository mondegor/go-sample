package usecase

import (
	"context"

	"github.com/mondegor/go-webcore/mrtool"
	"github.com/mondegor/go-webcore/mrtype"
)

type (
	CategoryAPI struct {
		storage       CategoryStorage
		serviceHelper *mrtool.ServiceHelper
	}
)

func NewCategoryAPI(
	storage CategoryStorage,
	serviceHelper *mrtool.ServiceHelper,
) *CategoryAPI {
	return &CategoryAPI{
		storage:       storage,
		serviceHelper: serviceHelper,
	}
}

func (uc *CategoryAPI) CheckingAvailability(ctx context.Context, id mrtype.KeyInt32) error {
	if id < 1 {
		return FactoryErrCategoryNotFound.New(id)
	}

	if err := uc.storage.IsExists(ctx, id); err != nil {
		if uc.serviceHelper.IsNotFound(err) {
			return FactoryErrCategoryNotFound.New(id)
		}

		return uc.serviceHelper.WrapErrorFailed(err, "CatalogCategoryAPI")
	}

	return nil
}
