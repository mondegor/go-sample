package usecase

import (
	"context"
	"go-sample/internal/modules/catalog/entity/admin-api"

	"github.com/mondegor/go-sysmess/mrerr"
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
		statusFlow    mrenum.StatusFlow
	}
)

func NewCategory(
	storage CategoryStorage,
	eventBox mrcore.EventBox,
	serviceHelper *mrtool.ServiceHelper,
) *Category {
	return &Category{
		storage:       storage,
		eventBox:      eventBox,
		serviceHelper: serviceHelper,
		statusFlow:    mrenum.ItemStatusFlow,
	}
}

func (uc *Category) GetList(ctx context.Context, params entity.CategoryParams) ([]entity.Category, int64, error) {
	fetchParams := uc.storage.NewFetchParams(params)
	total, err := uc.storage.FetchTotal(ctx, fetchParams.Where)

	if err != nil {
		return nil, 0, mrcore.FactoryErrServiceTemporarilyUnavailable.Wrap(err, entity.ModelNameCatalogCategory)
	}

	if total < 1 {
		return []entity.Category{}, 0, nil
	}

	items, err := uc.storage.Fetch(ctx, fetchParams)

	if err != nil {
		return nil, 0, mrcore.FactoryErrServiceTemporarilyUnavailable.Wrap(err, entity.ModelNameCatalogCategory)
	}

	return items, total, nil
}

func (uc *Category) GetItem(ctx context.Context, id mrtype.KeyInt32) (*entity.Category, error) {
	if id < 1 {
		return nil, mrcore.FactoryErrServiceEntityNotFound.New(entity.ModelNameCatalogCategory)
	}

	item := &entity.Category{ID: id}

	if err := uc.storage.LoadOne(ctx, item); err != nil {
		return nil, uc.serviceHelper.WrapErrorForSelect(err, entity.ModelNameCatalogCategory)
	}

	return item, nil
}

func (uc *Category) CheckAvailability(ctx context.Context, id mrtype.KeyInt32) error {
	if id < 1 {
		return mrcore.FactoryErrServiceEntityNotFound.New(entity.ModelNameCatalogCategory)
	}

	err := uc.storage.IsExists(ctx, id)

	return uc.serviceHelper.ReturnErrorIfItemNotFound(err, entity.ModelNameCatalogCategory)
}

// Create
// modifies: item{ID}
func (uc *Category) Create(ctx context.Context, item *entity.Category) error {
	item.Status = mrenum.ItemStatusDraft

	if err := uc.storage.Insert(ctx, item); err != nil {
		return mrcore.FactoryErrServiceEntityNotCreated.Wrap(err, entity.ModelNameCatalogCategory)
	}

	uc.eventBox.Emit(
		"%s::Create: id=%d",
		entity.ModelNameCatalogCategory,
		item.ID,
	)

	return nil
}

func (uc *Category) Store(ctx context.Context, item *entity.Category) error {
	if item.ID < 1 {
		return mrcore.FactoryErrServiceEntityNotFound.New(entity.ModelNameCatalogCategory)
	}

	if item.TagVersion < 1 {
		return mrcore.FactoryErrServiceIncorrectInputData.New(mrerr.Arg{"tagVersion": item.TagVersion})
	}

	if err := uc.storage.IsExists(ctx, item.ID); err != nil {
		return uc.serviceHelper.WrapErrorForSelect(err, entity.ModelNameCatalogCategory)
	}

	version, err := uc.storage.Update(ctx, item)

	if err != nil {
		return uc.serviceHelper.WrapErrorForUpdateWithVersion(err, entity.ModelNameCatalogCategory)
	}

	uc.eventBox.Emit(
		"%s::Store: id=%d, ver=%d",
		entity.ModelNameCatalogCategory,
		item.ID,
		version,
	)

	return nil
}

func (uc *Category) ChangeStatus(ctx context.Context, item *entity.Category) error {
	if item.ID < 1 {
		return mrcore.FactoryErrServiceEntityNotFound.New(entity.ModelNameCatalogCategory)
	}

	if item.TagVersion < 1 {
		return mrcore.FactoryErrServiceIncorrectInputData.New(mrerr.Arg{"tagVersion": item.TagVersion})
	}

	currentStatus, err := uc.storage.FetchStatus(ctx, item)

	if err != nil {
		return uc.serviceHelper.WrapErrorForSelect(err, entity.ModelNameCatalogCategory)
	}

	if currentStatus == item.Status {
		return nil
	}

	if !uc.statusFlow.Check(currentStatus, item.Status) {
		return mrcore.FactoryErrServiceIncorrectSwitchStatus.New(currentStatus, item.Status, entity.ModelNameCatalogCategory, item.ID)
	}

	version, err := uc.storage.UpdateStatus(ctx, item)

	if err != nil {
		return uc.serviceHelper.WrapErrorForUpdateWithVersion(err, entity.ModelNameCatalogCategory)
	}

	uc.eventBox.Emit(
		"%s::ChangeStatus: id=%d, ver=%d, status=%s",
		entity.ModelNameCatalogCategory,
		item.ID,
		version,
		item.Status,
	)

	return nil
}

func (uc *Category) Remove(ctx context.Context, id mrtype.KeyInt32) error {
	if id < 1 {
		return mrcore.FactoryErrServiceEntityNotFound.New(entity.ModelNameCatalogCategory)
	}

	if err := uc.storage.Delete(ctx, id); err != nil {
		return uc.serviceHelper.WrapErrorForRemove(err, entity.ModelNameCatalogCategory)
	}

	uc.eventBox.Emit(
		"%s::Remove: id=%d",
		entity.ModelNameCatalogCategory,
		id,
	)

	return nil
}
