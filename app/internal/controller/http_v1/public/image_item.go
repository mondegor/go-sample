package http_v1

import (
    "go-sample/internal/controller"
    usecase "go-sample/internal/usecase/public"
    "net/http"

    "github.com/mondegor/go-webcore/mrcore"
)

const (
    imageItemURL = "/public/v1/images/*path"
)

type (
    Image struct {
        section mrcore.ClientSection
        serviceFile usecase.FileItemService
    }
)

func NewImageItem(
    section mrcore.ClientSection,
    serviceFile usecase.FileItemService,
) *Image {
    return &Image{
        section: section,
        serviceFile: serviceFile,
    }
}

func (ht *Image) AddHandlers(router mrcore.HttpRouter) {
    moduleAccessFunc := func (next mrcore.HttpHandlerFunc) mrcore.HttpHandlerFunc {
        return ht.section.MiddlewareWithPermission(controller.PermissionImageItem, next)
    }

    router.HttpHandlerFunc(http.MethodGet, ht.section.Path(imageItemURL), moduleAccessFunc(ht.GetImage()))
}

func (ht *Image) GetImage() mrcore.HttpHandlerFunc {
    return func(c mrcore.ClientData) error {
        item, err := ht.serviceFile.Get(c.Context(), c.RequestPath().Get("path"))

        if err != nil {
            return err
        }

        defer item.Body.Close()

        return c.SendFile(item.ContentType, item.Body)
    }
}
