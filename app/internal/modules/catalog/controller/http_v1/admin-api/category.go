package http_v1

import (
	module "go-sample/internal/modules/catalog"
	"go-sample/internal/modules/catalog/controller/http_v1/admin-api/view"
	view_shared "go-sample/internal/modules/catalog/controller/http_v1/shared/view"
	"go-sample/internal/modules/catalog/entity/admin-api"
	usecase "go-sample/internal/modules/catalog/usecase/admin-api"
	"go-sample/pkg/modules/catalog"
	"net/http"
	"strconv"

	"github.com/mondegor/go-sysmess/mrerr"
	"github.com/mondegor/go-webcore/mrcore"
	"github.com/mondegor/go-webcore/mrctx"
	"github.com/mondegor/go-webcore/mrtype"
	"github.com/mondegor/go-webcore/mrview"
)

const (
	categoryURL             = "/v1/catalog/categories"
	categoryItemURL         = "/v1/catalog/categories/:id"
	categoryChangeStatusURL = "/v1/catalog/categories/:id/status"
)

type (
	Category struct {
		section    mrcore.ClientSection
		service    usecase.CategoryService
		listSorter mrview.ListSorter
	}
)

func NewCategory(
	section mrcore.ClientSection,
	service usecase.CategoryService,
	listSorter mrview.ListSorter,
) *Category {
	return &Category{
		section:    section,
		service:    service,
		listSorter: listSorter,
	}
}

func (ht *Category) AddHandlers(router mrcore.HttpRouter) {
	moduleAccessFunc := func(next mrcore.HttpHandlerFunc) mrcore.HttpHandlerFunc {
		return ht.section.MiddlewareWithPermission(module.UnitCategoryPermission, next)
	}

	router.HttpHandlerFunc(http.MethodGet, ht.section.Path(categoryURL), moduleAccessFunc(ht.GetList()))
	router.HttpHandlerFunc(http.MethodPost, ht.section.Path(categoryURL), moduleAccessFunc(ht.Create()))

	router.HttpHandlerFunc(http.MethodGet, ht.section.Path(categoryItemURL), moduleAccessFunc(ht.Get()))
	router.HttpHandlerFunc(http.MethodPut, ht.section.Path(categoryItemURL), moduleAccessFunc(ht.Store()))
	router.HttpHandlerFunc(http.MethodDelete, ht.section.Path(categoryItemURL), moduleAccessFunc(ht.Remove()))

	router.HttpHandlerFunc(http.MethodPut, ht.section.Path(categoryChangeStatusURL), moduleAccessFunc(ht.ChangeStatus()))
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
			Statuses:   view_shared.ParseFilterStatusList(c, module.ParamNameFilterStatuses),
		},
		Sorter: view_shared.ParseSortParams(c, ht.listSorter),
		Pager:  view_shared.ParsePageParams(c),
	}
}

func (ht *Category) Get() mrcore.HttpHandlerFunc {
	return func(c mrcore.ClientContext) error {
		item, err := ht.service.GetItem(c.Context(), ht.getItemID(c))

		if err != nil {
			return ht.wrapError(err, c)
		}

		return c.SendResponse(http.StatusOK, item)
	}
}

func (ht *Category) Create() mrcore.HttpHandlerFunc {
	return func(c mrcore.ClientContext) error {
		request := view.CreateCategoryRequest{}

		if err := c.Validate(&request); err != nil {
			return err
		}

		item := entity.Category{
			Caption: request.Caption,
		}

		if err := ht.service.Create(c.Context(), &item); err != nil {
			return ht.wrapError(err, c)
		}

		return c.SendResponse(
			http.StatusCreated,
			view.SuccessCreatedItemResponse{
				ItemID: strconv.Itoa(int(item.ID)),
				Message: mrctx.Locale(c.Context()).TranslateMessage(
					"msgCatalogCategorySuccessCreated",
					"entity has been success created",
				),
			},
		)
	}
}

func (ht *Category) Store() mrcore.HttpHandlerFunc {
	return func(c mrcore.ClientContext) error {
		request := view.StoreCategoryRequest{}

		if err := c.Validate(&request); err != nil {
			return err
		}

		item := entity.Category{
			ID:         ht.getItemID(c),
			TagVersion: request.Version,
			Caption:    request.Caption,
		}

		if err := ht.service.Store(c.Context(), &item); err != nil {
			return ht.wrapError(err, c)
		}

		return c.SendResponseNoContent()
	}
}

func (ht *Category) ChangeStatus() mrcore.HttpHandlerFunc {
	return func(c mrcore.ClientContext) error {
		request := view.ChangeItemStatusRequest{}

		if err := c.Validate(&request); err != nil {
			return err
		}

		item := entity.Category{
			ID:         ht.getItemID(c),
			TagVersion: request.TagVersion,
			Status:     request.Status,
		}

		if err := ht.service.ChangeStatus(c.Context(), &item); err != nil {
			return ht.wrapError(err, c)
		}

		return c.SendResponseNoContent()
	}
}

func (ht *Category) Remove() mrcore.HttpHandlerFunc {
	return func(c mrcore.ClientContext) error {
		if err := ht.service.Remove(c.Context(), ht.getItemID(c)); err != nil {
			return ht.wrapError(err, c)
		}

		return c.SendResponseNoContent()
	}
}

func (ht *Category) getItemID(c mrcore.ClientContext) mrtype.KeyInt32 {
	return view_shared.ParseKeyInt32FromPath(c, "id")
}

func (ht *Category) getRawItemID(c mrcore.ClientContext) string {
	return c.ParamFromPath("id")
}

func (ht *Category) wrapError(err error, c mrcore.ClientContext) error {
	if mrcore.FactoryErrServiceEntityNotFound.Is(err) {
		return catalog.FactoryErrCategoryNotFound.Wrap(err, ht.getRawItemID(c))
	}

	if mrcore.FactoryErrServiceEntityVersionInvalid.Is(err) {
		return mrerr.NewFieldError("version", err)
	}

	if mrcore.FactoryErrServiceSwitchStatusRejected.Is(err) {
		return mrerr.NewFieldError("status", err)
	}

	return err
}
