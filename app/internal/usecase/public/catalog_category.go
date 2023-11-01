package usecase

import (
    "context"
    "go-sample/internal/entity/public"

    "github.com/mondegor/go-storage/mrentity"
    "github.com/mondegor/go-sysmess/mrerr"
    "github.com/mondegor/go-webcore/mrcore"
    "github.com/mondegor/go-webcore/mrtool"
)

type (
    CatalogCategory struct {
        storage CatalogCategoryStorage
        serviceHelper *mrtool.ServiceHelper
    }
)

func NewCatalogCategory(
    storage CatalogCategoryStorage,
    serviceHelper *mrtool.ServiceHelper,
) *CatalogCategory {
    return &CatalogCategory{
        storage: storage,
        serviceHelper: serviceHelper,
    }
}

func (uc *CatalogCategory) GetList(ctx context.Context, params entity.CatalogCategoryParams) ([]entity.CatalogCategory, int64, error) {
    fetchParams := uc.storage.NewFetchParams(params)
    total, err := uc.storage.FetchTotal(ctx, fetchParams.Where)

    if err != nil {
        return nil, 0, mrcore.FactoryErrServiceTemporarilyUnavailable.Wrap(err, entity.ModelNameCatalogCategory)
    }

    if total < 1 {
        return []entity.CatalogCategory{}, 0, nil
    }

    items, err := uc.storage.Fetch(ctx, fetchParams)

    if err != nil {
        return nil, 0, mrcore.FactoryErrServiceTemporarilyUnavailable.Wrap(err, entity.ModelNameCatalogCategory)
    }

    return items, total, nil
}

func (uc *CatalogCategory) GetItem(ctx context.Context, id mrentity.KeyInt32) (*entity.CatalogCategory, error) {
    if id < 1 {
        return nil, mrcore.FactoryErrServiceIncorrectInputData.New(mrerr.Arg{"id": id})
    }

    item := &entity.CatalogCategory{Id: id}
    err := uc.storage.LoadOne(ctx, item)

    if err != nil {
        return nil, uc.serviceHelper.WrapErrorForSelect(err, entity.ModelNameCatalogCategory)
    }

    return item, nil
}
