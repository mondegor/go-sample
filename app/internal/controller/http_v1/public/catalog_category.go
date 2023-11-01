package http_v1

import (
    "go-sample/internal/controller"
    "go-sample/internal/controller/http_v1/public/view"
    view_shared "go-sample/internal/controller/http_v1/shared/view"
    "go-sample/internal/entity/public"
    usecase "go-sample/internal/usecase/public"
    "net/http"

    "github.com/mondegor/go-storage/mrentity"
    "github.com/mondegor/go-webcore/mrcore"
)

const (
    catalogCategoryListURL = "/public/v1/catalog/categories"
    catalogCategoryItemURL = "/public/v1/catalog/categories/:id"
)

type (
    CatalogCategory struct {
        section mrcore.ClientSection
        service usecase.CatalogCategoryService
    }
)

func NewCatalogCategory(
    section mrcore.ClientSection,
    service usecase.CatalogCategoryService,
) *CatalogCategory {
    return &CatalogCategory{
        section: section,
        service: service,
    }
}

func (ht *CatalogCategory) AddHandlers(router mrcore.HttpRouter) {
    moduleAccessFunc := func (next mrcore.HttpHandlerFunc) mrcore.HttpHandlerFunc {
        return ht.section.MiddlewareWithPermission(controller.PermissionCatalogCategory, next)
    }

    router.HttpHandlerFunc(http.MethodGet, catalogCategoryListURL, moduleAccessFunc(ht.GetList()))
    router.HttpHandlerFunc(http.MethodGet, catalogCategoryItemURL, moduleAccessFunc(ht.Get()))
}

func (ht *CatalogCategory) GetList() mrcore.HttpHandlerFunc {
    return func(c mrcore.ClientData) error {
        items, totalItems, err := ht.service.GetList(c.Context(), ht.listParams(c))

        if err != nil {
            return err
        }

        return c.SendResponse(
            http.StatusOK,
            view.CatalogCategoryListResponse{
                Items: items,
                Total: totalItems,
            },
        )
    }
}

func (ht *CatalogCategory) listParams(c mrcore.ClientData) entity.CatalogCategoryParams {
    return entity.CatalogCategoryParams{
        Filter: entity.CatalogCategoryListFilter{
            SearchText: view_shared.ParseFilterString(c, controller.ParamNameFilterQuery),
        },
        Sorter: view_shared.ParseListSorter(c),
        Pager: view_shared.ParseListPager(c),
    }
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

func (ht *CatalogCategory) getItemId(c mrcore.ClientData) mrentity.KeyInt32 {
    return mrentity.KeyInt32(c.RequestPath().GetInt("id"))
}
