package usecase

import (
	"context"
	"go-sample/internal/modules/catalog/entity/admin-api"

	"github.com/mondegor/go-sysmess/mrmsg"
	"github.com/mondegor/go-webcore/mrcore"
	"github.com/mondegor/go-webcore/mrenum"
	"github.com/mondegor/go-webcore/mrtool"
	"github.com/mondegor/go-webcore/mrtype"
)

type (
	Trademark struct {
		storage       TrademarkStorage
		eventBox      mrcore.EventBox
		serviceHelper *mrtool.ServiceHelper
		statusFlow    mrenum.StatusFlow
	}
)

func NewTrademark(
	storage TrademarkStorage,
	eventBox mrcore.EventBox,
	serviceHelper *mrtool.ServiceHelper,
) *Trademark {
	return &Trademark{
		storage:       storage,
		eventBox:      eventBox,
		serviceHelper: serviceHelper,
		statusFlow:    mrenum.ItemStatusFlow,
	}
}

func (uc *Trademark) GetList(ctx context.Context, params entity.TrademarkParams) ([]entity.Trademark, int64, error) {
	fetchParams := uc.storage.NewFetchParams(params)
	total, err := uc.storage.FetchTotal(ctx, fetchParams.Where)

	if err != nil {
		return nil, 0, uc.serviceHelper.WrapErrorFailed(err, entity.ModelNameTrademark)
	}

	if total < 1 {
		return []entity.Trademark{}, 0, nil
	}

	items, err := uc.storage.Fetch(ctx, fetchParams)

	if err != nil {
		return nil, 0, uc.serviceHelper.WrapErrorFailed(err, entity.ModelNameTrademark)
	}

	return items, total, nil
}

func (uc *Trademark) GetItem(ctx context.Context, id mrtype.KeyInt32) (*entity.Trademark, error) {
	if id < 1 {
		return nil, mrcore.FactoryErrServiceEntityNotFound.New()
	}

	item := &entity.Trademark{ID: id}

	if err := uc.storage.LoadOne(ctx, item); err != nil {
		return nil, uc.serviceHelper.WrapErrorEntityNotFoundOrFailed(err, entity.ModelNameTrademark, id)
	}

	return item, nil
}

// Create
// modifies: item{ID}
func (uc *Trademark) Create(ctx context.Context, item *entity.Trademark) error {
	item.Status = mrenum.ItemStatusDraft

	if err := uc.storage.Insert(ctx, item); err != nil {
		return uc.serviceHelper.WrapErrorFailed(err, entity.ModelNameTrademark)
	}

	uc.eventBoxEmitEntity(ctx, "Create", mrmsg.Data{"id": item.ID})

	return nil
}

func (uc *Trademark) Store(ctx context.Context, item *entity.Trademark) error {
	if item.ID < 1 {
		return mrcore.FactoryErrServiceEntityNotFound.New()
	}

	if item.TagVersion < 1 {
		return mrcore.FactoryErrServiceEntityVersionInvalid.New()
	}

	if err := uc.storage.IsExists(ctx, item.ID); err != nil {
		return uc.serviceHelper.WrapErrorEntityNotFoundOrFailed(err, entity.ModelNameTrademark, item.ID)
	}

	version, err := uc.storage.Update(ctx, item)

	if err != nil {
		return uc.serviceHelper.WrapErrorEntity(
			mrcore.FactoryErrServiceEntityVersionInvalid,
			err,
			entity.ModelNameTrademark,
			mrmsg.Data{"id": item.ID, "ver": item.TagVersion},
		)
	}

	uc.eventBoxEmitEntity(ctx, "Store", mrmsg.Data{"id": item.ID, "ver": version})

	return nil
}

func (uc *Trademark) ChangeStatus(ctx context.Context, item *entity.Trademark) error {
	if item.ID < 1 {
		return mrcore.FactoryErrServiceEntityNotFound.New()
	}

	if item.TagVersion < 1 {
		return mrcore.FactoryErrServiceEntityVersionInvalid.New()
	}

	currentStatus, err := uc.storage.FetchStatus(ctx, item)

	if err != nil {
		return uc.serviceHelper.WrapErrorEntityNotFoundOrFailed(err, entity.ModelNameTrademark, item.ID)
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
			entity.ModelNameTrademark,
			mrmsg.Data{"id": item.ID, "ver": item.TagVersion},
		)
	}

	uc.eventBoxEmitEntity(ctx, "ChangeStatus", mrmsg.Data{"id": item.ID, "ver": version, "status": item.Status})

	return nil
}

func (uc *Trademark) Remove(ctx context.Context, id mrtype.KeyInt32) error {
	if id < 1 {
		return mrcore.FactoryErrServiceEntityNotFound.New()
	}

	if err := uc.storage.Delete(ctx, id); err != nil {
		return uc.serviceHelper.WrapErrorEntityNotFoundOrFailed(err, entity.ModelNameTrademark, id)
	}

	uc.eventBoxEmitEntity(ctx, "Remove", mrmsg.Data{"id": id})

	return nil
}

func (uc *Trademark) eventBoxEmitEntity(ctx context.Context, eventName string, data mrmsg.Data) {
	uc.eventBox.Emit(
		"%s::%s: %s",
		entity.ModelNameTrademark,
		eventName,
		data,
	)
}
