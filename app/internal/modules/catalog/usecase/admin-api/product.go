package usecase

import (
	"context"
	"go-sample/internal/modules/catalog/entity/admin-api"
	usecase_shared "go-sample/internal/modules/catalog/usecase/shared"

	"github.com/mondegor/go-components/mrorderer"
	"github.com/mondegor/go-sysmess/mrmsg"
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
		categoryAPI      usecase_shared.CategoryServiceAPI
		trademarkAPI     usecase_shared.TrademarkServiceAPI
		eventBox         mrcore.EventBox
		serviceHelper    *mrtool.ServiceHelper
		statusFlow       mrenum.StatusFlow
	}
)

func NewProduct(
	componentOrderer mrorderer.Component,
	storage ProductStorage,
	categoryAPI usecase_shared.CategoryServiceAPI,
	trademarkAPI usecase_shared.TrademarkServiceAPI,
	eventBox mrcore.EventBox,
	serviceHelper *mrtool.ServiceHelper,
) *Product {
	return &Product{
		componentOrderer: componentOrderer,
		storage:          storage,
		categoryAPI:      categoryAPI,
		trademarkAPI:     trademarkAPI,
		eventBox:         eventBox,
		serviceHelper:    serviceHelper,
		statusFlow:       mrenum.ItemStatusFlow,
	}
}

func (uc *Product) GetList(ctx context.Context, params entity.ProductParams) ([]entity.Product, int64, error) {
	fetchParams := uc.storage.NewFetchParams(params)
	total, err := uc.storage.FetchTotal(ctx, fetchParams.Where)

	if err != nil {
		return nil, 0, uc.serviceHelper.WrapErrorFailed(err, entity.ModelNameProduct)
	}

	if total < 1 {
		return []entity.Product{}, 0, nil
	}

	items, err := uc.storage.Fetch(ctx, fetchParams)

	if err != nil {
		return nil, 0, uc.serviceHelper.WrapErrorFailed(err, entity.ModelNameProduct)
	}

	return items, total, nil
}

func (uc *Product) GetItem(ctx context.Context, id mrtype.KeyInt32) (*entity.Product, error) {
	if id < 1 {
		return nil, mrcore.FactoryErrServiceEntityNotFound.New()
	}

	item := &entity.Product{ID: id}

	if err := uc.storage.LoadOne(ctx, item); err != nil {
		return nil, uc.serviceHelper.WrapErrorEntityNotFoundOrFailed(err, entity.ModelNameProduct, id)
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
		return uc.serviceHelper.WrapErrorFailed(err, entity.ModelNameProduct)
	}

	uc.eventBoxEmitEntity(ctx, "Create", mrmsg.Data{"id": item.ID})

	meta := uc.storage.GetMetaData(item.CategoryID)
	component := uc.componentOrderer.WithMetaData(meta)

	if err := component.MoveToLast(ctx, item.ID); err != nil {
		mrctx.Logger(ctx).Err(err)
	}

	return nil
}

func (uc *Product) Store(ctx context.Context, item *entity.Product) error {
	if item.ID < 1 {
		return mrcore.FactoryErrServiceEntityNotFound.New()
	}

	if item.TagVersion < 1 {
		return mrcore.FactoryErrServiceEntityVersionInvalid.New()
	}

	if err := uc.storage.IsExists(ctx, item.ID); err != nil {
		return uc.serviceHelper.WrapErrorEntityNotFoundOrFailed(err, entity.ModelNameProduct, item.ID)
	}

	if err := uc.checkProduct(ctx, item); err != nil {
		return err
	}

	version, err := uc.storage.Update(ctx, item)

	if err != nil {
		return uc.serviceHelper.WrapErrorEntity(
			mrcore.FactoryErrServiceEntityVersionInvalid,
			err,
			entity.ModelNameProduct,
			mrmsg.Data{"id": item.ID, "ver": item.TagVersion},
		)
	}

	uc.eventBoxEmitEntity(ctx, "Store", mrmsg.Data{"id": item.ID, "ver": version})

	return nil
}

func (uc *Product) ChangeStatus(ctx context.Context, item *entity.Product) error {
	if item.ID < 1 {
		return mrcore.FactoryErrServiceEntityNotFound.New()
	}

	if item.TagVersion < 1 {
		return mrcore.FactoryErrServiceEntityVersionInvalid.New()
	}

	currentStatus, err := uc.storage.FetchStatus(ctx, item)

	if err != nil {
		return uc.serviceHelper.WrapErrorEntityNotFoundOrFailed(err, entity.ModelNameProduct, item.ID)
	}

	if currentStatus == item.Status {
		return nil
	}

	if !uc.statusFlow.Check(currentStatus, item.Status) {
		return mrcore.FactoryErrServiceSwitchStatusRejected.New(currentStatus, item.Status)
	}

	version, err := uc.storage.UpdateStatus(ctx, item)

	if err != nil {
		return uc.serviceHelper.WrapErrorEntity(
			mrcore.FactoryErrServiceEntityVersionInvalid,
			err,
			entity.ModelNameProduct,
			mrmsg.Data{"id": item.ID, "ver": item.TagVersion},
		)
	}

	uc.eventBoxEmitEntity(ctx, "ChangeStatus", mrmsg.Data{"id": item.ID, "ver": version, "status": item.Status})

	return nil
}

func (uc *Product) Remove(ctx context.Context, id mrtype.KeyInt32) error {
	if id < 1 {
		return mrcore.FactoryErrServiceEntityNotFound.New()
	}

	if err := uc.storage.Delete(ctx, id); err != nil {
		return uc.serviceHelper.WrapErrorEntityNotFoundOrFailed(err, entity.ModelNameProduct, id)
	}

	uc.eventBoxEmitEntity(ctx, "Remove", mrmsg.Data{"id": id})

	return nil
}

func (uc *Product) MoveAfterID(ctx context.Context, id mrtype.KeyInt32, afterID mrtype.KeyInt32) error {
	if id < 1 {
		return mrcore.FactoryErrServiceEntityNotFound.New()
	}

	item := entity.Product{
		ID: id,
	}

	if err := uc.storage.LoadOne(ctx, &item); err != nil {
		return uc.serviceHelper.WrapErrorEntityNotFoundOrFailed(err, entity.ModelNameProduct, id)
	}

	if item.CategoryID < 1 {
		return mrcore.FactoryErrInternalWithData.New(entity.ModelNameProduct, mrmsg.Data{"categoryId": item.CategoryID})
	}

	meta := uc.storage.GetMetaData(item.CategoryID)
	component := uc.componentOrderer.WithMetaData(meta)

	if err := component.MoveAfterID(ctx, id, afterID); err != nil {
		return err
	}

	uc.eventBoxEmitEntity(ctx, "Move", mrmsg.Data{"id": id, "afterId": afterID})

	return nil
}

func (uc *Product) checkProduct(ctx context.Context, item *entity.Product) error {
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

func (uc *Product) checkArticle(ctx context.Context, item *entity.Product) error {
	id, err := uc.storage.FetchIdByArticle(ctx, item.Article)

	if err != nil {
		if mrcore.FactoryErrStorageNoRowFound.Is(err) {
			return nil
		}

		return uc.serviceHelper.WrapErrorFailed(err, entity.ModelNameProduct)
	}

	if item.ID != id {
		return usecase_shared.FactoryErrProductArticleAlreadyExists.New(item.Article)
	}

	return nil
}

func (uc *Product) eventBoxEmitEntity(ctx context.Context, eventName string, data mrmsg.Data) {
	uc.eventBox.Emit(
		"%s::%s: %s",
		entity.ModelNameProduct,
		eventName,
		data,
	)
}
