package http_v1

import (
	module "go-sample/internal/modules/catalog/category"
	view_shared "go-sample/internal/modules/catalog/category/controller/http_v1/shared/view"
	"go-sample/internal/modules/catalog/category/entity/public-api"
	usecase "go-sample/internal/modules/catalog/category/usecase/public-api"
	"net/http"

	"github.com/google/uuid"
	"github.com/mondegor/go-webcore/mrserver"
)

const (
	categoryURL     = "/v1/catalog/categories"
	categoryItemURL = "/v1/catalog/categories/:id"
)

type (
	Category struct {
		parser  view_shared.RequestParser
		sender  mrserver.ResponseSender
		useCase usecase.CategoryUseCase
	}
)

func NewCategory(
	parser view_shared.RequestParser,
	sender mrserver.ResponseSender,
	useCase usecase.CategoryUseCase,
) *Category {
	return &Category{
		parser:  parser,
		sender:  sender,
		useCase: useCase,
	}
}

func (ht *Category) Handlers() []mrserver.HttpHandler {
	return []mrserver.HttpHandler{
		{http.MethodGet, categoryURL, "", ht.GetList},
		{http.MethodGet, categoryItemURL, "", ht.Get},
	}
}

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
		LanguageID: 1, // :TODO: mrlang.Ctx(r.Context()).LangID()
		Filter: entity.CategoryListFilter{
			SearchText: ht.parser.FilterString(r, module.ParamNameFilterSearchText),
		},
	}
}

func (ht *Category) Get(w http.ResponseWriter, r *http.Request) error {
	item, err := ht.useCase.GetItem(
		r.Context(),
		ht.getItemID(r),
		1, // :TODO: mrlang.Ctx(r.Context()).LangID(),
	)

	if err != nil {
		return err
	}

	return ht.sender.Send(w, http.StatusOK, item)
}

func (ht *Category) getItemID(r *http.Request) uuid.UUID {
	return ht.parser.PathParamUUID(r, "id")
}
