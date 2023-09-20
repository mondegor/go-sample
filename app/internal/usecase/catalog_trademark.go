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

type (
    CatalogTrademark struct {
        storage CatalogTrademarkStorage
        eventBox mrcore.EventBox
        serviceHelper *mrtool.ServiceHelper
        statusFlow mrcom.ItemStatusFlow
    }
)

func NewCatalogTrademark(storage CatalogTrademarkStorage,
                         eventBox mrcore.EventBox,
                         serviceHelper *mrtool.ServiceHelper) *CatalogTrademark {
    return &CatalogTrademark{
        storage: storage,
        eventBox: eventBox,
        serviceHelper: serviceHelper,
        statusFlow: mrcom.ItemStatusFlowDefault,
    }
}

func (uc *CatalogTrademark) GetList(ctx context.Context, listFilter *entity.CatalogTrademarkListFilter) ([]entity.CatalogTrademark, error) {
    items := make([]entity.CatalogTrademark, 0, 4)
    err := uc.storage.LoadAll(ctx, listFilter, &items)

    if err != nil {
        return nil, mrcore.FactoryErrServiceTemporarilyUnavailable.Wrap(err, entity.ModelNameCatalogTrademark)
    }

    return items, nil
}

func (uc *CatalogTrademark) GetItem(ctx context.Context, id mrentity.KeyInt32) (*entity.CatalogTrademark, error) {
    if id < 1 {
        return nil, mrcore.FactoryErrServiceIncorrectInputData.New(mrerr.Arg{"id": id})
    }

    item := &entity.CatalogTrademark{Id: id}
    err := uc.storage.LoadOne(ctx, item)

    if err != nil {
        return nil, uc.serviceHelper.WrapErrorForSelect(err, entity.ModelNameCatalogTrademark)
    }

    return item, nil
}

// Create
// modifies: item{Id}
func (uc *CatalogTrademark) Create(ctx context.Context, item *entity.CatalogTrademark) error {
    item.Status = mrcom.ItemStatusDraft
    err := uc.storage.Insert(ctx, item)

    if err != nil {
        return mrcore.FactoryErrServiceEntityNotCreated.Wrap(err, entity.ModelNameCatalogTrademark)
    }

    uc.eventBox.Emit(
        "%s::Create: id=%d",
        entity.ModelNameCatalogTrademark,
        item.Id,
    )

    return nil
}

func (uc *CatalogTrademark) Store(ctx context.Context, item *entity.CatalogTrademark) error {
    if item.Id < 1 || item.Version < 1 {
        return mrcore.FactoryErrServiceIncorrectInputData.New(mrerr.Arg{"item.Id": item.Id, "Version": item.Version})
    }

    err := uc.storage.Update(ctx, item)

    if err != nil {
        return uc.serviceHelper.WrapErrorForUpdate(err, entity.ModelNameCatalogTrademark)
    }

    uc.eventBox.Emit(
        "%s::Store: id=%d",
        entity.ModelNameCatalogTrademark,
        item.Id,
    )

    return nil
}

func (uc *CatalogTrademark) ChangeStatus(ctx context.Context, item *entity.CatalogTrademark) error {
    if item.Id < 1 || item.Version < 1 {
        return mrcore.FactoryErrServiceIncorrectInputData.New(mrerr.Arg{"item.Id": item.Id, "Version": item.Version})
    }

    currentStatus, err := uc.storage.FetchStatus(ctx, item)

    if err != nil {
        return uc.serviceHelper.WrapErrorForSelect(err, entity.ModelNameCatalogTrademark)
    }

    if !uc.statusFlow.Check(currentStatus, item.Status) {
        return mrcore.FactoryErrServiceIncorrectSwitchStatus.New(currentStatus, item.Status, entity.ModelNameCatalogTrademark, item.Id)
    }

    err = uc.storage.UpdateStatus(ctx, item)

    if err != nil {
        return uc.serviceHelper.WrapErrorForUpdate(err, entity.ModelNameCatalogTrademark)
    }

    uc.eventBox.Emit(
        "%s::ChangeStatus: id=%d, status=%s",
        entity.ModelNameCatalogTrademark,
        item.Id,
        item.Status,
    )

    return nil
}

func (uc *CatalogTrademark) Remove(ctx context.Context, id mrentity.KeyInt32) error {
    if id < 1 {
        return mrcore.FactoryErrServiceIncorrectInputData.New(mrerr.Arg{"id": id})
    }

    err := uc.storage.Delete(ctx, id)

    if err != nil {
        return uc.serviceHelper.WrapErrorForRemove(err, entity.ModelNameCatalogTrademark)
    }

    uc.eventBox.Emit(
        "%s::Remove: id=%d",
        entity.ModelNameCatalogTrademark,
        id,
    )

    return nil
}
