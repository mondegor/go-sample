package usecase

import (
	"context"

	"github.com/google/uuid"
	"github.com/mondegor/go-components/mrordering"
	"github.com/mondegor/go-sysmess/mrmsg"
	"github.com/mondegor/go-webcore/mrcore"
	"github.com/mondegor/go-webcore/mrenum"
	"github.com/mondegor/go-webcore/mrlog"
	"github.com/mondegor/go-webcore/mrsender"
	"github.com/mondegor/go-webcore/mrsender/decorator"
	"github.com/mondegor/go-webcore/mrstatus"
	"github.com/mondegor/go-webcore/mrstatus/mrflow"

	"github.com/mondegor/go-sample/internal/catalog/product/module"
	"github.com/mondegor/go-sample/internal/catalog/product/section/adm"
	"github.com/mondegor/go-sample/internal/catalog/product/section/adm/entity"
	"github.com/mondegor/go-sample/pkg/catalog/api"
)

type (
	// Product - comment struct.
	Product struct {
		storage      adm.ProductStorage
		categoryAPI  api.CategoryAvailability
		trademarkAPI api.TrademarkAvailability
		orderingAPI  mrordering.Mover
		eventEmitter mrsender.EventEmitter
		errorWrapper mrcore.UseCaseErrorWrapper
		statusFlow   mrstatus.Flow
	}
)

// NewProduct - создаёт объект Product.
func NewProduct(
	storage adm.ProductStorage,
	categoryAPI api.CategoryAvailability,
	trademarkAPI api.TrademarkAvailability,
	orderingAPI mrordering.Mover,
	eventEmitter mrsender.EventEmitter,
	errorWrapper mrcore.UseCaseErrorWrapper,
) *Product {
	return &Product{
		storage:      storage,
		categoryAPI:  categoryAPI,
		trademarkAPI: trademarkAPI,
		orderingAPI:  orderingAPI,
		eventEmitter: decorator.NewSourceEmitter(eventEmitter, entity.ModelNameProduct),
		errorWrapper: errorWrapper,
		statusFlow:   mrflow.ItemStatusFlow(),
	}
}

// GetList - comment method.
func (uc *Product) GetList(ctx context.Context, params entity.ProductParams) (items []entity.Product, countItems uint64, err error) {
	items, countItems, err = uc.storage.FetchWithTotal(ctx, params)
	if err != nil {
		return nil, 0, uc.errorWrapper.WrapErrorFailed(err, entity.ModelNameProduct)
	}

	if countItems == 0 {
		return make([]entity.Product, 0), 0, nil
	}

	return items, countItems, nil
}

// GetItem - comment method.
func (uc *Product) GetItem(ctx context.Context, itemID uint64) (entity.Product, error) {
	if itemID == 0 {
		return entity.Product{}, mrcore.ErrUseCaseEntityNotFound.New()
	}

	item, err := uc.storage.FetchOne(ctx, itemID)
	if err != nil {
		return entity.Product{}, uc.errorWrapper.WrapErrorEntityNotFoundOrFailed(err, entity.ModelNameProduct, itemID)
	}

	return item, nil
}

// Create - comment method.
func (uc *Product) Create(ctx context.Context, item entity.Product) (itemID uint64, err error) {
	if err = uc.checkItem(ctx, item); err != nil {
		return 0, err
	}

	item.Status = mrenum.ItemStatusDraft

	itemID, err = uc.storage.Insert(ctx, item)
	if err != nil {
		return 0, uc.errorWrapper.WrapErrorFailed(err, entity.ModelNameProduct)
	}

	uc.eventEmitter.Emit(ctx, "Create", mrmsg.Data{"id": itemID})

	if err = uc.orderingAPI.MoveToLast(ctx, itemID, uc.storage.NewCondition(item.CategoryID)); err != nil {
		mrlog.Ctx(ctx).Error().Err(err)
	}

	return itemID, nil
}

// Store - comment method.
func (uc *Product) Store(ctx context.Context, item entity.Product) error {
	if item.ID == 0 {
		return mrcore.ErrUseCaseEntityNotFound.New()
	}

	if item.TagVersion == 0 {
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

	uc.eventEmitter.Emit(ctx, "Store", mrmsg.Data{"id": item.ID, "ver": tagVersion})

	return nil
}

// ChangeStatus - comment method.
func (uc *Product) ChangeStatus(ctx context.Context, item entity.Product) error {
	if item.ID == 0 {
		return mrcore.ErrUseCaseEntityNotFound.New()
	}

	if item.TagVersion == 0 {
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

	uc.eventEmitter.Emit(ctx, "ChangeStatus", mrmsg.Data{"id": item.ID, "ver": tagVersion, "status": item.Status})

	return nil
}

// Remove - comment method.
func (uc *Product) Remove(ctx context.Context, itemID uint64) error {
	if itemID == 0 {
		return mrcore.ErrUseCaseEntityNotFound.New()
	}

	categoryID, err := uc.getCategoryID(ctx, itemID)
	if err != nil {
		return err
	}

	if err = uc.orderingAPI.Unlink(ctx, itemID, uc.storage.NewCondition(categoryID)); err != nil {
		return err
	}

	if err = uc.storage.Delete(ctx, itemID); err != nil {
		return uc.errorWrapper.WrapErrorEntityNotFoundOrFailed(err, entity.ModelNameProduct, itemID)
	}

	uc.eventEmitter.Emit(ctx, "Remove", mrmsg.Data{"id": itemID})

	return nil
}

// MoveAfterID - comment method.
func (uc *Product) MoveAfterID(ctx context.Context, itemID, afterID uint64) error {
	if itemID == 0 {
		return mrcore.ErrUseCaseEntityNotFound.New()
	}

	categoryID, err := uc.getCategoryID(ctx, itemID)
	if err != nil {
		return err
	}

	if err = uc.orderingAPI.MoveAfterID(ctx, itemID, afterID, uc.storage.NewCondition(categoryID)); err != nil {
		return err
	}

	uc.eventEmitter.Emit(ctx, "Move", mrmsg.Data{"id": itemID, "afterId": afterID})

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

func (uc *Product) getCategoryID(ctx context.Context, itemID uint64) (uuid.UUID, error) {
	item, err := uc.storage.FetchOne(ctx, itemID)
	if err != nil {
		return uuid.Nil, uc.errorWrapper.WrapErrorEntityNotFoundOrFailed(err, entity.ModelNameProduct, itemID)
	}

	if item.CategoryID == uuid.Nil {
		return uuid.Nil, mrcore.ErrInternal.New().WithAttr(entity.ModelNameProduct, mrmsg.Data{"categoryId": item.CategoryID})
	}

	return item.CategoryID, nil
}
