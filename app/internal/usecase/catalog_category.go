package usecase

import (
    "context"
    "go-sample/internal/entity"

    "github.com/mondegor/go-components/mrcom"
    "github.com/mondegor/go-storage/mrentity"
    "github.com/mondegor/go-sysmess/mrerr"
    "github.com/mondegor/go-webcore/mrcore"
    "github.com/mondegor/go-webcore/mrtool"
)

type CatalogCategory struct {
    storage CatalogCategoryStorage
    eventBox mrcore.EventBox
    serviceHelper *mrtool.ServiceHelper
    statusFlow mrcom.ItemStatusFlow
}

func NewCatalogCategory(storage CatalogCategoryStorage,
                        eventBox mrcore.EventBox,
                        serviceHelper *mrtool.ServiceHelper) *CatalogCategory {
    return &CatalogCategory{
        storage: storage,
        eventBox: eventBox,
        serviceHelper: serviceHelper,
        statusFlow: mrcom.ItemStatusFlowDefault,
    }
}

func (uc *CatalogCategory) GetList(ctx context.Context, listFilter *entity.CatalogCategoryListFilter) ([]entity.CatalogCategory, error) {
    items := make([]entity.CatalogCategory, 0, 16)
    err := uc.storage.LoadAll(ctx, listFilter, &items)

    if err != nil {
        return nil, mrcore.FactoryErrServiceEntityTemporarilyUnavailable.Wrap(err, entity.ModelNameCatalogCategory)
    }

    return items, nil
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

func (uc *CatalogCategory) CheckAvailability(ctx context.Context, id mrentity.KeyInt32) error {
    if id < 1 {
        return mrcore.FactoryErrServiceIncorrectInputData.New(mrerr.Arg{"id": id})
    }

    err := uc.storage.IsExists(ctx, id)

    return uc.serviceHelper.ReturnErrorIfItemNotFound(err, entity.ModelNameCatalogCategory)
}

// Create
// modifies: item{Id}
func (uc *CatalogCategory) Create(ctx context.Context, item *entity.CatalogCategory) error {
    item.Status = mrcom.ItemStatusDraft
    err := uc.storage.Insert(ctx, item)

    if err != nil {
        return mrcore.FactoryErrServiceEntityNotCreated.Wrap(err, entity.ModelNameCatalogCategory)
    }

    uc.eventBox.Emit(
        "%s::Create: id=%d",
        entity.ModelNameCatalogCategory,
        item.Id,
    )

    return nil
}

func (uc *CatalogCategory) Store(ctx context.Context, item *entity.CatalogCategory) error {
    if item.Id < 1 || item.Version < 1 {
        return mrcore.FactoryErrServiceIncorrectInputData.New(mrerr.Arg{"item.Id": item.Id, "Item.Version": item.Version})
    }

    err := uc.storage.Update(ctx, item)

    if err != nil {
        return uc.serviceHelper.WrapErrorForUpdate(err, entity.ModelNameCatalogCategory)
    }

    uc.eventBox.Emit(
        "%s::Store: id=%d",
        entity.ModelNameCatalogCategory,
        item.Id,
    )

    return nil
}

func (uc *CatalogCategory) ChangeStatus(ctx context.Context, item *entity.CatalogCategory) error {
    if item.Id < 1 || item.Version < 1 {
        return mrcore.FactoryErrServiceIncorrectInputData.New(mrerr.Arg{"item.Id": item.Id, "Item.Version": item.Version})
    }

    currentStatus, err := uc.storage.FetchStatus(ctx, item)

    if err != nil {
        return uc.serviceHelper.WrapErrorForSelect(err, entity.ModelNameCatalogCategory)
    }

    if !uc.statusFlow.Check(currentStatus, item.Status) {
        return mrcore.FactoryErrServiceIncorrectSwitchStatus.New(currentStatus, item.Status, entity.ModelNameCatalogCategory, item.Id)
    }

    err = uc.storage.UpdateStatus(ctx, item)

    if err != nil {
        return uc.serviceHelper.WrapErrorForUpdate(err, entity.ModelNameCatalogCategory)
    }

    uc.eventBox.Emit(
        "%s::ChangeStatus: id=%d, status=%s",
        entity.ModelNameCatalogCategory,
        item.Id,
        item.Status,
    )

    return nil
}

func (uc *CatalogCategory) Remove(ctx context.Context, id mrentity.KeyInt32) error {
    if id < 1 {
        return mrcore.FactoryErrServiceIncorrectInputData.New(mrerr.Arg{"id": id})
    }

    err := uc.storage.Delete(ctx, id)

    if err != nil {
        return uc.serviceHelper.WrapErrorForRemove(err, entity.ModelNameCatalogCategory)
    }

    uc.eventBox.Emit(
        "%s::Remove: id=%d",
        entity.ModelNameCatalogCategory,
        id,
    )

    return nil
}
