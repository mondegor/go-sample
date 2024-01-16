package http_v1

import (
	module "go-sample/internal/modules/catalog"
	"go-sample/internal/modules/catalog/controller/http_v1/public-api/view"
	view_shared "go-sample/internal/modules/catalog/controller/http_v1/shared/view"
	"go-sample/internal/modules/catalog/entity/public-api"
	usecase "go-sample/internal/modules/catalog/usecase/public-api"
	"net/http"

	"github.com/mondegor/go-webcore/mrcore"
	"github.com/mondegor/go-webcore/mrtype"
)

const (
	categoryURL     = "/v1/catalog/categories"
	categoryItemURL = "/v1/catalog/categories/:id"
)

type (
	Category struct {
		section mrcore.ClientSection
		service usecase.CategoryService
	}
)

func NewCategory(
	section mrcore.ClientSection,
	service usecase.CategoryService,
) *Category {
	return &Category{
		section: section,
		service: service,
	}
}

func (ht *Category) AddHandlers(router mrcore.HttpRouter) {
	moduleAccessFunc := func(next mrcore.HttpHandlerFunc) mrcore.HttpHandlerFunc {
		return ht.section.MiddlewareWithPermission(module.UnitCategoryPermission, next)
	}

	router.HttpHandlerFunc(http.MethodGet, ht.section.Path(categoryURL), moduleAccessFunc(ht.GetList()))
	router.HttpHandlerFunc(http.MethodGet, ht.section.Path(categoryItemURL), moduleAccessFunc(ht.Get()))
}

func (ht *Category) GetList() mrcore.HttpHandlerFunc {
	return func(c mrcore.ClientContext) error {
		items, totalItems, err := ht.service.GetList(c.Context(), ht.listParams(c))

		if err != nil {
			return err
		}

		return c.SendResponse(
			http.StatusOK,
			view.CategoryListResponse{
				Items: items,
				Total: totalItems,
			},
		)
	}
}

func (ht *Category) listParams(c mrcore.ClientContext) entity.CategoryParams {
	return entity.CategoryParams{
		Filter: entity.CategoryListFilter{
			SearchText: view_shared.ParseFilterString(c, module.ParamNameFilterSearchText),
		},
	}
}

func (ht *Category) Get() mrcore.HttpHandlerFunc {
	return func(c mrcore.ClientContext) error {
		item, err := ht.service.GetItem(c.Context(), ht.getItemID(c))

		if err != nil {
			return err
		}

		return c.SendResponse(http.StatusOK, item)
	}
}

func (ht *Category) getItemID(c mrcore.ClientContext) mrtype.KeyInt32 {
	return view_shared.ParseKeyInt32FromPath(c, "id")
}
