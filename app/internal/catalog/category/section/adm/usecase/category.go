package usecase

import (
	"context"

	"github.com/mondegor/go-sample/internal/catalog/category/section/adm/entity"

	"github.com/google/uuid"
	"github.com/mondegor/go-storage/mrentity"
	"github.com/mondegor/go-sysmess/mrmsg"
	"github.com/mondegor/go-webcore/mrcore"
	"github.com/mondegor/go-webcore/mrenum"
	"github.com/mondegor/go-webcore/mrpath"
	"github.com/mondegor/go-webcore/mrsender"
	"github.com/mondegor/go-webcore/mrstatus"
	"github.com/mondegor/go-webcore/mrstatus/mrflow"
)

type (
	// Category - comment struct.
	Category struct {
		storage      CategoryStorage
		eventEmitter mrsender.EventEmitter
		errorWrapper mrcore.UsecaseErrorWrapper
		imgBaseURL   mrpath.PathBuilder
		statusFlow   mrstatus.Flow
	}
)

// NewCategory - создаёт объект Category.
func NewCategory(
	storage CategoryStorage,
	eventEmitter mrsender.EventEmitter,
	errorWrapper mrcore.UsecaseErrorWrapper,
	imgBaseURL mrpath.PathBuilder,
) *Category {
	return &Category{
		storage:      storage,
		eventEmitter: eventEmitter,
		errorWrapper: errorWrapper,
		imgBaseURL:   imgBaseURL,
		statusFlow:   mrflow.ItemStatusFlow(),
	}
}

// GetList - comment method.
func (uc *Category) GetList(ctx context.Context, params entity.CategoryParams) ([]entity.Category, int64, error) {
	fetchParams := uc.storage.NewSelectParams(params)

	total, err := uc.storage.FetchTotal(ctx, fetchParams.Where)
	if err != nil {
		return nil, 0, uc.errorWrapper.WrapErrorFailed(err, entity.ModelNameCategory)
	}

	if total < 1 {
		return nil, 0, nil
	}

	items, err := uc.storage.Fetch(ctx, fetchParams)
	if err != nil {
		return nil, 0, uc.errorWrapper.WrapErrorFailed(err, entity.ModelNameCategory)
	}

	for i := range items {
		uc.prepareItem(&items[i])
	}

	return items, total, nil
}

// GetItem - comment method.
func (uc *Category) GetItem(ctx context.Context, itemID uuid.UUID) (entity.Category, error) {
	if itemID == uuid.Nil {
		return entity.Category{}, mrcore.ErrUseCaseEntityNotFound.New()
	}

	item, err := uc.storage.FetchOne(ctx, itemID)
	if err != nil {
		return entity.Category{}, uc.errorWrapper.WrapErrorEntityNotFoundOrFailed(err, entity.ModelNameCategory, itemID)
	}

	uc.prepareItem(&item)

	return item, nil
}

// Create - comment method.
func (uc *Category) Create(ctx context.Context, item entity.Category) (uuid.UUID, error) {
	item.Status = mrenum.ItemStatusDraft

	itemID, err := uc.storage.Insert(ctx, item)
	if err != nil {
		return uuid.Nil, uc.errorWrapper.WrapErrorFailed(err, entity.ModelNameCategory)
	}

	uc.emitEvent(ctx, "Create", mrmsg.Data{"id": itemID})

	return itemID, nil
}

// Store - comment method.
func (uc *Category) Store(ctx context.Context, item entity.Category) error {
	if item.ID == uuid.Nil {
		return mrcore.ErrUseCaseEntityNotFound.New()
	}

	if item.TagVersion < 1 {
		return mrcore.ErrUseCaseEntityVersionInvalid.New()
	}

	// предварительная проверка существования записи нужна для того,
	// чтобы при Update быть уверенным, что отсутствие записи из-за ошибки VersionInvalid
	if _, err := uc.storage.FetchStatus(ctx, item.ID); err != nil {
		return uc.errorWrapper.WrapErrorEntityNotFoundOrFailed(err, entity.ModelNameCategory, item.ID)
	}

	tagVersion, err := uc.storage.Update(ctx, item)
	if err != nil {
		if uc.errorWrapper.IsNotFoundError(err) {
			return mrcore.ErrUseCaseEntityVersionInvalid.Wrap(err)
		}

		return uc.errorWrapper.WrapErrorFailed(err, entity.ModelNameCategory)
	}

	uc.emitEvent(ctx, "Store", mrmsg.Data{"id": item.ID, "ver": tagVersion})

	return nil
}

// ChangeStatus - comment method.
func (uc *Category) ChangeStatus(ctx context.Context, item entity.Category) error {
	if item.ID == uuid.Nil {
		return mrcore.ErrUseCaseEntityNotFound.New()
	}

	if item.TagVersion < 1 {
		return mrcore.ErrUseCaseEntityVersionInvalid.New()
	}

	currentStatus, err := uc.storage.FetchStatus(ctx, item.ID)
	if err != nil {
		return uc.errorWrapper.WrapErrorEntityNotFoundOrFailed(err, entity.ModelNameCategory, item.ID)
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

		return uc.errorWrapper.WrapErrorFailed(err, entity.ModelNameCategory)
	}

	uc.emitEvent(ctx, "ChangeStatus", mrmsg.Data{"id": item.ID, "ver": tagVersion, "status": item.Status})

	return nil
}

// Remove - comment method.
func (uc *Category) Remove(ctx context.Context, itemID uuid.UUID) error {
	if itemID == uuid.Nil {
		return mrcore.ErrUseCaseEntityNotFound.New()
	}

	if err := uc.storage.Delete(ctx, itemID); err != nil {
		return uc.errorWrapper.WrapErrorEntityNotFoundOrFailed(err, entity.ModelNameCategory, itemID)
	}

	uc.emitEvent(ctx, "Remove", mrmsg.Data{"id": itemID})

	return nil
}

func (uc *Category) prepareItem(item *entity.Category) {
	if imageInfo := mrentity.ImageMetaToInfoPointer(item.ImageMeta, nil); imageInfo != nil {
		imageInfo.URL = uc.imgBaseURL.BuildPath(imageInfo.Path)
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
