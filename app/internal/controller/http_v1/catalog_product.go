package http_v1

import (
    "fmt"
    "go-sample/internal/controller/view"
    "go-sample/internal/entity"
    "go-sample/internal/usecase"
    "net/http"

    "github.com/mondegor/go-storage/mrentity"
    "github.com/mondegor/go-sysmess/mrerr"
    "github.com/mondegor/go-webcore/mrcore"
    "github.com/mondegor/go-webcore/mrctx"
    "github.com/mondegor/go-webcore/mrview"
)

const (
    catalogProductListURL = "/v1/catalog/cat/:cid/products"
    catalogProductItemURL = "/v1/catalog/cat/:cid/products/:id"
    catalogProductChangeStatusURL = "/v1/catalog/cat/:cid/products/:id/status"
    catalogProductMoveURL = "/v1/catalog/cat/:cid/products/:id/move"
)

type (
    CatalogProduct struct {
        service usecase.CatalogProductService
        serviceCategory usecase.CatalogCategoryService
        serviceTrademark usecase.CatalogTrademarkService
    }
)

func NewCatalogProduct(
    service usecase.CatalogProductService,
    serviceCategory usecase.CatalogCategoryService,
    serviceTrademark usecase.CatalogTrademarkService,
) *CatalogProduct {
    return &CatalogProduct{
        service: service,
        serviceCategory: serviceCategory,
        serviceTrademark: serviceTrademark,
    }
}

func (ht *CatalogProduct) AddHandlers(router mrcore.HttpRouter) {
    router.HttpHandlerFunc(http.MethodGet, catalogProductListURL, ht.CategoryMiddleware(ht.GetList()))
    router.HttpHandlerFunc(http.MethodPost, catalogProductListURL, ht.CategoryMiddleware(ht.Create()))

    router.HttpHandlerFunc(http.MethodGet, catalogProductItemURL, ht.CategoryMiddleware(ht.Get()))
    router.HttpHandlerFunc(http.MethodPut, catalogProductItemURL, ht.CategoryMiddleware(ht.Store()))
    router.HttpHandlerFunc(http.MethodDelete, catalogProductItemURL, ht.CategoryMiddleware(ht.Remove()))

    router.HttpHandlerFunc(http.MethodPut, catalogProductChangeStatusURL, ht.CategoryMiddleware(ht.ChangeStatus()))
    router.HttpHandlerFunc(http.MethodPatch, catalogProductMoveURL, ht.CategoryMiddleware(ht.Move()))
}

func (ht *CatalogProduct) GetList() mrcore.HttpHandlerFunc {
    return func(c mrcore.ClientData) error {
        items, err := ht.service.GetList(c.Context(), ht.newListFilter(c))

        if err != nil {
            return err
        }

        return c.SendResponse(http.StatusOK, items)
    }
}

func (ht *CatalogProduct) newListFilter(c mrcore.ClientData) *entity.CatalogProductListFilter {
    var listFilter entity.CatalogProductListFilter

    listFilter.CategoryId = ht.getCategoryId(c)
    view.ParseFilterTrademarkList(c, &listFilter.Trademarks)
    view.ParseFilterItemStatusList(c, &listFilter.Statuses)

    return &listFilter
}

func (ht *CatalogProduct) Get() mrcore.HttpHandlerFunc {
    return func(c mrcore.ClientData) error {
        item, err := ht.service.GetItem(c.Context(), ht.getItemId(c), ht.getCategoryId(c))

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
            CategoryId: ht.getCategoryId(c),
            TrademarkId: request.TrademarkId,
            Article: request.Article,
            Caption: request.Caption,
            Price: request.Price,
        }

        err := ht.service.Create(c.Context(), &item)

        if err != nil {
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

        err := ht.service.Store(c.Context(), &item)

        if err != nil {
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

        err := ht.service.ChangeStatus(c.Context(), &item)

        if err != nil {
            return err
        }

        return c.SendResponseNoContent()
    }
}

func (ht *CatalogProduct) Remove() mrcore.HttpHandlerFunc {
    return func(c mrcore.ClientData) error {
        err := ht.service.Remove(c.Context(), ht.getItemId(c))

        if err != nil {
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
            ht.getCategoryId(c),
        )

        if err != nil {
            return err
        }

        return c.SendResponseNoContent()
    }
}

func (ht *CatalogProduct) getItemId(c mrcore.ClientData) mrentity.KeyInt32 {
    id := mrentity.KeyInt32(c.RequestPath().GetInt("id"))

    if id > 0 {
        return id
    }

    return 0
}
