package usecase

import (
	"context"
	"go-sample/internal/modules/catalog/entity/admin-api"
	usecase_shared "go-sample/internal/modules/catalog/usecase/shared"

	"github.com/mondegor/go-components/mrorderer"
	"github.com/mondegor/go-sysmess/mrerr"
	"github.com/mondegor/go-webcore/mrcore"
	"github.com/mondegor/go-webcore/mrctx"
	"github.com/mondegor/go-webcore/mrenum"
	"github.com/mondegor/go-webcore/mrtool"
	"github.com/mondegor/go-webcore/mrtype"
)

type (
	Product struct {
		componentOrderer mrorderer.Component
		storage          ProductStorage
		storageCategory  CategoryStorage
		storageTrademark TrademarkStorage
		eventBox         mrcore.EventBox
		serviceHelper    *mrtool.ServiceHelper
		statusFlow       mrenum.StatusFlow
	}
)

func NewProduct(
	componentOrderer mrorderer.Component,
	storage ProductStorage,
	storageCategory CategoryStorage,
	storageTrademark TrademarkStorage,
	eventBox mrcore.EventBox,
	serviceHelper *mrtool.ServiceHelper,
) *Product {
	return &Product{
		componentOrderer: componentOrderer,
		storage:          storage,
		storageCategory:  storageCategory,
		storageTrademark: storageTrademark,
		eventBox:         eventBox,
		serviceHelper:    serviceHelper,
		statusFlow:       mrenum.ItemStatusFlow,
	}
}

func (uc *Product) GetList(ctx context.Context, params entity.ProductParams) ([]entity.Product, int64, error) {
	fetchParams := uc.storage.NewFetchParams(params)
	total, err := uc.storage.FetchTotal(ctx, fetchParams.Where)

	if err != nil {
		return nil, 0, mrcore.FactoryErrServiceTemporarilyUnavailable.Wrap(err, entity.ModelNameCatalogProduct)
	}

	if total < 1 {
		return []entity.Product{}, 0, nil
	}

	items, err := uc.storage.Fetch(ctx, fetchParams)

	if err != nil {
		return nil, 0, mrcore.FactoryErrServiceTemporarilyUnavailable.Wrap(err, entity.ModelNameCatalogProduct)
	}

	return items, total, nil
}

func (uc *Product) GetItem(ctx context.Context, id mrtype.KeyInt32) (*entity.Product, error) {
	if id < 1 {
		return nil, mrcore.FactoryErrServiceIncorrectInputData.New(mrerr.Arg{"id": id})
	}

	item := &entity.Product{ID: id}

	if err := uc.storage.LoadOne(ctx, item); err != nil {
		return nil, uc.serviceHelper.WrapErrorForSelect(err, entity.ModelNameCatalogProduct)
	}

	return item, nil
}

// Create
// modifies: item{ID}
func (uc *Product) Create(ctx context.Context, item *entity.Product) error {
	if err := uc.checkProduct(ctx, item); err != nil {
		return err
	}

	item.Status = mrenum.ItemStatusDraft

	if err := uc.storage.Insert(ctx, item); err != nil {
		return mrcore.FactoryErrServiceEntityNotCreated.Wrap(err, entity.ModelNameCatalogProduct)
	}

	uc.eventBox.Emit(
		"%s::Create: id=%d",
		entity.ModelNameCatalogProduct,
		item.ID,
	)

	meta := uc.storage.GetMetaData(item.CategoryID)
	component := uc.componentOrderer.WithMetaData(meta)

	if err := component.MoveToLast(ctx, item.ID); err != nil {
		mrctx.Logger(ctx).Err(err)
	}

	return nil
}

func (uc *Product) Store(ctx context.Context, item *entity.Product) error {
	if item.ID < 1 || item.TagVersion < 1 {
		return mrcore.FactoryErrServiceIncorrectInputData.New(mrerr.Arg{"item.id": item.ID, "version": item.TagVersion})
	}

	if err := uc.checkProduct(ctx, item); err != nil {
		return err
	}

	if err := uc.storage.Update(ctx, item); err != nil {
		return uc.serviceHelper.WrapErrorForUpdate(err, entity.ModelNameCatalogProduct)
	}

	uc.eventBox.Emit(
		"%s::Store: id=%d",
		entity.ModelNameCatalogProduct,
		item.ID,
	)

	return nil
}

func (uc *Product) ChangeStatus(ctx context.Context, item *entity.Product) error {
	if item.ID < 1 || item.TagVersion < 1 {
		return mrcore.FactoryErrServiceIncorrectInputData.New(mrerr.Arg{"item.id": item.ID, "version": item.TagVersion})
	}

	currentStatus, err := uc.storage.FetchStatus(ctx, item)

	if err != nil {
		return uc.serviceHelper.WrapErrorForSelect(err, entity.ModelNameCatalogProduct)
	}

	if currentStatus == item.Status {
		return nil
	}

	if !uc.statusFlow.Check(currentStatus, item.Status) {
		return mrcore.FactoryErrServiceIncorrectSwitchStatus.New(currentStatus, item.Status, entity.ModelNameCatalogProduct, item.ID)
	}

	err = uc.storage.UpdateStatus(ctx, item)

	if err != nil {
		return uc.serviceHelper.WrapErrorForUpdate(err, entity.ModelNameCatalogProduct)
	}

	uc.eventBox.Emit(
		"%s::ChangeStatus: id=%d, status=%s",
		entity.ModelNameCatalogProduct,
		item.ID,
		item.Status,
	)

	return nil
}

func (uc *Product) Remove(ctx context.Context, id mrtype.KeyInt32) error {
	if id < 1 {
		return mrcore.FactoryErrServiceIncorrectInputData.New(mrerr.Arg{"id": id})
	}

	if err := uc.storage.Delete(ctx, id); err != nil {
		return uc.serviceHelper.WrapErrorForRemove(err, entity.ModelNameCatalogProduct)
	}

	uc.eventBox.Emit(
		"%s::Remove: id=%d",
		entity.ModelNameCatalogProduct,
		id,
	)

	return nil
}

func (uc *Product) MoveAfterID(ctx context.Context, id mrtype.KeyInt32, afterID mrtype.KeyInt32) error {
	item := entity.Product{
		ID: id,
	}

	if err := uc.storage.LoadOne(ctx, &item); err != nil {
		return err
	}

	if item.CategoryID < 1 {
		return mrcore.FactoryErrServiceIncorrectInputData.New(mrerr.Arg{"categoryId": item.CategoryID})
	}

	meta := uc.storage.GetMetaData(item.CategoryID)
	component := uc.componentOrderer.WithMetaData(meta)

	return component.MoveAfterID(ctx, id, afterID)
}

func (uc *Product) checkProduct(ctx context.Context, item *entity.Product) error {
	if err := uc.checkArticle(ctx, item); err != nil {
		return err
	}

	if err := uc.storageCategory.IsExists(ctx, item.CategoryID); err != nil {
		if mrcore.FactoryErrStorageNoRowFound.Is(err) {
			return usecase_shared.FactoryErrCategoryNotFound.Wrap(err, item.CategoryID)
		}

		return err
	}

	if err := uc.storageTrademark.IsExists(ctx, item.TrademarkID); err != nil {
		if mrcore.FactoryErrStorageNoRowFound.Is(err) {
			return usecase_shared.FactoryErrTrademarkNotFound.Wrap(err, item.TrademarkID)
		}

		return err
	}

	return nil
}

func (uc *Product) checkArticle(ctx context.Context, item *entity.Product) error {
	id, err := uc.storage.FetchIdByArticle(ctx, item.Article)

	if err != nil {
		if mrcore.FactoryErrStorageNoRowFound.Is(err) {
			return nil
		}

		return mrcore.FactoryErrServiceTemporarilyUnavailable.Wrap(err, entity.ModelNameCatalogProduct)
	}

	if item.ID != id {
		return usecase_shared.FactoryErrProductArticleAlreadyExists.New(item.Article)
	}

	return nil
}
