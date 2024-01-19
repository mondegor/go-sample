package usecase

import (
	"context"
	"go-sample/internal/modules/catalog/entity/admin-api"

	"github.com/mondegor/go-storage/mrentity"
	"github.com/mondegor/go-sysmess/mrmsg"
	"github.com/mondegor/go-webcore/mrcore"
	"github.com/mondegor/go-webcore/mrenum"
	"github.com/mondegor/go-webcore/mrtool"
	"github.com/mondegor/go-webcore/mrtype"
)

type (
	Category struct {
		storage       CategoryStorage
		eventBox      mrcore.EventBox
		serviceHelper *mrtool.ServiceHelper
		imgBaseURL    mrcore.BuilderPath
		statusFlow    mrenum.StatusFlow
	}
)

func NewCategory(
	storage CategoryStorage,
	eventBox mrcore.EventBox,
	serviceHelper *mrtool.ServiceHelper,
	imgBaseURL mrcore.BuilderPath,
) *Category {
	return &Category{
		storage:       storage,
		eventBox:      eventBox,
		serviceHelper: serviceHelper,
		imgBaseURL:    imgBaseURL,
		statusFlow:    mrenum.ItemStatusFlow,
	}
}

func (uc *Category) GetList(ctx context.Context, params entity.CategoryParams) ([]entity.Category, int64, error) {
	fetchParams := uc.storage.NewFetchParams(params)
	total, err := uc.storage.FetchTotal(ctx, fetchParams.Where)

	if err != nil {
		return nil, 0, uc.serviceHelper.WrapErrorFailed(err, entity.ModelNameCategory)
	}

	if total < 1 {
		return []entity.Category{}, 0, nil
	}

	items, err := uc.storage.Fetch(ctx, fetchParams)

	if err != nil {
		return nil, 0, uc.serviceHelper.WrapErrorFailed(err, entity.ModelNameCategory)
	}

	for i := range items {
		uc.prepareItem(&items[i])
	}

	return items, total, nil
}

func (uc *Category) GetItem(ctx context.Context, id mrtype.KeyInt32) (*entity.Category, error) {
	if id < 1 {
		return nil, mrcore.FactoryErrServiceEntityNotFound.New()
	}

	item := &entity.Category{
		ID: id,
	}

	if err := uc.storage.LoadOne(ctx, item); err != nil {
		return nil, uc.serviceHelper.WrapErrorEntityNotFoundOrFailed(err, entity.ModelNameCategory, id)
	}

	uc.prepareItem(item)

	return item, nil
}

func (uc *Category) Create(ctx context.Context, item *entity.Category) error {
	item.Status = mrenum.ItemStatusDraft

	if err := uc.storage.Insert(ctx, item); err != nil {
		return uc.serviceHelper.WrapErrorFailed(err, entity.ModelNameCategory)
	}

	uc.eventBoxEmitEntity(ctx, "Create", mrmsg.Data{"id": item.ID})

	return nil
}

func (uc *Category) Store(ctx context.Context, item *entity.Category) error {
	if item.ID < 1 {
		return mrcore.FactoryErrServiceEntityNotFound.New()
	}

	if item.TagVersion < 1 {
		return mrcore.FactoryErrServiceEntityVersionInvalid.New()
	}

	if err := uc.storage.IsExists(ctx, item.ID); err != nil {
		return uc.serviceHelper.WrapErrorEntityNotFoundOrFailed(err, entity.ModelNameCategory, item.ID)
	}

	version, err := uc.storage.Update(ctx, item)

	if err != nil {
		if uc.serviceHelper.IsNotFoundError(err) {
			return mrcore.FactoryErrServiceEntityVersionInvalid.Wrap(err)
		}

		return uc.serviceHelper.WrapErrorFailed(err, entity.ModelNameCategory)
	}

	uc.eventBoxEmitEntity(ctx, "Store", mrmsg.Data{"id": item.ID, "ver": version})

	return nil
}

func (uc *Category) ChangeStatus(ctx context.Context, item *entity.Category) error {
	if item.ID < 1 {
		return mrcore.FactoryErrServiceEntityNotFound.New()
	}

	if item.TagVersion < 1 {
		return mrcore.FactoryErrServiceEntityVersionInvalid.New()
	}

	currentStatus, err := uc.storage.FetchStatus(ctx, item)

	if err != nil {
		return uc.serviceHelper.WrapErrorEntityNotFoundOrFailed(err, entity.ModelNameCategory, item.ID)
	}

	if currentStatus == item.Status {
		return nil
	}

	if !uc.statusFlow.Check(currentStatus, item.Status) {
		return mrcore.FactoryErrServiceSwitchStatusRejected.New(currentStatus, item.Status)
	}

	version, err := uc.storage.UpdateStatus(ctx, item)

	if err != nil {
		if uc.serviceHelper.IsNotFoundError(err) {
			return mrcore.FactoryErrServiceEntityVersionInvalid.Wrap(err)
		}

		return uc.serviceHelper.WrapErrorFailed(err, entity.ModelNameCategory)
	}

	uc.eventBoxEmitEntity(ctx, "ChangeStatus", mrmsg.Data{"id": item.ID, "ver": version, "status": item.Status})

	return nil
}

func (uc *Category) Remove(ctx context.Context, id mrtype.KeyInt32) error {
	if id < 1 {
		return mrcore.FactoryErrServiceEntityNotFound.New()
	}

	if err := uc.storage.Delete(ctx, id); err != nil {
		return uc.serviceHelper.WrapErrorEntityNotFoundOrFailed(err, entity.ModelNameCategory, id)
	}

	uc.eventBoxEmitEntity(ctx, "Remove", mrmsg.Data{"id": id})

	return nil
}

func (uc *Category) prepareItem(item *entity.Category) {
	if imageInfo := mrentity.ConvertImageMetaToInfo(item.ImageMeta); imageInfo != nil {
		imageInfo.URL = uc.imgBaseURL.FullPath(imageInfo.Path)
		item.ImageInfo = imageInfo
	}
}

func (uc *Category) eventBoxEmitEntity(ctx context.Context, eventName string, data mrmsg.Data) {
	uc.eventBox.Emit(
		"%s::%s: %s",
		entity.ModelNameCategory,
		eventName,
		data,
	)
}
