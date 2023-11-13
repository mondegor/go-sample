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
	Trademark struct {
		storage TrademarkStorage
		eventBox mrcore.EventBox
		serviceHelper *mrtool.ServiceHelper
		statusFlow mrenum.StatusFlow
	}
)

func NewTrademark(
	storage TrademarkStorage,
	eventBox mrcore.EventBox,
	serviceHelper *mrtool.ServiceHelper,
) *Trademark {
	return &Trademark{
		storage: storage,
		eventBox: eventBox,
		serviceHelper: serviceHelper,
		statusFlow: mrenum.ItemStatusFlow,
	}
}

func (uc *Trademark) GetList(ctx context.Context, params entity.TrademarkParams) ([]entity.Trademark, int64, error) {
	fetchParams := uc.storage.NewFetchParams(params)
	total, err := uc.storage.FetchTotal(ctx, fetchParams.Where)

	if err != nil {
		return nil, 0, mrcore.FactoryErrServiceTemporarilyUnavailable.Wrap(err, entity.ModelNameCatalogTrademark)
	}

	if total < 1 {
		return []entity.Trademark{}, 0, nil
	}

	items, err := uc.storage.Fetch(ctx, fetchParams)

	if err != nil {
		return nil, 0, mrcore.FactoryErrServiceTemporarilyUnavailable.Wrap(err, entity.ModelNameCatalogTrademark)
	}

	return items, total, nil
}

func (uc *Trademark) GetItem(ctx context.Context, id mrtype.KeyInt32) (*entity.Trademark, error) {
	if id < 1 {
		return nil, mrcore.FactoryErrServiceIncorrectInputData.New(mrerr.Arg{"id": id})
	}

	item := &entity.Trademark{ID: id}

	if err := uc.storage.LoadOne(ctx, item); err != nil {
		return nil, uc.serviceHelper.WrapErrorForSelect(err, entity.ModelNameCatalogTrademark)
	}

	return item, nil
}

// Create
// modifies: item{ID}
func (uc *Trademark) Create(ctx context.Context, item *entity.Trademark) error {
	item.Status = mrenum.ItemStatusDraft

	if err := uc.storage.Insert(ctx, item); err != nil {
		return mrcore.FactoryErrServiceEntityNotCreated.Wrap(err, entity.ModelNameCatalogTrademark)
	}

	uc.eventBox.Emit(
		"%s::Create: id=%d",
		entity.ModelNameCatalogTrademark,
		item.ID,
	)

	return nil
}

func (uc *Trademark) Store(ctx context.Context, item *entity.Trademark) error {
	if item.ID < 1 || item.TagVersion < 1 {
		return mrcore.FactoryErrServiceIncorrectInputData.New(mrerr.Arg{"item.id": item.ID, "version": item.TagVersion})
	}

	if err := uc.storage.Update(ctx, item); err != nil {
		return uc.serviceHelper.WrapErrorForUpdate(err, entity.ModelNameCatalogTrademark)
	}

	uc.eventBox.Emit(
		"%s::Store: id=%d",
		entity.ModelNameCatalogTrademark,
		item.ID,
	)

	return nil
}

func (uc *Trademark) ChangeStatus(ctx context.Context, item *entity.Trademark) error {
	if item.ID < 1 || item.TagVersion < 1 {
		return mrcore.FactoryErrServiceIncorrectInputData.New(mrerr.Arg{"item.id": item.ID, "version": item.TagVersion})
	}

	currentStatus, err := uc.storage.FetchStatus(ctx, item)

	if err != nil {
		return uc.serviceHelper.WrapErrorForSelect(err, entity.ModelNameCatalogTrademark)
	}

	if currentStatus == item.Status {
		return nil
	}

	if !uc.statusFlow.Check(currentStatus, item.Status) {
		return mrcore.FactoryErrServiceIncorrectSwitchStatus.New(currentStatus, item.Status, entity.ModelNameCatalogTrademark, item.ID)
	}

	err = uc.storage.UpdateStatus(ctx, item)

	if err != nil {
		return uc.serviceHelper.WrapErrorForUpdate(err, entity.ModelNameCatalogTrademark)
	}

	uc.eventBox.Emit(
		"%s::ChangeStatus: id=%d, status=%s",
		entity.ModelNameCatalogTrademark,
		item.ID,
		item.Status,
	)

	return nil
}

func (uc *Trademark) Remove(ctx context.Context, id mrtype.KeyInt32) error {
	if id < 1 {
		return mrcore.FactoryErrServiceIncorrectInputData.New(mrerr.Arg{"id": id})
	}

	if err := uc.storage.Delete(ctx, id); err != nil {
		return uc.serviceHelper.WrapErrorForRemove(err, entity.ModelNameCatalogTrademark)
	}

	uc.eventBox.Emit(
		"%s::Remove: id=%d",
		entity.ModelNameCatalogTrademark,
		id,
	)

	return nil
}
