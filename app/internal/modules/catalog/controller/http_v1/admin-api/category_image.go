package http_v1

import (
	module "go-sample/internal/modules/catalog"
	view_shared "go-sample/internal/modules/catalog/controller/http_v1/shared/view"
	usecase "go-sample/internal/modules/catalog/usecase/admin-api"
	usecase_shared "go-sample/internal/modules/catalog/usecase/shared"
	"net/http"

	"github.com/mondegor/go-webcore/mrcore"
	"github.com/mondegor/go-webcore/mrreq"
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
		return ht.section.MiddlewareWithPermission(module.UnitCategoryPermission, next)
	}

	router.HttpHandlerFunc(http.MethodGet, ht.section.Path(categoryItemImageURL), moduleAccessFunc(ht.GetImage()))
	router.HttpHandlerFunc(http.MethodPut, ht.section.Path(categoryItemImageURL), moduleAccessFunc(ht.UploadImage()))
	router.HttpHandlerFunc(http.MethodDelete, ht.section.Path(categoryItemImageURL), moduleAccessFunc(ht.RemoveImage()))
}

func (ht *CategoryImage) GetImage() mrcore.HttpHandlerFunc {
	return func(c mrcore.ClientContext) error {
		item, err := ht.service.GetFile(c.Context(), ht.getParentID(c))

		if err != nil {
			return ht.wrapError(err, c)
		}

		defer item.Body.Close()

		return c.SendFile(item.FileInfo, "", item.Body)
	}
}

func (ht *CategoryImage) UploadImage() mrcore.HttpHandlerFunc {
	return func(c mrcore.ClientContext) error {
		file, err := mrreq.File(c.Request(), module.ParamNameFileCatalogCategoryImage)

		if err != nil {
			return err
		}

		defer file.Body.Close()

		if err = ht.service.StoreFile(c.Context(), ht.getParentID(c), file); err != nil {
			return ht.wrapError(err, c)
		}

		return c.SendResponseNoContent()
	}
}

func (ht *CategoryImage) RemoveImage() mrcore.HttpHandlerFunc {
	return func(c mrcore.ClientContext) error {
		if err := ht.service.RemoveFile(c.Context(), ht.getParentID(c)); err != nil {
			return ht.wrapError(err, c)
		}

		return c.SendResponseNoContent()
	}
}

func (ht *CategoryImage) getParentID(c mrcore.ClientContext) mrtype.KeyInt32 {
	return view_shared.ParseKeyInt32FromPath(c, "id")
}

func (ht *CategoryImage) getRawParentID(c mrcore.ClientContext) string {
	return c.ParamFromPath("id")
}

func (ht *CategoryImage) wrapError(err error, c mrcore.ClientContext) error {
	if mrcore.FactoryErrServiceEntityNotFound.Is(err) {
		return usecase_shared.FactoryErrCategoryImageNotFound.Wrap(err, ht.getRawParentID(c))
	}

	return err
}
