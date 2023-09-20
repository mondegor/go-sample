package http_v1

import (
    "fmt"
    "go-sample/internal/controller/view"
    "go-sample/internal/entity"
    "go-sample/internal/usecase"
    "net/http"

    "github.com/mondegor/go-components/mrcom"
    "github.com/mondegor/go-storage/mrentity"
    "github.com/mondegor/go-webcore/mrcore"
    "github.com/mondegor/go-webcore/mrctx"
    "github.com/mondegor/go-webcore/mrview"
)

const (
    catalogTrademarkListURL = "/v1/catalog/trademarks"
    catalogTrademarkItemURL = "/v1/catalog/trademarks/:id"
    catalogTrademarkChangeStatusURL = "/v1/catalog/trademarks/:id/status"
)

type (
    CatalogTrademark struct {
        service usecase.CatalogTrademarkService
    }
)

func NewCatalogTrademark(service usecase.CatalogTrademarkService) *CatalogTrademark {
    return &CatalogTrademark{
        service: service,
    }
}

func (ht *CatalogTrademark) AddHandlers(router mrcore.HttpRouter) {
    router.HttpHandlerFunc(http.MethodGet, catalogTrademarkListURL, ht.GetList())
    router.HttpHandlerFunc(http.MethodPost, catalogTrademarkListURL, ht.Create())

    router.HttpHandlerFunc(http.MethodGet, catalogTrademarkItemURL, ht.Get())
    router.HttpHandlerFunc(http.MethodPut, catalogTrademarkItemURL, ht.Store())
    router.HttpHandlerFunc(http.MethodDelete, catalogTrademarkItemURL, ht.Remove())

    router.HttpHandlerFunc(http.MethodPut, catalogTrademarkChangeStatusURL, ht.ChangeStatus())
}

func (ht *CatalogTrademark) GetList() mrcore.HttpHandlerFunc {
    return func(c mrcore.ClientData) error {
        items, err := ht.service.GetList(c.Context(), ht.newListFilter(c))

        if err != nil {
            return err
        }

        return c.SendResponse(http.StatusOK, items)
    }
}

func (ht *CatalogTrademark) newListFilter(c mrcore.ClientData) *entity.CatalogTrademarkListFilter {
    var listFilter entity.CatalogTrademarkListFilter

    parseFilterStatuses(c, &listFilter.Statuses)

    return &listFilter
}

func (ht *CatalogTrademark) Get() mrcore.HttpHandlerFunc {
    return func(c mrcore.ClientData) error {
        item, err := ht.service.GetItem(c.Context(), ht.getItemId(c))

        if err != nil {
            return err
        }

        return c.SendResponse(http.StatusOK, item)
    }
}

func (ht *CatalogTrademark) Create() mrcore.HttpHandlerFunc {
    return func(c mrcore.ClientData) error {
        request := view.CreateCatalogTrademarkRequest{}

        if err := c.ParseAndValidate(&request); err != nil {
            return err
        }

        item := entity.CatalogTrademark{
            Caption: request.Caption,
        }

        err := ht.service.Create(c.Context(), &item)

        if err != nil {
            return err
        }

        response := mrview.CreateItemResponse{
            ItemId: fmt.Sprintf("%d", item.Id),
            Message: mrctx.Locale(c.Context()).TranslateMessage(
                "msgCatalogTrademarkSuccessCreated",
                "entity has been success created",
            ),
        }

        return c.SendResponse(http.StatusCreated, response)
    }
}

func (ht *CatalogTrademark) Store() mrcore.HttpHandlerFunc {
    return func(c mrcore.ClientData) error {
        request := view.StoreCatalogTrademarkRequest{}

        if err := c.ParseAndValidate(&request); err != nil {
            return err
        }

        item := entity.CatalogTrademark{
            Id:      ht.getItemId(c),
            Version: request.Version,
            Caption: request.Caption,
        }

        err := ht.service.Store(c.Context(), &item)

        if err != nil {
            return err
        }

        return c.SendResponseNoContent()
    }
}

func (ht *CatalogTrademark) ChangeStatus() mrcore.HttpHandlerFunc {
    return func(c mrcore.ClientData) error {
        request := mrcom.ChangeItemStatusRequest{}

        if err := c.ParseAndValidate(&request); err != nil {
            return err
        }

        item := entity.CatalogTrademark{
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

func (ht *CatalogTrademark) Remove() mrcore.HttpHandlerFunc {
    return func(c mrcore.ClientData) error {
        err := ht.service.Remove(c.Context(), ht.getItemId(c))

        if err != nil {
            return err
        }

        return c.SendResponseNoContent()
    }
}

func (ht *CatalogTrademark) getItemId(c mrcore.ClientData) mrentity.KeyInt32 {
    id := mrentity.KeyInt32(c.RequestPath().GetInt("id"))

    if id > 0 {
        return id
    }

    return 0
}
