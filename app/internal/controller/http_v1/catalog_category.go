package http_v1

import (
    "fmt"
    "go-sample/internal/controller/view"
    "go-sample/internal/entity"
    "go-sample/internal/usecase"
    "net/http"

    "github.com/mondegor/go-components/mrcom"
    "github.com/mondegor/go-storage/mrentity"
    "github.com/mondegor/go-storage/mrstorage"
    "github.com/mondegor/go-webcore/mrcore"
    "github.com/mondegor/go-webcore/mrctx"
    "github.com/mondegor/go-webcore/mrview"
)

const (
    catalogCategoryListURL = "/v1/catalog/categories"
    catalogCategoryItemURL = "/v1/catalog/categories/:id"
    catalogCategoryChangeStatusURL = "/v1/catalog/categories/:id/status"
    catalogCategoryItemImageURL = "/v1/catalog/categories/:id/image"
)

type (
    CatalogCategory struct {
        service usecase.CatalogCategoryService
        serviceImage usecase.CatalogCategoryImageService
    }
)

func NewCatalogCategory(service usecase.CatalogCategoryService,
                        serviceImage usecase.CatalogCategoryImageService) *CatalogCategory {
    return &CatalogCategory{
        service: service,
        serviceImage: serviceImage,
    }
}

func (ht *CatalogCategory) AddHandlers(router mrcore.HttpRouter) {
    router.HttpHandlerFunc(http.MethodGet, catalogCategoryListURL, ht.GetList())
    router.HttpHandlerFunc(http.MethodPost, catalogCategoryListURL, ht.Create())

    router.HttpHandlerFunc(http.MethodGet, catalogCategoryItemURL, ht.Get())
    router.HttpHandlerFunc(http.MethodPut, catalogCategoryItemURL, ht.Store())
    router.HttpHandlerFunc(http.MethodDelete, catalogCategoryItemURL, ht.Remove())

    router.HttpHandlerFunc(http.MethodPut, catalogCategoryChangeStatusURL, ht.ChangeStatus())

    router.HttpHandlerFunc(http.MethodGet, catalogCategoryItemImageURL, ht.GetImage())
    router.HttpHandlerFunc(http.MethodPatch, catalogCategoryItemImageURL, ht.UploadImage())
    router.HttpHandlerFunc(http.MethodDelete, catalogCategoryItemImageURL, ht.RemoveImage())
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
        request := mrcom.ChangeItemStatusRequest{}

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

func (ht *CatalogCategory) GetImage() mrcore.HttpHandlerFunc {
    return func(c mrcore.ClientData) error {
        item := entity.CatalogCategoryImageObject{
            CategoryId: ht.getItemId(c),
            File: mrstorage.File{},
        }

        err := ht.serviceImage.Load(c.Context(), &item)

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
            File: mrstorage.File{
                ContentType: hdr.Header.Get("Content-Type"),
                Name: hdr.Filename,
                Size: hdr.Size,
                Body: file,
            },
        }

        err = ht.serviceImage.Store(c.Context(), &item)

        if err != nil {
            return err
        }

        return c.SendResponseNoContent()
    }
}

func (ht *CatalogCategory) RemoveImage() mrcore.HttpHandlerFunc {
    return func(c mrcore.ClientData) error {
        err := ht.serviceImage.Remove(c.Context(), ht.getItemId(c))

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
