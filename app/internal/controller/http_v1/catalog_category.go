package http_v1

import (
    "fmt"
    "go-sample/internal/controller/view"
    "go-sample/internal/entity"
    "go-sample/internal/usecase"
    "net/http"

    "github.com/mondegor/go-storage/mrentity"
    "github.com/mondegor/go-webcore/mrcore"
    "github.com/mondegor/go-webcore/mrctx"
)

const (
    catalogCategoryListURL = "/v1/catalog/categories"
    catalogCategoryItemURL = "/v1/catalog/categories/:id"
    catalogCategoryChangeStatusURL = "/v1/catalog/categories/:id/status"
)

type CatalogCategory struct {
    service usecase.CatalogCategoryService
}

func NewCatalogCategory(service usecase.CatalogCategoryService) *CatalogCategory {
    return &CatalogCategory{
        service: service,
    }
}

func (ht *CatalogCategory) AddHandlers(router mrcore.HttpRouter) {
    router.HttpHandlerFunc(http.MethodGet, catalogCategoryListURL, ht.GetList())
    router.HttpHandlerFunc(http.MethodPost, catalogCategoryListURL, ht.Create())

    router.HttpHandlerFunc(http.MethodGet, catalogCategoryItemURL, ht.Get())
    router.HttpHandlerFunc(http.MethodPut, catalogCategoryItemURL, ht.Store())
    router.HttpHandlerFunc(http.MethodDelete, catalogCategoryItemURL, ht.Remove())

    router.HttpHandlerFunc(http.MethodPut, catalogCategoryChangeStatusURL, ht.ChangeStatus())
}

func (ht *CatalogCategory) GetList() mrcore.HttpHandlerFunc {
    return func(c mrcore.ClientData) error {
        items, err := ht.service.GetList(c.Context(), ht.newListFilter(c))

        if err != nil {
            return err
        }

        return c.SendResponse(http.StatusOK, items)
    }
}

func (ht *CatalogCategory) newListFilter(c mrcore.ClientData) *entity.CatalogCategoryListFilter {
    var listFilter entity.CatalogCategoryListFilter

    parseFilterStatuses(c, &listFilter.Statuses)

    return &listFilter
}

func (ht *CatalogCategory) Get() mrcore.HttpHandlerFunc {
    return func(c mrcore.ClientData) error {
        item, err := ht.service.GetItem(c.Context(), ht.getItemId(c))

        if err != nil {
            return err
        }

        return c.SendResponse(http.StatusOK, item)
    }
}

func (ht *CatalogCategory) Create() mrcore.HttpHandlerFunc {
    return func(c mrcore.ClientData) error {
        request := view.CreateCatalogCategory{}

        if err := c.ParseAndValidate(&request); err != nil {
            return err
        }

        item := entity.CatalogCategory{
            Caption: request.Caption,
        }

        err := ht.service.Create(c.Context(), &item)

        if err != nil {
            return err
        }

        response := view.CreateItemResponse{
            ItemId: fmt.Sprintf("%d", item.Id),
            Message: mrctx.Locale(c.Context()).TranslateMessage(
                "msgCatalogCategorySuccessCreated",
                "entity has been success created",
            ),
        }

        return c.SendResponse(http.StatusCreated, response)
    }
}

func (ht *CatalogCategory) Store() mrcore.HttpHandlerFunc {
    return func(c mrcore.ClientData) error {
        request := view.StoreCatalogCategory{}

        if err := c.ParseAndValidate(&request); err != nil {
            return err
        }

        item := entity.CatalogCategory{
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

func (ht *CatalogCategory) ChangeStatus() mrcore.HttpHandlerFunc {
    return func(c mrcore.ClientData) error {
        request := view.ChangeItemStatus{}

        if err := c.ParseAndValidate(&request); err != nil {
            return err
        }

        item := entity.CatalogCategory{
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

func (ht *CatalogCategory) Remove() mrcore.HttpHandlerFunc {
    return func(c mrcore.ClientData) error {
        err := ht.service.Remove(c.Context(), ht.getItemId(c))

        if err != nil {
            return err
        }

        return c.SendResponseNoContent()
    }
}

func (ht *CatalogCategory) getItemId(c mrcore.ClientData) mrentity.KeyInt32 {
    id := mrentity.KeyInt32(c.RequestPath().GetInt("id"))

    if id > 0 {
        return id
    }

    return 0
}
