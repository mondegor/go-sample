package usecase

import (
	"context"

	"go-sample/internal/modules/catalog/product/module"

	entity "go-sample/internal/modules/catalog/product/entity/admin_api"
	"go-sample/pkg/modules/catalog/api"

	"github.com/google/uuid"
	"github.com/mondegor/go-components/mrsort"
	"github.com/mondegor/go-sysmess/mrmsg"
	"github.com/mondegor/go-webcore/mrcore"
	"github.com/mondegor/go-webcore/mrenum"
	"github.com/mondegor/go-webcore/mrlog"
	"github.com/mondegor/go-webcore/mrsender"
	"github.com/mondegor/go-webcore/mrstatus"
	"github.com/mondegor/go-webcore/mrstatus/mrflow"
	"github.com/mondegor/go-webcore/mrtype"
)

type (
	// Product - comment struct.
	Product struct {
		storage      ProductStorage
		categoryAPI  api.CategoryAPI
		trademarkAPI api.TrademarkAPI
		ordererAPI   mrsort.Orderer
		eventEmitter mrsender.EventEmitter
		errorWrapper mrcore.UsecaseErrorWrapper
		statusFlow   mrstatus.Flow
	}
)

// NewProduct - comment func.
func NewProduct(
	storage ProductStorage,
	categoryAPI api.CategoryAPI,
	trademarkAPI api.TrademarkAPI,
	ordererAPI mrsort.Orderer,
	eventEmitter mrsender.EventEmitter,
	errorWrapper mrcore.UsecaseErrorWrapper,
) *Product {
	return &Product{
		storage:      storage,
		categoryAPI:  categoryAPI,
		trademarkAPI: trademarkAPI,
		ordererAPI:   ordererAPI,
		eventEmitter: eventEmitter,
		errorWrapper: errorWrapper,
		statusFlow:   mrflow.ItemStatusFlow(),
	}
}

// GetList - comment method.
func (uc *Product) GetList(ctx context.Context, params entity.ProductParams) ([]entity.Product, int64, error) {
	fetchParams := uc.storage.NewSelectParams(params)

	total, err := uc.storage.FetchTotal(ctx, fetchParams.Where)
	if err != nil {
		return nil, 0, uc.errorWrapper.WrapErrorFailed(err, entity.ModelNameProduct)
	}

	if total < 1 {
		return nil, 0, nil
	}

	items, err := uc.storage.Fetch(ctx, fetchParams)
	if err != nil {
		return nil, 0, uc.errorWrapper.WrapErrorFailed(err, entity.ModelNameProduct)
	}

	return items, total, nil
}

// GetItem - comment method.
func (uc *Product) GetItem(ctx context.Context, itemID mrtype.KeyInt32) (entity.Product, error) {
	if itemID < 1 {
		return entity.Product{}, mrcore.ErrUseCaseEntityNotFound.New()
	}

	item, err := uc.storage.FetchOne(ctx, itemID)
	if err != nil {
		return entity.Product{}, uc.errorWrapper.WrapErrorEntityNotFoundOrFailed(err, entity.ModelNameProduct, itemID)
	}

	return item, nil
}

// Create - comment method.
func (uc *Product) Create(ctx context.Context, item entity.Product) (mrtype.KeyInt32, error) {
	if err := uc.checkItem(ctx, item); err != nil {
		return 0, err
	}

	item.Status = mrenum.ItemStatusDraft

	itemID, err := uc.storage.Insert(ctx, item)
	if err != nil {
		return 0, uc.errorWrapper.WrapErrorFailed(err, entity.ModelNameProduct)
	}

	uc.emitEvent(ctx, "Create", mrmsg.Data{"id": itemID})

	if err := uc.getOrdererAPI(item.CategoryID).MoveToLast(ctx, itemID); err != nil {
		mrlog.Ctx(ctx).Error().Err(err)
	}

	return itemID, nil
}

// Store - comment method.
func (uc *Product) Store(ctx context.Context, item entity.Product) error {
	if item.ID < 1 {
		return mrcore.ErrUseCaseEntityNotFound.New()
	}

	if item.TagVersion < 1 {
		return mrcore.ErrUseCaseEntityVersionInvalid.New()
	}

	// предварительная проверка существования записи нужна для того,
	// чтобы при Update быть уверенным, что отсутствие записи из-за ошибки VersionInvalid
	if _, err := uc.storage.FetchStatus(ctx, item.ID); err != nil {
		return uc.errorWrapper.WrapErrorEntityNotFoundOrFailed(err, entity.ModelNameProduct, item.ID)
	}

	if err := uc.checkItem(ctx, item); err != nil {
		return err
	}

	tagVersion, err := uc.storage.Update(ctx, item)
	if err != nil {
		if uc.errorWrapper.IsNotFoundError(err) {
			return mrcore.ErrUseCaseEntityVersionInvalid.Wrap(err)
		}

		return uc.errorWrapper.WrapErrorFailed(err, entity.ModelNameProduct)
	}

	uc.emitEvent(ctx, "Store", mrmsg.Data{"id": item.ID, "ver": tagVersion})

	return nil
}

// ChangeStatus - comment method.
func (uc *Product) ChangeStatus(ctx context.Context, item entity.Product) error {
	if item.ID < 1 {
		return mrcore.ErrUseCaseEntityNotFound.New()
	}

	if item.TagVersion < 1 {
		return mrcore.ErrUseCaseEntityVersionInvalid.New()
	}

	currentStatus, err := uc.storage.FetchStatus(ctx, item.ID)
	if err != nil {
		return uc.errorWrapper.WrapErrorEntityNotFoundOrFailed(err, entity.ModelNameProduct, item.ID)
	}

	if currentStatus == item.Status {
		return nil
	}

	if !uc.statusFlow.Check(currentStatus, item.Status) {
		return mrcore.ErrUseCaseSwitchStatusRejected.New(currentStatus, item.Status)
	}

	tagVersion, err := uc.storage.UpdateStatus(ctx, item)
	if err != nil {
		if uc.errorWrapper.IsNotFoundError(err) {
			return mrcore.ErrUseCaseEntityVersionInvalid.Wrap(err)
		}

		return uc.errorWrapper.WrapErrorFailed(err, entity.ModelNameProduct)
	}

	uc.emitEvent(ctx, "ChangeStatus", mrmsg.Data{"id": item.ID, "ver": tagVersion, "status": item.Status})

	return nil
}

// Remove - comment method.
func (uc *Product) Remove(ctx context.Context, itemID mrtype.KeyInt32) error {
	if itemID < 1 {
		return mrcore.ErrUseCaseEntityNotFound.New()
	}

	categoryID, err := uc.getCategoryID(ctx, itemID)
	if err != nil {
		return err
	}

	if err = uc.getOrdererAPI(categoryID).Unlink(ctx, itemID); err != nil {
		return err
	}

	if err = uc.storage.Delete(ctx, itemID); err != nil {
		return uc.errorWrapper.WrapErrorEntityNotFoundOrFailed(err, entity.ModelNameProduct, itemID)
	}

	uc.emitEvent(ctx, "Remove", mrmsg.Data{"id": itemID})

	return nil
}

// MoveAfterID - comment method.
func (uc *Product) MoveAfterID(ctx context.Context, itemID, afterID mrtype.KeyInt32) error {
	if itemID < 1 {
		return mrcore.ErrUseCaseEntityNotFound.New()
	}

	categoryID, err := uc.getCategoryID(ctx, itemID)
	if err != nil {
		return err
	}

	if err = uc.getOrdererAPI(categoryID).MoveAfterID(ctx, itemID, afterID); err != nil {
		return err
	}

	uc.emitEvent(ctx, "Move", mrmsg.Data{"id": itemID, "afterId": afterID})

	return nil
}

func (uc *Product) checkItem(ctx context.Context, item entity.Product) error {
	if err := uc.checkArticle(ctx, item); err != nil {
		return err
	}

	if item.ID == 0 || item.CategoryID != uuid.Nil {
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
	itemID, err := uc.storage.FetchIDByArticle(ctx, item.Article)
	if err != nil {
		if uc.errorWrapper.IsNotFoundError(err) {
			return nil
		}

		return uc.errorWrapper.WrapErrorFailed(err, entity.ModelNameProduct)
	}

	if item.ID != itemID {
		return module.ErrUseCaseProductArticleAlreadyExists.New(item.Article)
	}

	return nil
}

func (uc *Product) getOrdererAPI(categoryID uuid.UUID) mrsort.Orderer {
	meta := uc.storage.NewOrderMeta(categoryID)

	return uc.ordererAPI.WithMetaData(meta)
}

func (uc *Product) getCategoryID(ctx context.Context, itemID mrtype.KeyInt32) (uuid.UUID, error) {
	item, err := uc.storage.FetchOne(ctx, itemID)
	if err != nil {
		return uuid.Nil, uc.errorWrapper.WrapErrorEntityNotFoundOrFailed(err, entity.ModelNameProduct, itemID)
	}

	if item.CategoryID == uuid.Nil {
		return uuid.Nil, mrcore.ErrInternal.New().WithAttr(entity.ModelNameProduct, mrmsg.Data{"categoryId": item.CategoryID})
	}

	return item.CategoryID, nil
}

func (uc *Product) emitEvent(ctx context.Context, eventName string, data mrmsg.Data) {
	uc.eventEmitter.EmitWithSource(
		ctx,
		eventName,
		entity.ModelNameProduct,
		data,
	)
}
