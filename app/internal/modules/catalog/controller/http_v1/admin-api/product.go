package http_v1

import (
	"fmt"
	module "go-sample/internal/modules/catalog"
	"go-sample/internal/modules/catalog/controller/http_v1/admin-api/view"
	view_shared "go-sample/internal/modules/catalog/controller/http_v1/shared/view"
	entity "go-sample/internal/modules/catalog/entity/admin-api"
	usecase "go-sample/internal/modules/catalog/usecase/admin-api"
	usecase_shared "go-sample/internal/modules/catalog/usecase/shared"
	"net/http"

	"github.com/mondegor/go-components/mrorderer"
	"github.com/mondegor/go-sysmess/mrerr"
	"github.com/mondegor/go-webcore/mrcore"
	"github.com/mondegor/go-webcore/mrctx"
	"github.com/mondegor/go-webcore/mrtype"
	"github.com/mondegor/go-webcore/mrview"
)

const (
	productURL             = "/v1/catalog/products"
	productItemURL         = "/v1/catalog/products/:id"
	productChangeStatusURL = "/v1/catalog/products/:id/status"
	productMoveURL         = "/v1/catalog/products/:id/move"
)

type (
	Product struct {
		section    mrcore.ClientSection
		service    usecase.ProductService
		listSorter mrview.ListSorter
	}
)

func NewProduct(
	section mrcore.ClientSection,
	service usecase.ProductService,
	listSorter mrview.ListSorter,
) *Product {
	return &Product{
		section:    section,
		service:    service,
		listSorter: listSorter,
	}
}

func (ht *Product) AddHandlers(router mrcore.HttpRouter) {
	moduleAccessFunc := func(next mrcore.HttpHandlerFunc) mrcore.HttpHandlerFunc {
		return ht.section.MiddlewareWithPermission(module.PermissionCatalogProduct, next)
	}

	router.HttpHandlerFunc(http.MethodGet, ht.section.Path(productURL), moduleAccessFunc(ht.GetList()))
	router.HttpHandlerFunc(http.MethodPost, ht.section.Path(productURL), moduleAccessFunc(ht.Create()))

	router.HttpHandlerFunc(http.MethodGet, ht.section.Path(productItemURL), moduleAccessFunc(ht.Get()))
	router.HttpHandlerFunc(http.MethodPut, ht.section.Path(productItemURL), moduleAccessFunc(ht.Store()))
	router.HttpHandlerFunc(http.MethodDelete, ht.section.Path(productItemURL), moduleAccessFunc(ht.Remove()))

	router.HttpHandlerFunc(http.MethodPut, ht.section.Path(productChangeStatusURL), moduleAccessFunc(ht.ChangeStatus()))
	router.HttpHandlerFunc(http.MethodPatch, ht.section.Path(productMoveURL), moduleAccessFunc(ht.Move()))
}

func (ht *Product) GetList() mrcore.HttpHandlerFunc {
	return func(c mrcore.ClientContext) error {
		items, totalItems, err := ht.service.GetList(c.Context(), ht.listParams(c))

		if err != nil {
			return err
		}

		return c.SendResponse(
			http.StatusOK,
			view.ProductListResponse{
				Items: items,
				Total: totalItems,
			},
		)
	}
}

func (ht *Product) listParams(c mrcore.ClientContext) entity.ProductParams {
	return entity.ProductParams{
		Filter: entity.ProductListFilter{
			CategoryID: view_shared.ParseFilterCategoryID(c, module.ParamNameFilterCategoryID),
			Trademarks: view_shared.ParseFilterTrademarkList(c, module.ParamNameFilterCatalogTrademarks),
			SearchText: view_shared.ParseFilterString(c, module.ParamNameFilterSearchText),
			Price:      view_shared.ParseFilterRangeInt64(c, module.ParamNameFilterPriceRange),
			Statuses:   view_shared.ParseFilterStatusList(c, module.ParamNameFilterStatuses),
		},
		Sorter: view_shared.ParseSortParams(c, ht.listSorter),
		Pager:  view_shared.ParsePageParams(c),
	}
}

func (ht *Product) Get() mrcore.HttpHandlerFunc {
	return func(c mrcore.ClientContext) error {
		item, err := ht.service.GetItem(c.Context(), ht.getItemID(c))

		if err != nil {
			return ht.wrapError(err, ht.getRawItemID(c))
		}

		return c.SendResponse(http.StatusOK, item)
	}
}

func (ht *Product) Create() mrcore.HttpHandlerFunc {
	return func(c mrcore.ClientContext) error {
		request := view.CreateProductRequest{}

		if err := c.Validate(&request); err != nil {
			return err
		}

		item := entity.Product{
			CategoryID:  request.CategoryID,
			TrademarkID: request.TrademarkID,
			Article:     request.Article,
			Caption:     request.Caption,
			Price:       request.Price,
		}

		if err := ht.service.Create(c.Context(), &item); err != nil {
			return ht.wrapError(err, ht.getRawItemID(c))
		}

		return c.SendResponse(
			http.StatusCreated,
			view.SuccessCreatedItemResponse{
				ItemID: fmt.Sprintf("%d", item.ID),
				Message: mrctx.Locale(c.Context()).TranslateMessage(
					"msgProductSuccessCreated",
					"entity has been success created",
				),
			},
		)
	}
}

func (ht *Product) Store() mrcore.HttpHandlerFunc {
	return func(c mrcore.ClientContext) error {
		request := view.StoreProductRequest{}

		if err := c.Validate(&request); err != nil {
			return err
		}

		item := entity.Product{
			ID:          ht.getItemID(c),
			TagVersion:  request.Version,
			TrademarkID: request.TrademarkID,
			Article:     request.Article,
			Caption:     request.Caption,
			Price:       request.Price,
		}

		if err := ht.service.Store(c.Context(), &item); err != nil {
			return ht.wrapError(err, ht.getRawItemID(c))
		}

		return c.SendResponseNoContent()
	}
}

func (ht *Product) ChangeStatus() mrcore.HttpHandlerFunc {
	return func(c mrcore.ClientContext) error {
		request := view.ChangeItemStatusRequest{}

		if err := c.Validate(&request); err != nil {
			return err
		}

		item := entity.Product{
			ID:         ht.getItemID(c),
			TagVersion: request.TagVersion,
			Status:     request.Status,
		}

		if err := ht.service.ChangeStatus(c.Context(), &item); err != nil {
			return ht.wrapError(err, ht.getRawItemID(c))
		}

		return c.SendResponseNoContent()
	}
}

func (ht *Product) Remove() mrcore.HttpHandlerFunc {
	return func(c mrcore.ClientContext) error {
		if err := ht.service.Remove(c.Context(), ht.getItemID(c)); err != nil {
			return ht.wrapError(err, ht.getRawItemID(c))
		}

		return c.SendResponseNoContent()
	}
}

func (ht *Product) Move() mrcore.HttpHandlerFunc {
	return func(c mrcore.ClientContext) error {
		request := view.MoveItemRequest{}

		if err := c.Validate(&request); err != nil {
			return err
		}

		err := ht.service.MoveAfterID(
			c.Context(),
			ht.getItemID(c),
			request.AfterNodeID,
		)

		if err != nil {
			return ht.wrapErrorNode(err, ht.getRawItemID(c))
		}

		return c.SendResponseNoContent()
	}
}

func (ht *Product) getItemID(c mrcore.ClientContext) mrtype.KeyInt32 {
	return view_shared.ParseIDFromPath(c, "id")
}

func (ht *Product) getRawItemID(c mrcore.ClientContext) string {
	return c.ParamFromPath("id")
}

func (ht *Product) wrapError(err error, rawItemID string) error {
	if mrcore.FactoryErrServiceEntityNotFound.Is(err) {
		return usecase_shared.FactoryErrProductNotFound.Wrap(err, rawItemID)
	}

	if mrcore.FactoryErrServiceEntityVersionInvalid.Is(err) {
		return mrerr.NewFieldError("version", err)
	}

	if mrcore.FactoryErrServiceSwitchStatusRejected.Is(err) {
		return mrerr.NewFieldError("status", err)
	}

	if usecase_shared.FactoryErrProductArticleAlreadyExists.Is(err) {
		return mrerr.NewFieldError("article", err)
	}

	if usecase_shared.FactoryErrCategoryNotFound.Is(err) {
		return mrerr.NewFieldError("categoryId", err)
	}

	if usecase_shared.FactoryErrTrademarkNotFound.Is(err) {
		return mrerr.NewFieldError("trademarkId", err)
	}

	return err
}

func (ht *Product) wrapErrorNode(err error, rawItemID string) error {
	if mrcore.FactoryErrServiceEntityNotFound.Is(err) {
		return usecase_shared.FactoryErrProductNotFound.Wrap(err, rawItemID)
	}

	if mrorderer.FactoryErrAfterNodeNotFound.Is(err) {
		return mrerr.NewFieldError("afterNodeId", err)
	}

	return err
}
