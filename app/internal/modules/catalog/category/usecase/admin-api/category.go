package usecase

import (
	"context"
	"go-sample/internal/modules/catalog/category/entity/admin-api"

	"github.com/mondegor/go-storage/mrentity"
	"github.com/mondegor/go-sysmess/mrmsg"
	"github.com/mondegor/go-webcore/mrcore"
	"github.com/mondegor/go-webcore/mrenum"
	"github.com/mondegor/go-webcore/mrlib"
	"github.com/mondegor/go-webcore/mrsender"
	"github.com/mondegor/go-webcore/mrtype"
)

type (
	Category struct {
		storage       CategoryStorage
		eventEmitter  mrsender.EventEmitter
		usecaseHelper *mrcore.UsecaseHelper
		imgBaseURL    mrlib.BuilderPath
		statusFlow    mrenum.StatusFlow
	}
)

func NewCategory(
	storage CategoryStorage,
	eventEmitter mrsender.EventEmitter,
	usecaseHelper *mrcore.UsecaseHelper,
	imgBaseURL mrlib.BuilderPath,
) *Category {
	return &Category{
		storage:       storage,
		eventEmitter:  eventEmitter,
		usecaseHelper: usecaseHelper,
		imgBaseURL:    imgBaseURL,
		statusFlow:    mrenum.ItemStatusFlow,
	}
}

func (uc *Category) GetList(ctx context.Context, params entity.CategoryParams) ([]entity.Category, int64, error) {
	fetchParams := uc.storage.NewFetchParams(params)
	total, err := uc.storage.FetchTotal(ctx, fetchParams.Where)

	if err != nil {
		return nil, 0, uc.usecaseHelper.WrapErrorFailed(err, entity.ModelNameCategory)
	}

	if total < 1 {
		return []entity.Category{}, 0, nil
	}

	items, err := uc.storage.Fetch(ctx, fetchParams)

	if err != nil {
		return nil, 0, uc.usecaseHelper.WrapErrorFailed(err, entity.ModelNameCategory)
	}

	for i := range items {
		uc.prepareItem(&items[i])
	}

	return items, total, nil
}

func (uc *Category) GetItem(ctx context.Context, itemID mrtype.KeyInt32) (entity.Category, error) {
	if itemID < 1 {
		return entity.Category{}, mrcore.FactoryErrUseCaseEntityNotFound.New()
	}

	item, err := uc.storage.FetchOne(ctx, itemID)

	if err != nil {
		return entity.Category{}, uc.usecaseHelper.WrapErrorEntityNotFoundOrFailed(err, entity.ModelNameCategory, itemID)
	}

	uc.prepareItem(&item)

	return item, nil
}

func (uc *Category) Create(ctx context.Context, item entity.Category) (mrtype.KeyInt32, error) {
	item.Status = mrenum.ItemStatusDraft
	itemID, err := uc.storage.Insert(ctx, item)

	if err != nil {
		return 0, uc.usecaseHelper.WrapErrorFailed(err, entity.ModelNameCategory)
	}

	uc.emitEvent(ctx, "Create", mrmsg.Data{"id": itemID})

	return itemID, nil
}

func (uc *Category) Store(ctx context.Context, item entity.Category) error {
	if item.ID < 1 {
		return mrcore.FactoryErrUseCaseEntityNotFound.New()
	}

	if item.TagVersion < 1 {
		return mrcore.FactoryErrUseCaseEntityVersionInvalid.New()
	}

	if err := uc.storage.IsExists(ctx, item.ID); err != nil {
		return uc.usecaseHelper.WrapErrorEntityNotFoundOrFailed(err, entity.ModelNameCategory, item.ID)
	}

	version, err := uc.storage.Update(ctx, item)

	if err != nil {
		if uc.usecaseHelper.IsNotFoundError(err) {
			return mrcore.FactoryErrUseCaseEntityVersionInvalid.Wrap(err)
		}

		return uc.usecaseHelper.WrapErrorFailed(err, entity.ModelNameCategory)
	}

	uc.emitEvent(ctx, "Store", mrmsg.Data{"id": item.ID, "ver": version})

	return nil
}

func (uc *Category) ChangeStatus(ctx context.Context, item entity.Category) error {
	if item.ID < 1 {
		return mrcore.FactoryErrUseCaseEntityNotFound.New()
	}

	if item.TagVersion < 1 {
		return mrcore.FactoryErrUseCaseEntityVersionInvalid.New()
	}

	currentStatus, err := uc.storage.FetchStatus(ctx, item)

	if err != nil {
		return uc.usecaseHelper.WrapErrorEntityNotFoundOrFailed(err, entity.ModelNameCategory, item.ID)
	}

	if currentStatus == item.Status {
		return nil
	}

	if !uc.statusFlow.Check(currentStatus, item.Status) {
		return mrcore.FactoryErrUseCaseSwitchStatusRejected.New(currentStatus, item.Status)
	}

	version, err := uc.storage.UpdateStatus(ctx, item)

	if err != nil {
		if uc.usecaseHelper.IsNotFoundError(err) {
			return mrcore.FactoryErrUseCaseEntityVersionInvalid.Wrap(err)
		}

		return uc.usecaseHelper.WrapErrorFailed(err, entity.ModelNameCategory)
	}

	uc.emitEvent(ctx, "ChangeStatus", mrmsg.Data{"id": item.ID, "ver": version, "status": item.Status})

	return nil
}

func (uc *Category) Remove(ctx context.Context, itemID mrtype.KeyInt32) error {
	if itemID < 1 {
		return mrcore.FactoryErrUseCaseEntityNotFound.New()
	}

	if err := uc.storage.Delete(ctx, itemID); err != nil {
		return uc.usecaseHelper.WrapErrorEntityNotFoundOrFailed(err, entity.ModelNameCategory, itemID)
	}

	uc.emitEvent(ctx, "Remove", mrmsg.Data{"id": itemID})

	return nil
}

func (uc *Category) prepareItem(item *entity.Category) {
	if imageInfo := mrentity.ImageMetaToInfoPointer(item.ImageMeta); imageInfo != nil {
		imageInfo.URL = uc.imgBaseURL.FullPath(imageInfo.Path)
		item.ImageInfo = imageInfo
	}
}

func (uc *Category) emitEvent(ctx context.Context, eventName string, data mrmsg.Data) {
	uc.eventEmitter.EmitWithSource(
		ctx,
		eventName,
		entity.ModelNameCategory,
		data,
	)
}
