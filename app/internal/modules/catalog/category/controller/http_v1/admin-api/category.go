package http_v1

import (
	module "go-sample/internal/modules/catalog/category"
	view_shared "go-sample/internal/modules/catalog/category/controller/http_v1/shared/view"
	"go-sample/internal/modules/catalog/category/entity/admin-api"
	usecase "go-sample/internal/modules/catalog/category/usecase/admin-api"
	"go-sample/pkg/modules/catalog"
	"net/http"
	"strconv"

	"github.com/mondegor/go-sysmess/mrerr"
	"github.com/mondegor/go-sysmess/mrlang"
	"github.com/mondegor/go-webcore/mrcore"

	"github.com/mondegor/go-webcore/mrserver"
	"github.com/mondegor/go-webcore/mrtype"
	"github.com/mondegor/go-webcore/mrview"
)

const (
	categoryListURL             = "/v1/catalog/categories"
	categoryItemURL             = "/v1/catalog/categories/:id"
	categoryItemChangeStatusURL = "/v1/catalog/categories/:id/status"
)

type (
	Category struct {
		parser     view_shared.RequestParser
		sender     mrserver.ResponseSender
		useCase    usecase.CategoryUseCase
		listSorter mrview.ListSorter
	}
)

func NewCategory(
	parser view_shared.RequestParser,
	sender mrserver.ResponseSender,
	useCase usecase.CategoryUseCase,
	listSorter mrview.ListSorter,
) *Category {
	return &Category{
		parser:     parser,
		sender:     sender,
		useCase:    useCase,
		listSorter: listSorter,
	}
}

func (ht *Category) Handlers() []mrserver.HttpHandler {
	return []mrserver.HttpHandler{
		{http.MethodGet, categoryListURL, "", ht.GetList},
		{http.MethodPost, categoryListURL, "", ht.Create},

		{http.MethodGet, categoryItemURL, "", ht.Get},
		{http.MethodPut, categoryItemURL, "", ht.Store},
		{http.MethodDelete, categoryItemURL, "", ht.Remove},

		{http.MethodPut, categoryItemChangeStatusURL, "", ht.ChangeStatus},
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
		Filter: entity.CategoryListFilter{
			SearchText: ht.parser.FilterString(r, module.ParamNameFilterSearchText),
			Statuses:   ht.parser.FilterStatusList(r, module.ParamNameFilterStatuses),
		},
		Sorter: ht.parser.SortParams(r, ht.listSorter),
		Pager:  ht.parser.PageParams(r),
	}
}

func (ht *Category) Get(w http.ResponseWriter, r *http.Request) error {
	item, err := ht.useCase.GetItem(r.Context(), ht.getItemID(r))

	if err != nil {
		return ht.wrapError(err, r)
	}

	return ht.sender.Send(w, http.StatusOK, item)
}

func (ht *Category) Create(w http.ResponseWriter, r *http.Request) error {
	request := CreateCategoryRequest{}

	if err := ht.parser.Validate(r, &request); err != nil {
		return err
	}

	item := entity.Category{
		Caption: request.Caption,
	}

	if itemID, err := ht.useCase.Create(r.Context(), item); err != nil {
		return ht.wrapError(err, r)
	} else {
		return ht.sender.Send(
			w,
			http.StatusCreated,
			SuccessCreatedItemResponse{
				ItemID: strconv.Itoa(int(itemID)),
				Message: mrlang.Ctx(r.Context()).TranslateMessage(
					"msgCatalogCategorySuccessCreated",
					"entity has been success created",
				),
			},
		)
	}
}

func (ht *Category) Store(w http.ResponseWriter, r *http.Request) error {
	request := StoreCategoryRequest{}

	if err := ht.parser.Validate(r, &request); err != nil {
		return err
	}

	item := entity.Category{
		ID:         ht.getItemID(r),
		TagVersion: request.TagVersion,
		Caption:    request.Caption,
	}

	if err := ht.useCase.Store(r.Context(), item); err != nil {
		return ht.wrapError(err, r)
	}

	return ht.sender.SendNoContent(w)
}

func (ht *Category) ChangeStatus(w http.ResponseWriter, r *http.Request) error {
	request := ChangeItemStatusRequest{}

	if err := ht.parser.Validate(r, &request); err != nil {
		return err
	}

	item := entity.Category{
		ID:         ht.getItemID(r),
		TagVersion: request.TagVersion,
		Status:     request.Status,
	}

	if err := ht.useCase.ChangeStatus(r.Context(), item); err != nil {
		return ht.wrapError(err, r)
	}

	return ht.sender.SendNoContent(w)
}

func (ht *Category) Remove(w http.ResponseWriter, r *http.Request) error {
	if err := ht.useCase.Remove(r.Context(), ht.getItemID(r)); err != nil {
		return ht.wrapError(err, r)
	}

	return ht.sender.SendNoContent(w)
}

func (ht *Category) getItemID(r *http.Request) mrtype.KeyInt32 {
	return ht.parser.PathKeyInt32(r, "id")
}

func (ht *Category) getRawItemID(r *http.Request) string {
	return ht.parser.PathParamString(r, "id")
}

func (ht *Category) wrapError(err error, r *http.Request) error {
	if mrcore.FactoryErrUseCaseEntityNotFound.Is(err) {
		return catalog.FactoryErrCategoryNotFound.Wrap(err, ht.getRawItemID(r))
	}

	if mrcore.FactoryErrUseCaseEntityVersionInvalid.Is(err) {
		return mrerr.NewCustomError("tagVersion", err)
	}

	if mrcore.FactoryErrUseCaseSwitchStatusRejected.Is(err) {
		return mrerr.NewCustomError("status", err)
	}

	return err
}
