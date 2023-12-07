package http_v1

import (
	module "go-sample/internal/modules/catalog"
	view_shared "go-sample/internal/modules/catalog/controller/http_v1/shared/view"
	usecase "go-sample/internal/modules/catalog/usecase/admin-api"
	usecase_shared "go-sample/internal/modules/catalog/usecase/shared"
	"net/http"

	"github.com/mondegor/go-webcore/mrcore"
	"github.com/mondegor/go-webcore/mrctx"
	"github.com/mondegor/go-webcore/mrdebug"
	"github.com/mondegor/go-webcore/mrtype"
)

const (
	categoryItemImageURL = "/v1/catalog/categories/:id/image"
)

type (
	CategoryImage struct {
		section mrcore.ClientSection
		service usecase.CategoryImageService
	}
)

func NewCategoryImage(
	section mrcore.ClientSection,
	service usecase.CategoryImageService,
) *CategoryImage {
	return &CategoryImage{
		section: section,
		service: service,
	}
}

func (ht *CategoryImage) AddHandlers(router mrcore.HttpRouter) {
	moduleAccessFunc := func(next mrcore.HttpHandlerFunc) mrcore.HttpHandlerFunc {
		return ht.section.MiddlewareWithPermission(module.PermissionCatalogCategory, next)
	}

	router.HttpHandlerFunc(http.MethodGet, ht.section.Path(categoryItemImageURL), moduleAccessFunc(ht.GetImage()))
	router.HttpHandlerFunc(http.MethodPut, ht.section.Path(categoryItemImageURL), moduleAccessFunc(ht.UploadImage()))
	router.HttpHandlerFunc(http.MethodDelete, ht.section.Path(categoryItemImageURL), moduleAccessFunc(ht.RemoveImage()))
}

func (ht *CategoryImage) GetImage() mrcore.HttpHandlerFunc {
	return func(c mrcore.ClientContext) error {
		item, err := ht.service.Get(c.Context(), ht.getItemID(c))

		if err != nil {
			return ht.wrapError(err, ht.getRawItemID(c))
		}

		defer item.Body.Close()

		return c.SendFile(item.FileInfo, "", item.Body)
	}
}

func (ht *CategoryImage) UploadImage() mrcore.HttpHandlerFunc {
	return func(c mrcore.ClientContext) error {
		logger := mrctx.Logger(c.Context())

		file, hdr, err := c.Request().FormFile(module.ParamNameFileCatalogCategoryImage)

		if err != nil {
			mrdebug.MultipartForm(logger, c.Request().MultipartForm)
			return mrcore.FactoryErrHttpMultipartFormFile.Caller(-1).Wrap(err, module.ParamNameFileCatalogCategoryImage)
		}

		defer file.Close()

		mrdebug.MultipartFileHeader(logger, hdr)

		item := mrtype.File{
			FileInfo: mrtype.FileInfo{
				ContentType:  hdr.Header.Get("Content-Type"),
				OriginalName: hdr.Filename,
				Size:         hdr.Size,
			},
			Body: file,
		}

		if err = ht.service.Store(c.Context(), ht.getItemID(c), &item); err != nil {
			return ht.wrapError(err, ht.getRawItemID(c))
		}

		return c.SendResponseNoContent()
	}
}

func (ht *CategoryImage) RemoveImage() mrcore.HttpHandlerFunc {
	return func(c mrcore.ClientContext) error {
		if err := ht.service.Remove(c.Context(), ht.getItemID(c)); err != nil {
			return ht.wrapError(err, ht.getRawItemID(c))
		}

		return c.SendResponseNoContent()
	}
}

func (ht *CategoryImage) getItemID(c mrcore.ClientContext) mrtype.KeyInt32 {
	return view_shared.ParseIDFromPath(c, "id")
}

func (ht *CategoryImage) getRawItemID(c mrcore.ClientContext) string {
	return c.ParamFromPath("id")
}

func (ht *CategoryImage) wrapError(err error, rawItemID string) error {
	if mrcore.FactoryErrServiceEntityNotFound.Is(err) {
		return usecase_shared.FactoryErrCategoryImageNotFound.Wrap(err, rawItemID)
	}

	return err
}
