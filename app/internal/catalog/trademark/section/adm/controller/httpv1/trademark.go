package httpv1

import (
	"net/http"

	"github.com/mondegor/go-sample/internal/catalog/trademark/module"
	"github.com/mondegor/go-sample/internal/catalog/trademark/section/adm/entity"
	"github.com/mondegor/go-sample/internal/catalog/trademark/section/adm/usecase"
	"github.com/mondegor/go-sample/pkg/catalog/api"
	"github.com/mondegor/go-sample/pkg/validate"
	"github.com/mondegor/go-sample/pkg/view"

	"github.com/mondegor/go-sysmess/mrerr"
	"github.com/mondegor/go-webcore/mrcore"
	"github.com/mondegor/go-webcore/mrserver"
	"github.com/mondegor/go-webcore/mrtype"
	"github.com/mondegor/go-webcore/mrview"
)

const (
	trademarkListURL             = "/v1/catalog/trademarks"
	trademarkItemURL             = "/v1/catalog/trademarks/{id}"
	trademarkItemChangeStatusURL = "/v1/catalog/trademarks/{id}/status"
)

type (
	// Trademark - comment struct.
	Trademark struct {
		parser     validate.RequestExtendParser
		sender     mrserver.ResponseSender
		useCase    usecase.TrademarkUseCase
		listSorter mrview.ListSorter
	}
)

// NewTrademark - создаёт контроллер Trademark.
func NewTrademark(
	parser validate.RequestExtendParser,
	sender mrserver.ResponseSender,
	useCase usecase.TrademarkUseCase,
	listSorter mrview.ListSorter,
) *Trademark {
	return &Trademark{
		parser:     parser,
		sender:     sender,
		useCase:    useCase,
		listSorter: listSorter,
	}
}

// Handlers - возвращает обработчики контроллера Trademark.
func (ht *Trademark) Handlers() []mrserver.HttpHandler {
	return []mrserver.HttpHandler{
		{Method: http.MethodGet, URL: trademarkListURL, Func: ht.GetList},
		{Method: http.MethodPost, URL: trademarkListURL, Func: ht.Create},

		{Method: http.MethodGet, URL: trademarkItemURL, Func: ht.Get},
		{Method: http.MethodPut, URL: trademarkItemURL, Func: ht.Store},
		{Method: http.MethodDelete, URL: trademarkItemURL, Func: ht.Remove},

		{Method: http.MethodPatch, URL: trademarkItemChangeStatusURL, Func: ht.ChangeStatus},
	}
}

// GetList - comment method.
func (ht *Trademark) GetList(w http.ResponseWriter, r *http.Request) error {
	items, totalItems, err := ht.useCase.GetList(r.Context(), ht.listParams(r))
	if err != nil {
		return err
	}

	return ht.sender.Send(
		w,
		http.StatusOK,
		TrademarkListResponse{
			Items: items,
			Total: totalItems,
		},
	)
}

func (ht *Trademark) listParams(r *http.Request) entity.TrademarkParams {
	return entity.TrademarkParams{
		Filter: entity.TrademarkListFilter{
			SearchText: ht.parser.FilterString(r, module.ParamNameFilterSearchText),
			Statuses:   ht.parser.FilterStatusList(r, module.ParamNameFilterStatuses),
		},
		Sorter: ht.parser.SortParams(r, ht.listSorter),
		Pager:  ht.parser.PageParams(r),
	}
}

// Get - comment method.
func (ht *Trademark) Get(w http.ResponseWriter, r *http.Request) error {
	item, err := ht.useCase.GetItem(r.Context(), ht.getItemID(r))
	if err != nil {
		return ht.wrapError(err, r)
	}

	return ht.sender.Send(w, http.StatusOK, item)
}

// Create - comment method.
func (ht *Trademark) Create(w http.ResponseWriter, r *http.Request) error {
	request := CreateTrademarkRequest{}

	if err := ht.parser.Validate(r, &request); err != nil {
		return err
	}

	item := entity.Trademark{
		Caption: request.Caption,
	}

	itemID, err := ht.useCase.Create(r.Context(), item)
	if err != nil {
		return ht.wrapError(err, r)
	}

	return ht.sender.Send(
		w,
		http.StatusCreated,
		view.SuccessCreatedItemInt32Response{
			ItemID: itemID,
		},
	)
}

// Store - comment method.
func (ht *Trademark) Store(w http.ResponseWriter, r *http.Request) error {
	request := StoreTrademarkRequest{}

	if err := ht.parser.Validate(r, &request); err != nil {
		return err
	}

	item := entity.Trademark{
		ID:         ht.getItemID(r),
		TagVersion: request.TagVersion,
		Caption:    request.Caption,
	}

	if err := ht.useCase.Store(r.Context(), item); err != nil {
		return ht.wrapError(err, r)
	}

	return ht.sender.SendNoContent(w)
}

// ChangeStatus - comment method.
func (ht *Trademark) ChangeStatus(w http.ResponseWriter, r *http.Request) error {
	request := view.ChangeItemStatusRequest{}

	if err := ht.parser.Validate(r, &request); err != nil {
		return err
	}

	item := entity.Trademark{
		ID:         ht.getItemID(r),
		TagVersion: request.TagVersion,
		Status:     request.Status,
	}

	if err := ht.useCase.ChangeStatus(r.Context(), item); err != nil {
		return ht.wrapError(err, r)
	}

	return ht.sender.SendNoContent(w)
}

// Remove - comment method.
func (ht *Trademark) Remove(w http.ResponseWriter, r *http.Request) error {
	if err := ht.useCase.Remove(r.Context(), ht.getItemID(r)); err != nil {
		return ht.wrapError(err, r)
	}

	return ht.sender.SendNoContent(w)
}

func (ht *Trademark) getItemID(r *http.Request) mrtype.KeyInt32 {
	return ht.parser.PathKeyInt32(r, "id")
}

func (ht *Trademark) getRawItemID(r *http.Request) string {
	return ht.parser.PathParamString(r, "id")
}

func (ht *Trademark) wrapError(err error, r *http.Request) error {
	if mrcore.ErrUseCaseEntityNotFound.Is(err) {
		return api.ErrTrademarkNotFound.Wrap(err, ht.getRawItemID(r))
	}

	if mrcore.ErrUseCaseEntityVersionInvalid.Is(err) {
		return mrerr.NewCustomError("tagVersion", err)
	}

	if mrcore.ErrUseCaseSwitchStatusRejected.Is(err) {
		return mrerr.NewCustomError("status", err)
	}

	return err
}
