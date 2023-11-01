package http_v1

import (
    "fmt"
    "go-sample/internal/controller"
    "go-sample/internal/controller/http_v1/admin-panel/view"
    view_shared "go-sample/internal/controller/http_v1/shared/view"
    "go-sample/internal/entity/admin-panel"
    usecase "go-sample/internal/usecase/admin-panel"
    "net/http"

    "github.com/mondegor/go-storage/mrentity"
    "github.com/mondegor/go-sysmess/mrerr"
    "github.com/mondegor/go-webcore/mrcore"
    "github.com/mondegor/go-webcore/mrctx"
    "github.com/mondegor/go-webcore/mrview"
)

const (
    catalogProductListURL = "/adm/v1/catalog/products"
    catalogProductItemURL = "/adm/v1/catalog/products/:id"
    catalogProductChangeStatusURL = "/adm/v1/catalog/products/:id/status"
    catalogProductMoveURL = "/adm/v1/catalog/products/:id/move"
)

type (
    CatalogProduct struct {
        section mrcore.ClientSection
        service usecase.CatalogProductService
        serviceCategory usecase.CatalogCategoryService
        serviceTrademark usecase.CatalogTrademarkService
    }
)

func NewCatalogProduct(
    section mrcore.ClientSection,
    service usecase.CatalogProductService,
    serviceCategory usecase.CatalogCategoryService,
    serviceTrademark usecase.CatalogTrademarkService,
) *CatalogProduct {
    return &CatalogProduct{
        section: section,
        service: service,
        serviceCategory: serviceCategory,
        serviceTrademark: serviceTrademark,
    }
}

func (ht *CatalogProduct) AddHandlers(router mrcore.HttpRouter) {
    moduleAccessFunc := func (next mrcore.HttpHandlerFunc) mrcore.HttpHandlerFunc {
        return ht.section.MiddlewareWithPermission(controller.PermissionCatalogProduct, next)
    }

    router.HttpHandlerFunc(http.MethodGet, catalogProductListURL, moduleAccessFunc(ht.GetList()))
    router.HttpHandlerFunc(http.MethodPost, catalogProductListURL, moduleAccessFunc(ht.Create()))

    router.HttpHandlerFunc(http.MethodGet, catalogProductItemURL, moduleAccessFunc(ht.Get()))
    router.HttpHandlerFunc(http.MethodPut, catalogProductItemURL, moduleAccessFunc(ht.Store()))
    router.HttpHandlerFunc(http.MethodDelete, catalogProductItemURL, moduleAccessFunc(ht.Remove()))

    router.HttpHandlerFunc(http.MethodPut, catalogProductChangeStatusURL, moduleAccessFunc(ht.ChangeStatus()))
    router.HttpHandlerFunc(http.MethodPatch, catalogProductMoveURL, moduleAccessFunc(ht.Move()))
}

func (ht *CatalogProduct) GetList() mrcore.HttpHandlerFunc {
    return func(c mrcore.ClientData) error {
        items, totalItems, err := ht.service.GetList(c.Context(), ht.listParams(c))

        if err != nil {
            return err
        }

        return c.SendResponse(
            http.StatusOK,
            view.CatalogProductListResponse{
                Items: items,
                Total: totalItems,
            },
        )
    }
}

func (ht *CatalogProduct) listParams(c mrcore.ClientData) entity.CatalogProductParams {
    return entity.CatalogProductParams{
        Filter: entity.CatalogProductListFilter{
            CategoryId: view_shared.ParseFilterCategoryId(c, controller.ParamNameFilterCategoryId),
            Trademarks: view_shared.ParseFilterTrademarkList(c, controller.ParamNameFilterTrademarks),
            SearchText: view_shared.ParseFilterString(c, controller.ParamNameFilterQuery),
            Price: view_shared.ParseFilterRangeInt64(c, controller.ParamNameFilterPriceRange),
            Statuses: view_shared.ParseFilterItemStatusList(c, controller.ParamNameFilterStatuses),
        },
        Sorter: view_shared.ParseListSorter(c),
        Pager: view_shared.ParseListPager(c),
    }
}

func (ht *CatalogProduct) Get() mrcore.HttpHandlerFunc {
    return func(c mrcore.ClientData) error {
        item, err := ht.service.GetItem(c.Context(), ht.getItemId(c))

        if err != nil {
            return err
        }

        return c.SendResponse(http.StatusOK, item)
    }
}

func (ht *CatalogProduct) Create() mrcore.HttpHandlerFunc {
    return func(c mrcore.ClientData) error {
        request := view.CreateCatalogProductRequest{}

        if err := c.ParseAndValidate(&request); err != nil {
            return err
        }

        item := entity.CatalogProduct{
            CategoryId: request.CategoryId,
            TrademarkId: request.TrademarkId,
            Article: request.Article,
            Caption: request.Caption,
            Price: request.Price,
        }

        if err := ht.service.Create(c.Context(), &item); err != nil {
            if usecase.FactoryErrCatalogProductArticleAlreadyExists.Is(err) {
                return mrerr.NewListWith("article", err)
            }

            if usecase.FactoryErrCatalogTrademarkNotFound.Is(err) {
                return mrerr.NewListWith("trademarkId", err)
            }

            return err
        }

        response := mrview.CreateItemResponse{
            ItemId: fmt.Sprintf("%d", item.Id),
            Message: mrctx.Locale(c.Context()).TranslateMessage(
                "msgCatalogProductSuccessCreated",
                "entity has been success created",
            ),
        }

        return c.SendResponse(http.StatusCreated, response)
    }
}

func (ht *CatalogProduct) Store() mrcore.HttpHandlerFunc {
    return func(c mrcore.ClientData) error {
        request := view.StoreCatalogProductRequest{}

        if err := c.ParseAndValidate(&request); err != nil {
            return err
        }

        item := entity.CatalogProduct{
            Id:          ht.getItemId(c),
            Version:     request.Version,
            TrademarkId: request.TrademarkId,
            Article:     request.Article,
            Caption:     request.Caption,
            Price:       request.Price,
        }

        if err := ht.service.Store(c.Context(), &item); err != nil {
            if usecase.FactoryErrCatalogTrademarkNotFound.Is(err) {
                return mrerr.NewListWith("trademarkId", err)
            }

            return err
        }

        return c.SendResponseNoContent()
    }
}

func (ht *CatalogProduct) ChangeStatus() mrcore.HttpHandlerFunc {
    return func(c mrcore.ClientData) error {
        request := view.ChangeItemStatusRequest{}

        if err := c.ParseAndValidate(&request); err != nil {
            return err
        }

        item := entity.CatalogProduct{
            Id:      ht.getItemId(c),
            Version: request.Version,
            Status:  request.Status,
        }

        if err := ht.service.ChangeStatus(c.Context(), &item); err != nil {
            return err
        }

        return c.SendResponseNoContent()
    }
}

func (ht *CatalogProduct) Remove() mrcore.HttpHandlerFunc {
    return func(c mrcore.ClientData) error {
        if err := ht.service.Remove(c.Context(), ht.getItemId(c)); err != nil {
            return err
        }

        return c.SendResponseNoContent()
    }
}

func (ht *CatalogProduct) Move() mrcore.HttpHandlerFunc {
    return func(c mrcore.ClientData) error {
        request := view.MoveCatalogProductRequest{}

        if err := c.ParseAndValidate(&request); err != nil {
            return err
        }

        err := ht.service.MoveAfterId(
            c.Context(),
            ht.getItemId(c),
            request.AfterNodeId,
        )

        if err != nil {
            return err
        }

        return c.SendResponseNoContent()
    }
}

func (ht *CatalogProduct) getItemId(c mrcore.ClientData) mrentity.KeyInt32 {
    return mrentity.KeyInt32(c.RequestPath().GetInt("id"))
}
