package http_v1

import (
	module "go-sample/internal/modules/catalog/product"
	view_shared "go-sample/internal/modules/catalog/product/controller/http_v1/shared/view"
	entity "go-sample/internal/modules/catalog/product/entity/admin-api"
	usecase "go-sample/internal/modules/catalog/product/usecase/admin-api"
	usecase_shared "go-sample/internal/modules/catalog/product/usecase/shared"
	"go-sample/pkg/modules/catalog"
	"net/http"
	"strconv"

	"github.com/mondegor/go-components/mrorderer"
	"github.com/mondegor/go-sysmess/mrerr"
	"github.com/mondegor/go-sysmess/mrlang"
	"github.com/mondegor/go-webcore/mrcore"

	"github.com/mondegor/go-webcore/mrserver"
	"github.com/mondegor/go-webcore/mrtype"
	"github.com/mondegor/go-webcore/mrview"
)

const (
	productListURL             = "/v1/catalog/products"
	productItemURL             = "/v1/catalog/products/:id"
	productItemChangeStatusURL = "/v1/catalog/products/:id/status"
	productItemMoveURL         = "/v1/catalog/products/:id/move"
)

type (
	Product struct {
		parser     view_shared.RequestParser
		sender     mrserver.ResponseSender
		service    usecase.ProductService
		listSorter mrview.ListSorter
	}
)

func NewProduct(
	parser view_shared.RequestParser,
	sender mrserver.ResponseSender,
	service usecase.ProductService,
	listSorter mrview.ListSorter,
) *Product {
	return &Product{
		parser:     parser,
		sender:     sender,
		service:    service,
		listSorter: listSorter,
	}
}

func (ht *Product) Handlers() []mrserver.HttpHandler {
	return []mrserver.HttpHandler{
		{http.MethodGet, productListURL, "", ht.GetList},
		{http.MethodPost, productListURL, "", ht.Create},

		{http.MethodGet, productItemURL, "", ht.Get},
		{http.MethodPut, productItemURL, "", ht.Store},
		{http.MethodDelete, productItemURL, "", ht.Remove},

		{http.MethodPut, productItemChangeStatusURL, "", ht.ChangeStatus},
		{http.MethodPatch, productItemMoveURL, "", ht.Move},
	}
}

func (ht *Product) GetList(w http.ResponseWriter, r *http.Request) error {
	items, totalItems, err := ht.service.GetList(r.Context(), ht.listParams(r))

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
			CategoryID:   ht.parser.FilterKeyInt32(r, module.ParamNameFilterCatalogCategoryID),
			SearchText:   ht.parser.FilterString(r, module.ParamNameFilterSearchText),
			TrademarkIDs: ht.parser.FilterKeyInt32List(r, module.ParamNameFilterCatalogTrademarkIDs),
			Price:        ht.parser.FilterRangeInt64(r, module.ParamNameFilterPriceRange),
			Statuses:     ht.parser.FilterStatusList(r, module.ParamNameFilterStatuses),
		},
		Sorter: ht.parser.SortParams(r, ht.listSorter),
		Pager:  ht.parser.PageParams(r),
	}
}

func (ht *Product) Get(w http.ResponseWriter, r *http.Request) error {
	item, err := ht.service.GetItem(r.Context(), ht.getItemID(r))

	if err != nil {
		return ht.wrapError(err, r)
	}

	return ht.sender.Send(w, http.StatusOK, item)
}

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
		Price:       request.Price,
	}

	if err := ht.service.Create(r.Context(), &item); err != nil {
		return ht.wrapError(err, r)
	}

	return ht.sender.Send(
		w,
		http.StatusCreated,
		SuccessCreatedItemResponse{
			ItemID: strconv.Itoa(int(item.ID)),
			Message: mrlang.Ctx(r.Context()).TranslateMessage(
				"msgCatalogProductSuccessCreated",
				"entity has been success created",
			),
		},
	)
}

func (ht *Product) Store(w http.ResponseWriter, r *http.Request) error {
	request := StoreProductRequest{}

	if err := ht.parser.Validate(r, &request); err != nil {
		return err
	}

	item := entity.Product{
		ID:          ht.getItemID(r),
		TagVersion:  request.Version,
		Article:     request.Article,
		Caption:     request.Caption,
		TrademarkID: request.TrademarkID,
		Price:       request.Price,
	}

	if err := ht.service.Store(r.Context(), &item); err != nil {
		return ht.wrapError(err, r)
	}

	return ht.sender.SendNoContent(w)
}

func (ht *Product) ChangeStatus(w http.ResponseWriter, r *http.Request) error {
	request := ChangeItemStatusRequest{}

	if err := ht.parser.Validate(r, &request); err != nil {
		return err
	}

	item := entity.Product{
		ID:         ht.getItemID(r),
		TagVersion: request.TagVersion,
		Status:     request.Status,
	}

	if err := ht.service.ChangeStatus(r.Context(), &item); err != nil {
		return ht.wrapError(err, r)
	}

	return ht.sender.SendNoContent(w)
}

func (ht *Product) Remove(w http.ResponseWriter, r *http.Request) error {
	if err := ht.service.Remove(r.Context(), ht.getItemID(r)); err != nil {
		return ht.wrapError(err, r)
	}

	return ht.sender.SendNoContent(w)
}

func (ht *Product) Move(w http.ResponseWriter, r *http.Request) error {
	request := MoveItemRequest{}

	if err := ht.parser.Validate(r, &request); err != nil {
		return err
	}

	if err := ht.service.MoveAfterID(r.Context(), ht.getItemID(r), request.AfterNodeID); err != nil {
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
	if mrcore.FactoryErrServiceEntityNotFound.Is(err) {
		return usecase_shared.FactoryErrProductNotFound.Wrap(err, ht.getRawItemID(r))
	}

	if mrcore.FactoryErrServiceEntityVersionInvalid.Is(err) {
		return mrerr.NewCustomError("version", err)
	}

	if mrcore.FactoryErrServiceSwitchStatusRejected.Is(err) {
		return mrerr.NewCustomError("status", err)
	}

	if usecase_shared.FactoryErrProductArticleAlreadyExists.Is(err) {
		return mrerr.NewCustomError("article", err)
	}

	if catalog.FactoryErrCategoryNotFound.Is(err) {
		return mrerr.NewCustomError("categoryId", err)
	}

	if catalog.FactoryErrTrademarkNotFound.Is(err) {
		return mrerr.NewCustomError("trademarkId", err)
	}

	return err
}

func (ht *Product) wrapErrorNode(err error, rawItemID string) error {
	if mrcore.FactoryErrServiceEntityNotFound.Is(err) {
		return usecase_shared.FactoryErrProductNotFound.Wrap(err, rawItemID)
	}

	if mrorderer.FactoryErrAfterNodeNotFound.Is(err) {
		return mrerr.NewCustomError("afterNodeId", err)
	}

	return err
}
