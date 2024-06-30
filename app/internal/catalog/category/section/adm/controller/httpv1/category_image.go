package httpv1

import (
	"net/http"

	"github.com/mondegor/go-sample/internal/catalog/category/module"
	"github.com/mondegor/go-sample/internal/catalog/category/section/adm/usecase"
	"github.com/mondegor/go-sample/internal/catalog/category/shared/validate"

	"github.com/google/uuid"
	"github.com/mondegor/go-webcore/mrcore"
	"github.com/mondegor/go-webcore/mrserver"
	"github.com/mondegor/go-webcore/mrserver/mrparser"
)

const (
	categoryItemImageURL = "/v1/catalog/categories/{id}/image"
)

type (
	// CategoryImage - comment struct.
	CategoryImage struct {
		parser  validate.RequestCategoryParser
		sender  mrserver.FileResponseSender
		useCase usecase.CategoryImageUseCase
	}
)

// NewCategoryImage - создаёт объект CategoryImage.
func NewCategoryImage(parser validate.RequestCategoryParser, sender mrserver.FileResponseSender, useCase usecase.CategoryImageUseCase) *CategoryImage {
	return &CategoryImage{
		parser:  parser,
		sender:  sender,
		useCase: useCase,
	}
}

// Handlers - comment method.
func (ht *CategoryImage) Handlers() []mrserver.HttpHandler {
	return []mrserver.HttpHandler{
		{Method: http.MethodGet, URL: categoryItemImageURL, Func: ht.GetImage},
		{Method: http.MethodPatch, URL: categoryItemImageURL, Func: ht.UploadImage},
		{Method: http.MethodDelete, URL: categoryItemImageURL, Func: ht.RemoveImage},
	}
}

// GetImage - comment method.
func (ht *CategoryImage) GetImage(w http.ResponseWriter, r *http.Request) error {
	item, err := ht.useCase.GetFile(r.Context(), ht.getParentID(r))
	if err != nil {
		return ht.wrapError(err, r)
	}

	defer item.Body.Close()

	return ht.sender.SendFile(r.Context(), w, item.ToFile())
}

// UploadImage - comment method.
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

// RemoveImage - comment method.
func (ht *CategoryImage) RemoveImage(w http.ResponseWriter, r *http.Request) error {
	if err := ht.useCase.RemoveFile(r.Context(), ht.getParentID(r)); err != nil {
		return ht.wrapError(err, r)
	}

	return ht.sender.SendNoContent(w)
}

func (ht *CategoryImage) getParentID(r *http.Request) uuid.UUID {
	return ht.parser.PathParamUUID(r, "id")
}

func (ht *CategoryImage) getRawParentID(r *http.Request) string {
	return ht.parser.PathParamString(r, "id")
}

func (ht *CategoryImage) wrapError(err error, r *http.Request) error {
	if mrcore.ErrUseCaseEntityNotFound.Is(err) {
		return module.ErrUseCaseCategoryImageNotFound.Wrap(err, ht.getRawParentID(r))
	}

	return err
}
