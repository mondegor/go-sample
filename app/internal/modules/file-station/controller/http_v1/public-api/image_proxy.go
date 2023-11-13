package http_v1

import (
	"fmt"
	"go-sample/internal/global"
	usecase "go-sample/internal/modules/file-station/usecase/public-api"
	"net/http"
	"strings"

	"github.com/mondegor/go-webcore/mrcore"
)

type (
	ImageProxy struct {
		section mrcore.ClientSection
		service usecase.FileProviderAdapterService
		imagesURL string
	}
)

func NewImageProxy(
	section mrcore.ClientSection,
	service usecase.FileProviderAdapterService,
	basePath string, // :TODO: to URL
) *ImageProxy {
	return &ImageProxy{
		section: section,
		service: service,
		imagesURL: fmt.Sprintf("/%s/*path", strings.Trim(basePath, "/")),
	}
}

func (ht *ImageProxy) AddHandlers(router mrcore.HttpRouter) {
	moduleAccessFunc := func (next mrcore.HttpHandlerFunc) mrcore.HttpHandlerFunc {
		return ht.section.MiddlewareWithPermission(global.PermissionFileStationImageProxy, next)
	}

	router.HttpHandlerFunc(http.MethodGet, ht.section.Path(ht.imagesURL), moduleAccessFunc(ht.Get()))
}

func (ht *ImageProxy) Get() mrcore.HttpHandlerFunc {
	return func(c mrcore.ClientData) error {
		item, err := ht.service.Get(c.Context(), c.RequestPath().Get("path"))

		if err != nil {
			return err
		}

		defer item.Body.Close()

		return c.SendFile(item.FileInfo, "", item.Body)
	}
}
