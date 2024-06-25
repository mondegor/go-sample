package httpv1

import (
	"net/http"

	"github.com/mondegor/go-sample/internal/catalog/category/module"
	"github.com/mondegor/go-sample/internal/catalog/category/section/pub/entity"
	"github.com/mondegor/go-sample/internal/catalog/category/section/pub/usecase"
	"github.com/mondegor/go-sample/internal/catalog/category/shared/validate"

	"github.com/google/uuid"
	"github.com/mondegor/go-webcore/mrserver"
)

const (
	categoryURL     = "/v1/catalog/categories"
	categoryItemURL = "/v1/catalog/categories/{id}"
)

type (
	// Category - comment struct.
	Category struct {
		parser  validate.RequestParser
		sender  mrserver.ResponseSender
		useCase usecase.CategoryUseCase
	}
)

// NewCategory - comment func.
func NewCategory(parser validate.RequestParser, sender mrserver.ResponseSender, useCase usecase.CategoryUseCase) *Category {
	return &Category{
		parser:  parser,
		sender:  sender,
		useCase: useCase,
	}
}

// Handlers - comment method.
func (ht *Category) Handlers() []mrserver.HttpHandler {
	return []mrserver.HttpHandler{
		{Method: http.MethodGet, URL: categoryURL, Func: ht.GetList},
		{Method: http.MethodGet, URL: categoryItemURL, Func: ht.Get},
	}
}

// GetList - comment method.
func (ht *Category) GetList(w http.ResponseWriter, r *http.Request) error {
	items, totalItems, err := ht.useCase.GetList(r.Context(), ht.listParams(r))
	if err != nil {
		return err
	}

	return ht.sender.Send(
		w,
		http.StatusOK,
		CategoryListResponse{
			Items: items,
			Total: totalItems,
		},
	)
}

func (ht *Category) listParams(r *http.Request) entity.CategoryParams {
	return entity.CategoryParams{
		LanguageID: 1, // TODO: mrlang.Ctx(r.Context()).LangID()
		Filter: entity.CategoryListFilter{
			SearchText: ht.parser.FilterString(r, module.ParamNameFilterSearchText),
		},
	}
}

// Get - comment method.
func (ht *Category) Get(w http.ResponseWriter, r *http.Request) error {
	item, err := ht.useCase.GetItem(
		r.Context(),
		ht.getItemID(r),
		1, // TODO: mrlang.Ctx(r.Context()).LangID(),
	)
	if err != nil {
		return err
	}

	return ht.sender.Send(w, http.StatusOK, item)
}

func (ht *Category) getItemID(r *http.Request) uuid.UUID {
	return ht.parser.PathParamUUID(r, "id")
}
