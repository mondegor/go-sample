package http_v1

import (
    "fmt"
    "go-sample/internal/global"
    "go-sample/internal/modules/catalog/controller/http_v1/admin-api/view"
    view_shared "go-sample/internal/modules/catalog/controller/http_v1/shared/view"
    entity "go-sample/internal/modules/catalog/entity/admin-api"
    usecase "go-sample/internal/modules/catalog/usecase/admin-api"
    usecase_shared "go-sample/internal/modules/catalog/usecase/shared"
    "net/http"

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
        section mrcore.ClientSection
        service usecase.ProductService
        listSorter mrview.ListSorter
    }
)

func NewProduct(
    section mrcore.ClientSection,
    service usecase.ProductService,
    listSorter mrview.ListSorter,
) *Product {
    return &Product{
        section: section,
        service: service,
        listSorter: listSorter,
    }
}

func (ht *Product) AddHandlers(router mrcore.HttpRouter) {
    moduleAccessFunc := func (next mrcore.HttpHandlerFunc) mrcore.HttpHandlerFunc {
        return ht.section.MiddlewareWithPermission(global.PermissionCatalogProduct, next)
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
    return func(c mrcore.ClientData) error {
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

func (ht *Product) listParams(c mrcore.ClientData) entity.ProductParams {
    return entity.ProductParams{
        Filter: entity.ProductListFilter{
            CategoryID: view_shared.ParseFilterCategoryID(c, global.ParamNameFilterCategoryID),
            Trademarks: view_shared.ParseFilterTrademarkList(c, global.ParamNameFilterCatalogTrademarks),
            SearchText: view_shared.ParseFilterString(c, global.ParamNameFilterSearchText),
            Price: view_shared.ParseFilterRangeInt64(c, global.ParamNameFilterPriceRange),
            Statuses: view_shared.ParseFilterStatusList(c, global.ParamNameFilterStatuses),
        },
        Sorter: view_shared.ParseListSorter(c, ht.listSorter),
        Pager: view_shared.ParseListPager(c),
    }
}

func (ht *Product) Get() mrcore.HttpHandlerFunc {
    return func(c mrcore.ClientData) error {
        item, err := ht.service.GetItem(c.Context(), ht.getItemID(c))

        if err != nil {
            return err
        }

        return c.SendResponse(http.StatusOK, item)
    }
}

func (ht *Product) Create() mrcore.HttpHandlerFunc {
    return func(c mrcore.ClientData) error {
        request := view.CreateProductRequest{}

        if err := c.ParseAndValidate(&request); err != nil {
            return err
        }

        item := entity.Product{
            CategoryID: request.CategoryID,
            TrademarkID: request.TrademarkID,
            Article: request.Article,
            Caption: request.Caption,
            Price: request.Price,
        }

        if err := ht.service.Create(c.Context(), &item); err != nil {
            return ht.getWrapError(err)
        }

        return c.SendResponse(
            http.StatusCreated,
            view.CreateItemResponse{
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
    return func(c mrcore.ClientData) error {
        request := view.StoreProductRequest{}

        if err := c.ParseAndValidate(&request); err != nil {
            return err
        }

        item := entity.Product{
            ID:          ht.getItemID(c),
            TagVersion:  request.Version,
            CategoryID:  request.CategoryID,
            TrademarkID: request.TrademarkID,
            Article:     request.Article,
            Caption:     request.Caption,
            Price:       request.Price,
        }

        if err := ht.service.Store(c.Context(), &item); err != nil {
            return ht.getWrapError(err)
        }

        return c.SendResponseNoContent()
    }
}

func (ht *Product) ChangeStatus() mrcore.HttpHandlerFunc {
    return func(c mrcore.ClientData) error {
        request := view.ChangeItemStatusRequest{}

        if err := c.ParseAndValidate(&request); err != nil {
            return err
        }

        item := entity.Product{
            ID:         ht.getItemID(c),
            TagVersion: request.Version,
            Status:     request.Status,
        }

        if err := ht.service.ChangeStatus(c.Context(), &item); err != nil {
            return err
        }

        return c.SendResponseNoContent()
    }
}

func (ht *Product) Remove() mrcore.HttpHandlerFunc {
    return func(c mrcore.ClientData) error {
        if err := ht.service.Remove(c.Context(), ht.getItemID(c)); err != nil {
            return err
        }

        return c.SendResponseNoContent()
    }
}

func (ht *Product) Move() mrcore.HttpHandlerFunc {
    return func(c mrcore.ClientData) error {
        request := view.MoveItemRequest{}

        if err := c.ParseAndValidate(&request); err != nil {
            return err
        }

        err := ht.service.MoveAfterID(
            c.Context(),
            ht.getItemID(c),
            request.AfterNodeID,
        )

        if err != nil {
            return err
        }

        return c.SendResponseNoContent()
    }
}

func (ht *Product) getItemID(c mrcore.ClientData) mrtype.KeyInt32 {
    return mrtype.KeyInt32(c.RequestPath().GetInt64("id"))
}

func (ht *Product) getWrapError(err error) error {
    if usecase_shared.FactoryErrProductArticleAlreadyExists.Is(err) {
        return mrerr.NewListWith("article", err)
    }

    if usecase_shared.FactoryErrCategoryNotFound.Is(err) {
        return mrerr.NewListWith("categoryId", err)
    }

    if usecase_shared.FactoryErrTrademarkNotFound.Is(err) {
        return mrerr.NewListWith("trademarkId", err)
    }

    return err
}
