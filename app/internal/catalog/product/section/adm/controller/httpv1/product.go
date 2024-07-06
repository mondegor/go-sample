package httpv1

import (
	"net/http"

	"github.com/mondegor/go-components/mrsort"

	"github.com/mondegor/go-sample/internal/catalog/product/section/adm"

	"github.com/mondegor/go-sample/pkg/validate"

	"github.com/mondegor/go-sample/internal/catalog/product/module"
	"github.com/mondegor/go-sample/internal/catalog/product/section/adm/entity"
	"github.com/mondegor/go-sample/pkg/catalog/api"
	"github.com/mondegor/go-sample/pkg/view"

	"github.com/mondegor/go-sysmess/mrerr"
	"github.com/mondegor/go-webcore/mrcore"
	"github.com/mondegor/go-webcore/mrserver"
	"github.com/mondegor/go-webcore/mrtype"
	"github.com/mondegor/go-webcore/mrview"
)

const (
	productListURL             = "/v1/catalog/products"
	productItemURL             = "/v1/catalog/products/{id}"
	productItemChangeStatusURL = "/v1/catalog/products/{id}/status"
	productItemMoveURL         = "/v1/catalog/products/{id}/move"
)

type (
	// Product - comment struct.
	Product struct {
		parser     validate.RequestExtendParser
		sender     mrserver.ResponseSender
		useCase    adm.ProductUseCase
		listSorter mrview.ListSorter
	}
)

// NewProduct - создаёт контроллер Product.
func NewProduct(parser validate.RequestExtendParser, sender mrserver.ResponseSender, useCase adm.ProductUseCase, listSorter mrview.ListSorter) *Product {
	return &Product{
		parser:     parser,
		sender:     sender,
		useCase:    useCase,
		listSorter: listSorter,
	}
}

// Handlers - возвращает обработчики контроллера Product.
func (ht *Product) Handlers() []mrserver.HttpHandler {
	return []mrserver.HttpHandler{
		{Method: http.MethodGet, URL: productListURL, Func: ht.GetList},
		{Method: http.MethodPost, URL: productListURL, Func: ht.Create},

		{Method: http.MethodGet, URL: productItemURL, Func: ht.Get},
		{Method: http.MethodPatch, URL: productItemURL, Func: ht.Store},
		{Method: http.MethodDelete, URL: productItemURL, Func: ht.Remove},

		{Method: http.MethodPatch, URL: productItemChangeStatusURL, Func: ht.ChangeStatus},
		{Method: http.MethodPatch, URL: productItemMoveURL, Func: ht.Move},
	}
}

// GetList - comment method.
func (ht *Product) GetList(w http.ResponseWriter, r *http.Request) error {
	items, totalItems, err := ht.useCase.GetList(r.Context(), ht.listParams(r))
	if err != nil {
		return err
	}

	return ht.sender.Send(
		w,
		http.StatusOK,
		ProductListResponse{
			Items: items,
			Total: totalItems,
		},
	)
}

func (ht *Product) listParams(r *http.Request) entity.ProductParams {
	return entity.ProductParams{
		Filter: entity.ProductListFilter{
			CategoryID:   ht.parser.FilterUUID(r, module.ParamNameFilterCatalogCategoryID),
			SearchText:   ht.parser.FilterString(r, module.ParamNameFilterSearchText),
			TrademarkIDs: ht.parser.FilterKeyInt32List(r, module.ParamNameFilterCatalogTrademarkIDs),
			Price:        ht.parser.FilterRangeInt64(r, module.ParamNameFilterPriceRange),
			Statuses:     ht.parser.FilterStatusList(r, module.ParamNameFilterStatuses),
		},
		Sorter: ht.parser.SortParams(r, ht.listSorter),
		Pager:  ht.parser.PageParams(r),
	}
}

// Get - comment method.
func (ht *Product) Get(w http.ResponseWriter, r *http.Request) error {
	item, err := ht.useCase.GetItem(r.Context(), ht.getItemID(r))
	if err != nil {
		return ht.wrapError(err, r)
	}

	return ht.sender.Send(w, http.StatusOK, item)
}

// Create - comment method.
func (ht *Product) Create(w http.ResponseWriter, r *http.Request) error {
	request := CreateProductRequest{}

	if err := ht.parser.Validate(r, &request); err != nil {
		return err
	}

	item := entity.Product{
		CategoryID:  request.CategoryID,
		Article:     request.Article,
		TrademarkID: request.TrademarkID,
		Caption:     request.Caption,
		Price:       &request.Price,
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
func (ht *Product) Store(w http.ResponseWriter, r *http.Request) error {
	request := StoreProductRequest{}

	if err := ht.parser.Validate(r, &request); err != nil {
		return err
	}

	item := entity.Product{
		ID:          ht.getItemID(r),
		TagVersion:  request.TagVersion,
		Article:     request.Article,
		Caption:     request.Caption,
		TrademarkID: request.TrademarkID,
		Price:       request.Price,
	}

	if err := ht.useCase.Store(r.Context(), item); err != nil {
		return ht.wrapError(err, r)
	}

	return ht.sender.SendNoContent(w)
}

// ChangeStatus - comment method.
func (ht *Product) ChangeStatus(w http.ResponseWriter, r *http.Request) error {
	request := view.ChangeItemStatusRequest{}

	if err := ht.parser.Validate(r, &request); err != nil {
		return err
	}

	item := entity.Product{
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
func (ht *Product) Remove(w http.ResponseWriter, r *http.Request) error {
	if err := ht.useCase.Remove(r.Context(), ht.getItemID(r)); err != nil {
		return ht.wrapError(err, r)
	}

	return ht.sender.SendNoContent(w)
}

// Move - comment method.
func (ht *Product) Move(w http.ResponseWriter, r *http.Request) error {
	request := view.MoveItemRequest{}

	if err := ht.parser.Validate(r, &request); err != nil {
		return err
	}

	if err := ht.useCase.MoveAfterID(r.Context(), ht.getItemID(r), request.AfterNodeID); err != nil {
		return ht.wrapErrorNode(err, ht.getRawItemID(r))
	}

	return ht.sender.SendNoContent(w)
}

func (ht *Product) getItemID(r *http.Request) mrtype.KeyInt32 {
	return ht.parser.PathKeyInt32(r, "id")
}

func (ht *Product) getRawItemID(r *http.Request) string {
	return ht.parser.PathParamString(r, "id")
}

func (ht *Product) wrapError(err error, r *http.Request) error {
	if mrcore.ErrUseCaseEntityNotFound.Is(err) {
		return module.ErrUseCaseProductNotFound.Wrap(err, ht.getRawItemID(r))
	}

	if mrcore.ErrUseCaseEntityVersionInvalid.Is(err) {
		return mrerr.NewCustomError("tagVersion", err)
	}

	if mrcore.ErrUseCaseSwitchStatusRejected.Is(err) {
		return mrerr.NewCustomError("status", err)
	}

	if module.ErrUseCaseProductArticleAlreadyExists.Is(err) {
		return mrerr.NewCustomError("article", err)
	}

	if api.ErrCategoryRequired.Is(err) ||
		api.ErrCategoryNotFound.Is(err) {
		return mrerr.NewCustomError("categoryId", err)
	}

	if api.ErrTrademarkRequired.Is(err) ||
		api.ErrTrademarkNotFound.Is(err) {
		return mrerr.NewCustomError("trademarkId", err)
	}

	return err
}

func (ht *Product) wrapErrorNode(err error, rawItemID string) error {
	if mrcore.ErrUseCaseEntityNotFound.Is(err) {
		return module.ErrUseCaseProductNotFound.Wrap(err, rawItemID)
	}

	if mrsort.ErrAfterNodeNotFound.Is(err) {
		return mrerr.NewCustomError("afterNodeId", err)
	}

	return err
}
