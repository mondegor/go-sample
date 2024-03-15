package http_v1

import (
	module "go-sample/internal/modules/catalog/category"
	view_shared "go-sample/internal/modules/catalog/category/controller/http_v1/shared/view"
	usecase "go-sample/internal/modules/catalog/category/usecase/admin-api"
	usecase_shared "go-sample/internal/modules/catalog/category/usecase/shared"
	"net/http"

	"github.com/mondegor/go-webcore/mrcore"
	"github.com/mondegor/go-webcore/mrserver"
	"github.com/mondegor/go-webcore/mrserver/mrparser"
	"github.com/mondegor/go-webcore/mrtype"
)

const (
	categoryItemImageURL = "/v1/catalog/categories/:id/image"
)

type (
	CategoryImage struct {
		parser  view_shared.RequestParser
		sender  mrserver.FileResponseSender
		useCase usecase.CategoryImageUseCase
	}
)

func NewCategoryImage(
	parser view_shared.RequestParser,
	sender mrserver.FileResponseSender,
	useCase usecase.CategoryImageUseCase,
) *CategoryImage {
	return &CategoryImage{
		parser:  parser,
		sender:  sender,
		useCase: useCase,
	}
}

func (ht *CategoryImage) Handlers() []mrserver.HttpHandler {
	return []mrserver.HttpHandler{
		{http.MethodGet, categoryItemImageURL, "", ht.GetImage},
		{http.MethodPut, categoryItemImageURL, "", ht.UploadImage},
		{http.MethodDelete, categoryItemImageURL, "", ht.RemoveImage},
	}
}

func (ht *CategoryImage) GetImage(w http.ResponseWriter, r *http.Request) error {
	item, err := ht.useCase.GetFile(r.Context(), ht.getParentID(r))

	if err != nil {
		return ht.wrapError(err, r)
	}

	defer item.Body.Close()

	return ht.sender.SendFile(r.Context(), w, item.ToFile())
}

func (ht *CategoryImage) UploadImage(w http.ResponseWriter, r *http.Request) error {
	image, err := ht.parser.FormImage(r, module.ParamNameFileCatalogCategoryImage)

	if err != nil {
		return mrparser.WrapImageError(err, module.ParamNameFileCatalogCategoryImage)
	}

	defer image.Body.Close()

	if err = ht.useCase.StoreFile(r.Context(), ht.getParentID(r), image); err != nil {
		return ht.wrapError(err, r)
	}

	return ht.sender.SendNoContent(w)
}

func (ht *CategoryImage) RemoveImage(w http.ResponseWriter, r *http.Request) error {
	if err := ht.useCase.RemoveFile(r.Context(), ht.getParentID(r)); err != nil {
		return ht.wrapError(err, r)
	}

	return ht.sender.SendNoContent(w)
}

func (ht *CategoryImage) getParentID(r *http.Request) mrtype.KeyInt32 {
	return ht.parser.PathKeyInt32(r, "id")
}

func (ht *CategoryImage) getRawParentID(r *http.Request) string {
	return ht.parser.PathParamString(r, "id")
}

func (ht *CategoryImage) wrapError(err error, r *http.Request) error {
	if mrcore.FactoryErrUseCaseEntityNotFound.Is(err) {
		return usecase_shared.FactoryErrCategoryImageNotFound.Wrap(err, ht.getRawParentID(r))
	}

	return err
}
