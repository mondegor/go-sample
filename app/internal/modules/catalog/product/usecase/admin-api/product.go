package usecase

import (
	"context"
	"go-sample/internal/modules/catalog/product/entity/admin-api"
	usecase_shared "go-sample/internal/modules/catalog/product/usecase/shared"
	"go-sample/pkg/modules/catalog"

	"github.com/mondegor/go-components/mrorderer"
	"github.com/mondegor/go-sysmess/mrmsg"
	"github.com/mondegor/go-webcore/mrcore"
	"github.com/mondegor/go-webcore/mrenum"
	"github.com/mondegor/go-webcore/mrlog"
	"github.com/mondegor/go-webcore/mrsender"
	"github.com/mondegor/go-webcore/mrtype"
)

type (
	Product struct {
		storage       ProductStorage
		categoryAPI   catalog.CategoryAPI
		trademarkAPI  catalog.TrademarkAPI
		ordererAPI    mrorderer.API
		eventEmitter  mrsender.EventEmitter
		usecaseHelper *mrcore.UsecaseHelper
		statusFlow    mrenum.StatusFlow
	}
)

func NewProduct(
	storage ProductStorage,
	categoryAPI catalog.CategoryAPI,
	trademarkAPI catalog.TrademarkAPI,
	ordererAPI mrorderer.API,
	eventEmitter mrsender.EventEmitter,
	usecaseHelper *mrcore.UsecaseHelper,
) *Product {
	return &Product{
		storage:       storage,
		categoryAPI:   categoryAPI,
		trademarkAPI:  trademarkAPI,
		ordererAPI:    ordererAPI,
		eventEmitter:  eventEmitter,
		usecaseHelper: usecaseHelper,
		statusFlow:    mrenum.ItemStatusFlow,
	}
}

func (uc *Product) GetList(ctx context.Context, params entity.ProductParams) ([]entity.Product, int64, error) {
	fetchParams := uc.storage.NewFetchParams(params)
	total, err := uc.storage.FetchTotal(ctx, fetchParams.Where)

	if err != nil {
		return nil, 0, uc.usecaseHelper.WrapErrorFailed(err, entity.ModelNameProduct)
	}

	if total < 1 {
		return []entity.Product{}, 0, nil
	}

	items, err := uc.storage.Fetch(ctx, fetchParams)

	if err != nil {
		return nil, 0, uc.usecaseHelper.WrapErrorFailed(err, entity.ModelNameProduct)
	}

	return items, total, nil
}

func (uc *Product) GetItem(ctx context.Context, itemID mrtype.KeyInt32) (entity.Product, error) {
	if itemID < 1 {
		return entity.Product{}, mrcore.FactoryErrUseCaseEntityNotFound.New()
	}

	item, err := uc.storage.FetchOne(ctx, itemID)

	if err != nil {
		return entity.Product{}, uc.usecaseHelper.WrapErrorEntityNotFoundOrFailed(err, entity.ModelNameProduct, itemID)
	}

	return item, nil
}

func (uc *Product) Create(ctx context.Context, item entity.Product) (mrtype.KeyInt32, error) {
	if err := uc.checkItem(ctx, item); err != nil {
		return 0, err
	}

	item.Status = mrenum.ItemStatusDraft
	itemID, err := uc.storage.Insert(ctx, item)

	if err != nil {
		return 0, uc.usecaseHelper.WrapErrorFailed(err, entity.ModelNameProduct)
	}

	uc.emitEvent(ctx, "Create", mrmsg.Data{"id": itemID})

	meta := uc.storage.GetMetaData(item.CategoryID)
	ordererAPI := uc.ordererAPI.WithMetaData(meta)

	if err := ordererAPI.MoveToLast(ctx, itemID); err != nil {
		mrlog.Ctx(ctx).Error().Err(err)
	}

	return itemID, nil
}

func (uc *Product) Store(ctx context.Context, item entity.Product) error {
	if item.ID < 1 {
		return mrcore.FactoryErrUseCaseEntityNotFound.New()
	}

	if item.TagVersion < 1 {
		return mrcore.FactoryErrUseCaseEntityVersionInvalid.New()
	}

	if err := uc.storage.IsExists(ctx, item.ID); err != nil {
		return uc.usecaseHelper.WrapErrorEntityNotFoundOrFailed(err, entity.ModelNameProduct, item.ID)
	}

	if err := uc.checkItem(ctx, item); err != nil {
		return err
	}

	tagVersion, err := uc.storage.Update(ctx, item)

	if err != nil {
		if uc.usecaseHelper.IsNotFoundError(err) {
			return mrcore.FactoryErrUseCaseEntityVersionInvalid.Wrap(err)
		}

		return uc.usecaseHelper.WrapErrorFailed(err, entity.ModelNameProduct)
	}

	uc.emitEvent(ctx, "Store", mrmsg.Data{"id": item.ID, "ver": tagVersion})

	return nil
}

func (uc *Product) ChangeStatus(ctx context.Context, item entity.Product) error {
	if item.ID < 1 {
		return mrcore.FactoryErrUseCaseEntityNotFound.New()
	}

	if item.TagVersion < 1 {
		return mrcore.FactoryErrUseCaseEntityVersionInvalid.New()
	}

	currentStatus, err := uc.storage.FetchStatus(ctx, item)

	if err != nil {
		return uc.usecaseHelper.WrapErrorEntityNotFoundOrFailed(err, entity.ModelNameProduct, item.ID)
	}

	if currentStatus == item.Status {
		return nil
	}

	if !uc.statusFlow.Check(currentStatus, item.Status) {
		return mrcore.FactoryErrUseCaseSwitchStatusRejected.New(currentStatus, item.Status)
	}

	tagVersion, err := uc.storage.UpdateStatus(ctx, item)

	if err != nil {
		if uc.usecaseHelper.IsNotFoundError(err) {
			return mrcore.FactoryErrUseCaseEntityVersionInvalid.Wrap(err)
		}

		return uc.usecaseHelper.WrapErrorFailed(err, entity.ModelNameProduct)
	}

	uc.emitEvent(ctx, "ChangeStatus", mrmsg.Data{"id": item.ID, "ver": tagVersion, "status": item.Status})

	return nil
}

func (uc *Product) Remove(ctx context.Context, itemID mrtype.KeyInt32) error {
	if itemID < 1 {
		return mrcore.FactoryErrUseCaseEntityNotFound.New()
	}

	if err := uc.storage.Delete(ctx, itemID); err != nil {
		return uc.usecaseHelper.WrapErrorEntityNotFoundOrFailed(err, entity.ModelNameProduct, itemID)
	}

	uc.emitEvent(ctx, "Remove", mrmsg.Data{"id": itemID})

	return nil
}

func (uc *Product) MoveAfterID(ctx context.Context, itemID mrtype.KeyInt32, afterID mrtype.KeyInt32) error {
	if itemID < 1 {
		return mrcore.FactoryErrUseCaseEntityNotFound.New()
	}

	item, err := uc.storage.FetchOne(ctx, itemID)

	if err != nil {
		return uc.usecaseHelper.WrapErrorEntityNotFoundOrFailed(err, entity.ModelNameProduct, itemID)
	}

	if item.CategoryID < 1 {
		return mrcore.FactoryErrInternal.WithAttr(entity.ModelNameProduct, mrmsg.Data{"categoryId": item.CategoryID}).New()
	}

	meta := uc.storage.GetMetaData(item.CategoryID)
	ordererAPI := uc.ordererAPI.WithMetaData(meta)

	if err := ordererAPI.MoveAfterID(ctx, itemID, afterID); err != nil {
		return err
	}

	uc.emitEvent(ctx, "Move", mrmsg.Data{"id": itemID, "afterId": afterID})

	return nil
}

func (uc *Product) checkItem(ctx context.Context, item entity.Product) error {
	if err := uc.checkArticle(ctx, item); err != nil {
		return err
	}

	if item.ID == 0 || item.CategoryID > 0 {
		if err := uc.categoryAPI.CheckingAvailability(ctx, item.CategoryID); err != nil {
			return err
		}
	}

	if item.ID == 0 || item.TrademarkID > 0 {
		if err := uc.trademarkAPI.CheckingAvailability(ctx, item.TrademarkID); err != nil {
			return err
		}
	}

	return nil
}

func (uc *Product) checkArticle(ctx context.Context, item entity.Product) error {
	itemID, err := uc.storage.FetchIdByArticle(ctx, item.Article)

	if err != nil {
		if uc.usecaseHelper.IsNotFoundError(err) {
			return nil
		}

		return uc.usecaseHelper.WrapErrorFailed(err, entity.ModelNameProduct)
	}

	if item.ID != itemID {
		return usecase_shared.FactoryErrProductArticleAlreadyExists.New(item.Article)
	}

	return nil
}

func (uc *Product) emitEvent(ctx context.Context, eventName string, data mrmsg.Data) {
	uc.eventEmitter.EmitWithSource(
		ctx,
		eventName,
		entity.ModelNameProduct,
		data,
	)
}
