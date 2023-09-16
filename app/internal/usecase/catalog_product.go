package usecase

import (
    "context"
    "go-sample/internal/entity"

    "github.com/mondegor/go-components/mrcom"
    mrcom_orderer "github.com/mondegor/go-components/mrcom/orderer"
    "github.com/mondegor/go-storage/mrentity"
    "github.com/mondegor/go-sysmess/mrerr"
    "github.com/mondegor/go-webcore/mrcore"
    "github.com/mondegor/go-webcore/mrctx"
    "github.com/mondegor/go-webcore/mrtool"
)

type (
	CatalogProduct struct {
        componentOrderer mrcom_orderer.Component
        storage CatalogProductStorage
        storageCatalogTrademark CatalogTrademarkStorage
        eventBox mrcore.EventBox
        serviceHelper *mrtool.ServiceHelper
        statusFlow mrcom.ItemStatusFlow
    }
)

func NewCatalogProduct(componentOrderer mrcom_orderer.Component,
                       storage CatalogProductStorage,
                       storageCatalogTrademark CatalogTrademarkStorage,
                       eventBox mrcore.EventBox,
                       serviceHelper *mrtool.ServiceHelper) *CatalogProduct {
    return &CatalogProduct{
        componentOrderer: componentOrderer,
        storage: storage,
        storageCatalogTrademark: storageCatalogTrademark,
        eventBox: eventBox,
        serviceHelper: serviceHelper,
        statusFlow: mrcom.ItemStatusFlowDefault,
    }
}

func (uc *CatalogProduct) GetList(ctx context.Context, listFilter *entity.CatalogProductListFilter) ([]entity.CatalogProduct, error) {
    items := make([]entity.CatalogProduct, 0, 16)
    err := uc.storage.LoadAll(ctx, listFilter, &items)

    if err != nil {
        return nil, mrcore.FactoryErrServiceEntityTemporarilyUnavailable.Wrap(err, entity.ModelNameCatalogProduct)
    }

    return items, nil
}

func (uc *CatalogProduct) GetItem(ctx context.Context, id mrentity.KeyInt32, categoryId mrentity.KeyInt32) (*entity.CatalogProduct, error) {
    if id < 1 {
        return nil, mrcore.FactoryErrServiceIncorrectInputData.New(mrerr.Arg{"id": id})
    }

    item := &entity.CatalogProduct{Id: id, CategoryId: categoryId}
    err := uc.storage.LoadOne(ctx, item)

    if err != nil {
        return nil, uc.serviceHelper.WrapErrorForSelect(err, entity.ModelNameCatalogProduct)
    }

    return item, nil
}

// Create
// modifies: item{Id}
func (uc *CatalogProduct) Create(ctx context.Context, item *entity.CatalogProduct) error {
    err := uc.checkArticle(ctx, item)

    if err != nil {
        return err
    }

    err = uc.storageCatalogTrademark.IsExists(ctx, item.TrademarkId)

    if err != nil {
        if mrcore.FactoryErrStorageNoRowFound.Is(err) {
            return FactoryErrCatalogTrademarkNotFound.Wrap(err, item.TrademarkId)
        }

        return err
    }

    item.Status = mrcom.ItemStatusDraft
    err = uc.storage.Insert(ctx, item)

    if err != nil {
        return mrcore.FactoryErrServiceEntityNotCreated.Wrap(err, entity.ModelNameCatalogProduct)
    }

    uc.eventBox.Emit(
        "%s::Create: id=%d",
        entity.ModelNameCatalogProduct,
        item.Id,
    )

    meta := uc.storage.GetMetaData(item.CategoryId)
    component := uc.componentOrderer.WithMetaData(meta)

    err = component.MoveToLast(
        ctx,
        item.Id,
    )

    if err != nil {
        mrctx.Logger(ctx).Err(err)
    }

    return nil
}

func (uc *CatalogProduct) Store(ctx context.Context, item *entity.CatalogProduct) error {
    if item.Id < 1 || item.Version < 1 {
        return mrcore.FactoryErrServiceIncorrectInputData.New(mrerr.Arg{"item.Id": item.Id, "Item.Version": item.Version})
    }

    err := uc.checkArticle(ctx, item)

    if err != nil {
        return err
    }

    err = uc.storage.Update(ctx, item)

    if err != nil {
        return uc.serviceHelper.WrapErrorForUpdate(err, entity.ModelNameCatalogProduct)
    }

    uc.eventBox.Emit(
        "%s::Store: id=%d",
        entity.ModelNameCatalogProduct,
        item.Id,
    )

    return nil
}

func (uc *CatalogProduct) ChangeStatus(ctx context.Context, item *entity.CatalogProduct) error {
    if item.Id < 1 || item.Version < 1 {
        return mrcore.FactoryErrServiceIncorrectInputData.New(mrerr.Arg{"item.Id": item.Id, "Item.Version": item.Version})
    }

    currentStatus, err := uc.storage.FetchStatus(ctx, item)

    if err != nil {
        return uc.serviceHelper.WrapErrorForSelect(err, entity.ModelNameCatalogProduct)
    }

    if !uc.statusFlow.Check(currentStatus, item.Status) {
        return mrcore.FactoryErrServiceIncorrectSwitchStatus.New(currentStatus, item.Status, entity.ModelNameCatalogProduct, item.Id)
    }

    err = uc.storage.UpdateStatus(ctx, item)

    if err != nil {
        return uc.serviceHelper.WrapErrorForUpdate(err, entity.ModelNameCatalogProduct)
    }

    uc.eventBox.Emit(
        "%s::ChangeStatus: id=%d, status=%s",
        entity.ModelNameCatalogProduct,
        item.Id,
        item.Status,
    )

    return nil
}

func (uc *CatalogProduct) Remove(ctx context.Context, id mrentity.KeyInt32) error {
    if id < 1 {
        return mrcore.FactoryErrServiceIncorrectInputData.New(mrerr.Arg{"id": id})
    }

    err := uc.storage.Delete(ctx, id)

    if err != nil {
        return uc.serviceHelper.WrapErrorForRemove(err, entity.ModelNameCatalogProduct)
    }

    uc.eventBox.Emit(
        "%s::Remove: id=%d",
        entity.ModelNameCatalogProduct,
        id,
    )

    return nil
}

func (uc *CatalogProduct) MoveAfterId(ctx context.Context, id mrentity.KeyInt32, afterId mrentity.KeyInt32, categoryId mrentity.KeyInt32) error {
    if categoryId < 1 {
        return mrcore.FactoryErrServiceIncorrectInputData.New(mrerr.Arg{"categoryId": categoryId})
    }

    meta := uc.storage.GetMetaData(categoryId)
    component := uc.componentOrderer.WithMetaData(meta)

    return component.MoveAfterId(ctx, id, afterId)
}

func (uc *CatalogProduct) checkArticle(ctx context.Context, item *entity.CatalogProduct) error {
    id, err := uc.storage.FetchIdByArticle(ctx, item.Article)

    if err != nil {
        if mrcore.FactoryErrStorageNoRowFound.Is(err) {
            return nil
        }

        return mrcore.FactoryErrServiceEntityTemporarilyUnavailable.Wrap(err, entity.ModelNameCatalogProduct)
    }

    if item.Id == id {
        return nil
    }

    return FactoryErrCatalogProductArticleAlreadyExists.New(item.Article)
}
