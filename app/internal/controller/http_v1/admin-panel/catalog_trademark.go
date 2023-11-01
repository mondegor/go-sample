package http_v1

import (
    "fmt"
    "go-sample/internal/controller"
    "go-sample/internal/controller/http_v1/admin-panel/view"
    view_shared "go-sample/internal/controller/http_v1/shared/view"
    entity "go-sample/internal/entity/admin-panel"
    usecase "go-sample/internal/usecase/admin-panel"
    "net/http"

    "github.com/mondegor/go-storage/mrentity"
    "github.com/mondegor/go-webcore/mrcore"
    "github.com/mondegor/go-webcore/mrctx"
    "github.com/mondegor/go-webcore/mrview"
)

const (
    catalogTrademarkListURL = "/adm/v1/catalog/trademarks"
    catalogTrademarkItemURL = "/adm/v1/catalog/trademarks/:id"
    catalogTrademarkChangeStatusURL = "/adm/v1/catalog/trademarks/:id/status"
)

type (
    CatalogTrademark struct {
        section mrcore.ClientSection
        service usecase.CatalogTrademarkService
    }
)

func NewCatalogTrademark(
    section mrcore.ClientSection,
    service usecase.CatalogTrademarkService,
) *CatalogTrademark {
    return &CatalogTrademark{
        section: section,
        service: service,
    }
}

func (ht *CatalogTrademark) AddHandlers(router mrcore.HttpRouter) {
    moduleAccessFunc := func (next mrcore.HttpHandlerFunc) mrcore.HttpHandlerFunc {
        return ht.section.MiddlewareWithPermission(controller.PermissionCatalogTrademark, next)
    }

    router.HttpHandlerFunc(http.MethodGet, catalogTrademarkListURL, moduleAccessFunc(ht.GetList()))
    router.HttpHandlerFunc(http.MethodPost, catalogTrademarkListURL, moduleAccessFunc(ht.Create()))

    router.HttpHandlerFunc(http.MethodGet, catalogTrademarkItemURL, moduleAccessFunc(ht.Get()))
    router.HttpHandlerFunc(http.MethodPut, catalogTrademarkItemURL, moduleAccessFunc(ht.Store()))
    router.HttpHandlerFunc(http.MethodDelete, catalogTrademarkItemURL, moduleAccessFunc(ht.Remove()))

    router.HttpHandlerFunc(http.MethodPut, catalogTrademarkChangeStatusURL, moduleAccessFunc(ht.ChangeStatus()))
}

func (ht *CatalogTrademark) GetList() mrcore.HttpHandlerFunc {
    return func(c mrcore.ClientData) error {
        items, totalItems, err := ht.service.GetList(c.Context(), ht.listParams(c))

        if err != nil {
            return err
        }

        return c.SendResponse(
            http.StatusOK,
            view.CatalogTrademarkListResponse{
                Items: items,
                Total: totalItems,
            },
        )
    }
}

func (ht *CatalogTrademark) listParams(c mrcore.ClientData) entity.CatalogTrademarkParams {
    return entity.CatalogTrademarkParams{
        Filter: entity.CatalogTrademarkListFilter{
            SearchText: view_shared.ParseFilterString(c, controller.ParamNameFilterQuery),
            Statuses: view_shared.ParseFilterItemStatusList(c, controller.ParamNameFilterStatuses),
        },
        Sorter: view_shared.ParseListSorter(c),
        Pager: view_shared.ParseListPager(c),
    }
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

        if err := ht.service.Create(c.Context(), &item); err != nil {
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

        if err := ht.service.Store(c.Context(), &item); err != nil {
            return err
        }

        return c.SendResponseNoContent()
    }
}

func (ht *CatalogTrademark) ChangeStatus() mrcore.HttpHandlerFunc {
    return func(c mrcore.ClientData) error {
        request := view.ChangeItemStatusRequest{}

        if err := c.ParseAndValidate(&request); err != nil {
            return err
        }

        item := entity.CatalogTrademark{
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

func (ht *CatalogTrademark) Remove() mrcore.HttpHandlerFunc {
    return func(c mrcore.ClientData) error {
        if err := ht.service.Remove(c.Context(), ht.getItemId(c)); err != nil {
            return err
        }

        return c.SendResponseNoContent()
    }
}

func (ht *CatalogTrademark) getItemId(c mrcore.ClientData) mrentity.KeyInt32 {
    return mrentity.KeyInt32(c.RequestPath().GetInt("id"))
}
