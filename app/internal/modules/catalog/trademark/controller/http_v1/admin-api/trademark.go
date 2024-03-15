package http_v1

import (
	module "go-sample/internal/modules/catalog/trademark"
	view_shared "go-sample/internal/modules/catalog/trademark/controller/http_v1/shared/view"
	entity "go-sample/internal/modules/catalog/trademark/entity/admin-api"
	usecase "go-sample/internal/modules/catalog/trademark/usecase/admin-api"
	"go-sample/pkg/modules/catalog"
	"net/http"
	"strconv"

	"github.com/mondegor/go-sysmess/mrerr"
	"github.com/mondegor/go-sysmess/mrlang"
	"github.com/mondegor/go-webcore/mrcore"

	"github.com/mondegor/go-webcore/mrserver"
	"github.com/mondegor/go-webcore/mrtype"
	"github.com/mondegor/go-webcore/mrview"
)

const (
	trademarkListURL             = "/v1/catalog/trademarks"
	trademarkItemURL             = "/v1/catalog/trademarks/:id"
	trademarkItemChangeStatusURL = "/v1/catalog/trademarks/:id/status"
)

type (
	Trademark struct {
		parser     view_shared.RequestParser
		sender     mrserver.ResponseSender
		useCase    usecase.TrademarkUseCase
		listSorter mrview.ListSorter
	}
)

func NewTrademark(
	parser view_shared.RequestParser,
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

func (ht *Trademark) Handlers() []mrserver.HttpHandler {
	return []mrserver.HttpHandler{
		{http.MethodGet, trademarkListURL, "", ht.GetList},
		{http.MethodPost, trademarkListURL, "", ht.Create},

		{http.MethodGet, trademarkItemURL, "", ht.Get},
		{http.MethodPut, trademarkItemURL, "", ht.Store},
		{http.MethodDelete, trademarkItemURL, "", ht.Remove},

		{http.MethodPut, trademarkItemChangeStatusURL, "", ht.ChangeStatus},
	}
}

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

func (ht *Trademark) Get(w http.ResponseWriter, r *http.Request) error {
	item, err := ht.useCase.GetItem(r.Context(), ht.getItemID(r))

	if err != nil {
		return ht.wrapError(err, r)
	}

	return ht.sender.Send(w, http.StatusOK, item)
}

func (ht *Trademark) Create(w http.ResponseWriter, r *http.Request) error {
	request := CreateTrademarkRequest{}

	if err := ht.parser.Validate(r, &request); err != nil {
		return err
	}

	item := entity.Trademark{
		Caption: request.Caption,
	}

	if itemID, err := ht.useCase.Create(r.Context(), item); err != nil {
		return ht.wrapError(err, r)
	} else {
		return ht.sender.Send(
			w,
			http.StatusCreated,
			SuccessCreatedItemResponse{
				ItemID: strconv.Itoa(int(itemID)),
				Message: mrlang.Ctx(r.Context()).TranslateMessage(
					"msgCatalogTrademarkSuccessCreated",
					"entity has been success created",
				),
			},
		)
	}
}

func (ht *Trademark) Store(w http.ResponseWriter, r *http.Request) error {
	request := StoreTrademarkRequest{}

	if err := ht.parser.Validate(r, &request); err != nil {
		return err
	}

	item := entity.Trademark{
		ID:         ht.getItemID(r),
		TagVersion: request.Version,
		Caption:    request.Caption,
	}

	if err := ht.useCase.Store(r.Context(), item); err != nil {
		return ht.wrapError(err, r)
	}

	return ht.sender.SendNoContent(w)
}

func (ht *Trademark) ChangeStatus(w http.ResponseWriter, r *http.Request) error {
	request := ChangeItemStatusRequest{}

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
	if mrcore.FactoryErrUseCaseEntityNotFound.Is(err) {
		return catalog.FactoryErrTrademarkNotFound.Wrap(err, ht.getRawItemID(r))
	}

	if mrcore.FactoryErrUseCaseEntityVersionInvalid.Is(err) {
		return mrerr.NewCustomError("version", err)
	}

	if mrcore.FactoryErrUseCaseSwitchStatusRejected.Is(err) {
		return mrerr.NewCustomError("status", err)
	}

	return err
}
