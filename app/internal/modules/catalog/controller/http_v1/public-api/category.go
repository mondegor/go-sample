package http_v1

import (
	module "go-sample/internal/modules/catalog"
	"go-sample/internal/modules/catalog/controller/http_v1/public-api/view"
	view_shared "go-sample/internal/modules/catalog/controller/http_v1/shared/view"
	"go-sample/internal/modules/catalog/entity/public-api"
	usecase "go-sample/internal/modules/catalog/usecase/public-api"
	"net/http"

	"github.com/mondegor/go-sysmess/mrlang"
	"github.com/mondegor/go-webcore/mrserver"
	"github.com/mondegor/go-webcore/mrtype"
)

const (
	categoryURL     = "/v1/catalog/categories"
	categoryItemURL = "/v1/catalog/categories/:id"
)

type (
	Category struct {
		parser  view_shared.RequestParser
		sender  mrserver.ResponseSender
		service usecase.CategoryService
	}
)

func NewCategory(
	parser view_shared.RequestParser,
	sender mrserver.ResponseSender,
	service usecase.CategoryService,
) *Category {
	return &Category{
		parser:  parser,
		sender:  sender,
		service: service,
	}
}

func (ht *Category) Handlers() []mrserver.HttpHandler {
	return []mrserver.HttpHandler{
		{http.MethodGet, categoryURL, "", ht.GetList},
		{http.MethodGet, categoryItemURL, "", ht.Get},
	}
}

func (ht *Category) GetList(w http.ResponseWriter, r *http.Request) error {
	items, totalItems, err := ht.service.GetList(r.Context(), ht.listParams(r))

	if err != nil {
		return err
	}

	return ht.sender.Send(
		w,
		http.StatusOK,
		view.CategoryListResponse{
			Items: items,
			Total: totalItems,
		},
	)
}

func (ht *Category) listParams(r *http.Request) entity.CategoryParams {
	return entity.CategoryParams{
		Filter: entity.CategoryListFilter{
			SearchText: ht.parser.FilterString(r, module.ParamNameFilterSearchText),
		},
	}
}

func (ht *Category) Get(w http.ResponseWriter, r *http.Request) error {
	item, err := ht.service.GetItem(
		r.Context(),
		ht.getItemID(r),
		mrlang.Ctx(r.Context()).LangID(),
	)

	if err != nil {
		return err
	}

	return ht.sender.Send(w, http.StatusOK, item)
}

func (ht *Category) getItemID(r *http.Request) mrtype.KeyInt32 {
	return ht.parser.FilterKeyInt32(r, "id")
}
