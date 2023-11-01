package usecase

import (
    "context"
    "go-sample/internal/entity/admin-panel"

    "github.com/mondegor/go-storage/mrentity"
    "github.com/mondegor/go-sysmess/mrerr"
    "github.com/mondegor/go-webcore/mrcore"
    "github.com/mondegor/go-webcore/mrenum"
    "github.com/mondegor/go-webcore/mrtool"
)

type (
    CatalogCategory struct {
        storage CatalogCategoryStorage
        eventBox mrcore.EventBox
        serviceHelper *mrtool.ServiceHelper
        statusFlow mrenum.StatusFlow
    }
)

func NewCatalogCategory(
    storage CatalogCategoryStorage,
    eventBox mrcore.EventBox,
    serviceHelper *mrtool.ServiceHelper,
) *CatalogCategory {
    return &CatalogCategory{
        storage: storage,
        eventBox: eventBox,
        serviceHelper: serviceHelper,
        statusFlow: mrenum.ItemStatusFlow,
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
    item.Status = mrenum.ItemStatusDraft
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
        return mrcore.FactoryErrServiceIncorrectInputData.New(mrerr.Arg{"item.Id": item.Id, "Version": item.Version})
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
        return mrcore.FactoryErrServiceIncorrectInputData.New(mrerr.Arg{"item.Id": item.Id, "Version": item.Version})
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
