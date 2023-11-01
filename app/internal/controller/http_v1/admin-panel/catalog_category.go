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
    "github.com/mondegor/go-webcore/mrcore"
    "github.com/mondegor/go-webcore/mrctx"
    "github.com/mondegor/go-webcore/mrview"
)

const (
    catalogCategoryListURL = "/adm/v1/catalog/categories"
    catalogCategoryItemURL = "/adm/v1/catalog/categories/:id"
    catalogCategoryChangeStatusURL = "/adm/v1/catalog/categories/:id/status"
    catalogCategoryItemImageURL = "/adm/v1/catalog/categories/:id/image"
)

type (
    CatalogCategory struct {
        section mrcore.ClientSection
        service usecase.CatalogCategoryService
        serviceImage usecase.CatalogCategoryImageService
    }
)

func NewCatalogCategory(
    section mrcore.ClientSection,
    service usecase.CatalogCategoryService,
    serviceImage usecase.CatalogCategoryImageService,
) *CatalogCategory {
    return &CatalogCategory{
        section: section,
        service: service,
        serviceImage: serviceImage,
    }
}

func (ht *CatalogCategory) AddHandlers(router mrcore.HttpRouter) {
    moduleAccessFunc := func (next mrcore.HttpHandlerFunc) mrcore.HttpHandlerFunc {
        return ht.section.MiddlewareWithPermission(controller.PermissionCatalogCategory, next)
    }

    router.HttpHandlerFunc(http.MethodGet, catalogCategoryListURL, moduleAccessFunc(ht.GetList()))
    router.HttpHandlerFunc(http.MethodPost, catalogCategoryListURL, moduleAccessFunc(ht.Create()))

    router.HttpHandlerFunc(http.MethodGet, catalogCategoryItemURL, moduleAccessFunc(ht.Get()))
    router.HttpHandlerFunc(http.MethodPut, catalogCategoryItemURL, moduleAccessFunc(ht.Store()))
    router.HttpHandlerFunc(http.MethodDelete, catalogCategoryItemURL, moduleAccessFunc(ht.Remove()))

    router.HttpHandlerFunc(http.MethodPut, catalogCategoryChangeStatusURL, moduleAccessFunc(ht.ChangeStatus()))

    router.HttpHandlerFunc(http.MethodGet, catalogCategoryItemImageURL, moduleAccessFunc(ht.GetImage()))
    router.HttpHandlerFunc(http.MethodPut, catalogCategoryItemImageURL, moduleAccessFunc(ht.UploadImage()))
    router.HttpHandlerFunc(http.MethodDelete, catalogCategoryItemImageURL, moduleAccessFunc(ht.RemoveImage()))
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
            Statuses: view_shared.ParseFilterItemStatusList(c, controller.ParamNameFilterStatuses),
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

func (ht *CatalogCategory) Create() mrcore.HttpHandlerFunc {
    return func(c mrcore.ClientData) error {
        request := view.CreateCatalogCategoryRequest{}

        if err := c.ParseAndValidate(&request); err != nil {
            return err
        }

        item := entity.CatalogCategory{
            Caption: request.Caption,
        }

        if err := ht.service.Create(c.Context(), &item); err != nil {
            return err
        }

        response := mrview.CreateItemResponse{
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
        request := view.StoreCatalogCategoryRequest{}

        if err := c.ParseAndValidate(&request); err != nil {
            return err
        }

        item := entity.CatalogCategory{
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

func (ht *CatalogCategory) ChangeStatus() mrcore.HttpHandlerFunc {
    return func(c mrcore.ClientData) error {
        request := view.ChangeItemStatusRequest{}

        if err := c.ParseAndValidate(&request); err != nil {
            return err
        }

        item := entity.CatalogCategory{
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

func (ht *CatalogCategory) Remove() mrcore.HttpHandlerFunc {
    return func(c mrcore.ClientData) error {
        if err := ht.service.Remove(c.Context(), ht.getItemId(c)); err != nil {
            return err
        }

        return c.SendResponseNoContent()
    }
}

func (ht *CatalogCategory) GetImage() mrcore.HttpHandlerFunc {
    return func(c mrcore.ClientData) error {
        item, err := ht.serviceImage.Get(c.Context(), ht.getItemId(c))

        if err != nil {
            return err
        }

        defer item.File.Body.Close()

        return c.SendFile(item.File.ContentType, item.File.Body)
    }
}

func (ht *CatalogCategory) UploadImage() mrcore.HttpHandlerFunc {
    return func(c mrcore.ClientData) error {
        file, hdr, err := c.Request().FormFile("categoryImage")

        if err != nil {
            return err
        }

        defer file.Close()

        logger := mrctx.Logger(c.Context())

        logger.Debug(
            "uploaded file: name=%s, size=%d, header=%#v",
            hdr.Filename, hdr.Size, hdr.Header,
        )

        item := entity.CatalogCategoryImageObject{
            CategoryId: ht.getItemId(c),
            File: mrentity.File{
                ContentType: hdr.Header.Get("Content-Type"),
                Name: hdr.Filename,
                Size: hdr.Size,
                Body: file,
            },
        }

        if err = ht.serviceImage.Store(c.Context(), &item); err != nil {
            return err
        }

        return c.SendResponseNoContent()
    }
}

func (ht *CatalogCategory) RemoveImage() mrcore.HttpHandlerFunc {
    return func(c mrcore.ClientData) error {
        if err := ht.serviceImage.Remove(c.Context(), ht.getItemId(c)); err != nil {
            return err
        }

        return c.SendResponseNoContent()
    }
}

func (ht *CatalogCategory) getItemId(c mrcore.ClientData) mrentity.KeyInt32 {
    return mrentity.KeyInt32(c.RequestPath().GetInt("id"))
}
