package usecase

import (
	"context"

	"github.com/mondegor/go-sample/internal/catalog/trademark/section/adm"
	"github.com/mondegor/go-sample/internal/catalog/trademark/section/adm/entity"

	"github.com/mondegor/go-sysmess/mrmsg"
	"github.com/mondegor/go-webcore/mrcore"
	"github.com/mondegor/go-webcore/mrenum"
	"github.com/mondegor/go-webcore/mrsender"
	"github.com/mondegor/go-webcore/mrstatus"
	"github.com/mondegor/go-webcore/mrstatus/mrflow"
	"github.com/mondegor/go-webcore/mrtype"
)

type (
	// Trademark - comment struct.
	Trademark struct {
		storage      adm.TrademarkStorage
		eventEmitter mrsender.EventEmitter
		errorWrapper mrcore.UseCaseErrorWrapper
		statusFlow   mrstatus.Flow
	}
)

// NewTrademark - создаёт объект Trademark.
func NewTrademark(storage adm.TrademarkStorage, eventEmitter mrsender.EventEmitter, errorWrapper mrcore.UseCaseErrorWrapper) *Trademark {
	return &Trademark{
		storage:      storage,
		eventEmitter: eventEmitter,
		errorWrapper: errorWrapper,
		statusFlow:   mrflow.ItemStatusFlow(),
	}
}

// GetList - comment method.
func (uc *Trademark) GetList(ctx context.Context, params entity.TrademarkParams) ([]entity.Trademark, int64, error) {
	fetchParams := uc.storage.NewSelectParams(params)

	total, err := uc.storage.FetchTotal(ctx, fetchParams.Where)
	if err != nil {
		return nil, 0, uc.errorWrapper.WrapErrorFailed(err, entity.ModelNameTrademark)
	}

	if total < 1 {
		return make([]entity.Trademark, 0), 0, nil
	}

	items, err := uc.storage.Fetch(ctx, fetchParams)
	if err != nil {
		return nil, 0, uc.errorWrapper.WrapErrorFailed(err, entity.ModelNameTrademark)
	}

	return items, total, nil
}

// GetItem - comment method.
func (uc *Trademark) GetItem(ctx context.Context, itemID mrtype.KeyInt32) (entity.Trademark, error) {
	if itemID < 1 {
		return entity.Trademark{}, mrcore.ErrUseCaseEntityNotFound.New()
	}

	item, err := uc.storage.FetchOne(ctx, itemID)
	if err != nil {
		return entity.Trademark{}, uc.errorWrapper.WrapErrorEntityNotFoundOrFailed(err, entity.ModelNameTrademark, itemID)
	}

	return item, nil
}

// Create - comment method.
func (uc *Trademark) Create(ctx context.Context, item entity.Trademark) (mrtype.KeyInt32, error) {
	item.Status = mrenum.ItemStatusDraft

	itemID, err := uc.storage.Insert(ctx, item)
	if err != nil {
		return 0, uc.errorWrapper.WrapErrorFailed(err, entity.ModelNameTrademark)
	}

	uc.emitEvent(ctx, "Create", mrmsg.Data{"id": itemID})

	return itemID, nil
}

// Store - comment method.
func (uc *Trademark) Store(ctx context.Context, item entity.Trademark) error {
	if item.ID < 1 {
		return mrcore.ErrUseCaseEntityNotFound.New()
	}

	if item.TagVersion < 1 {
		return mrcore.ErrUseCaseEntityVersionInvalid.New()
	}

	// предварительная проверка существования записи нужна для того,
	// чтобы при Update быть уверенным, что отсутствие записи из-за ошибки VersionInvalid
	if _, err := uc.storage.FetchStatus(ctx, item.ID); err != nil {
		return uc.errorWrapper.WrapErrorEntityNotFoundOrFailed(err, entity.ModelNameTrademark, item.ID)
	}

	tagVersion, err := uc.storage.Update(ctx, item)
	if err != nil {
		if uc.errorWrapper.IsNotFoundError(err) {
			return mrcore.ErrUseCaseEntityVersionInvalid.Wrap(err)
		}

		return uc.errorWrapper.WrapErrorFailed(err, entity.ModelNameTrademark)
	}

	uc.emitEvent(ctx, "Store", mrmsg.Data{"id": item.ID, "ver": tagVersion})

	return nil
}

// ChangeStatus - comment method.
func (uc *Trademark) ChangeStatus(ctx context.Context, item entity.Trademark) error {
	if item.ID < 1 {
		return mrcore.ErrUseCaseEntityNotFound.New()
	}

	if item.TagVersion < 1 {
		return mrcore.ErrUseCaseEntityVersionInvalid.New()
	}

	currentStatus, err := uc.storage.FetchStatus(ctx, item.ID)
	if err != nil {
		return uc.errorWrapper.WrapErrorEntityNotFoundOrFailed(err, entity.ModelNameTrademark, item.ID)
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

		return uc.errorWrapper.WrapErrorFailed(err, entity.ModelNameTrademark)
	}

	uc.emitEvent(ctx, "ChangeStatus", mrmsg.Data{"id": item.ID, "ver": tagVersion, "status": item.Status})

	return nil
}

// Remove - comment method.
func (uc *Trademark) Remove(ctx context.Context, itemID mrtype.KeyInt32) error {
	if itemID < 1 {
		return mrcore.ErrUseCaseEntityNotFound.New()
	}

	if err := uc.storage.Delete(ctx, itemID); err != nil {
		return uc.errorWrapper.WrapErrorEntityNotFoundOrFailed(err, entity.ModelNameTrademark, itemID)
	}

	uc.emitEvent(ctx, "Remove", mrmsg.Data{"id": itemID})

	return nil
}

func (uc *Trademark) emitEvent(ctx context.Context, eventName string, data mrmsg.Data) {
	uc.eventEmitter.EmitWithSource(
		ctx,
		eventName,
		entity.ModelNameTrademark,
		data,
	)
}
