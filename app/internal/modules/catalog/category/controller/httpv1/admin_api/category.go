package httpv1

import (
	"net/http"

	view_shared "go-sample/internal/modules/catalog/category/controller/httpv1/shared/view"
	entity "go-sample/internal/modules/catalog/category/entity/admin_api"
	"go-sample/internal/modules/catalog/category/module"
	usecase "go-sample/internal/modules/catalog/category/usecase/admin_api"
	"go-sample/pkg/modules/catalog/api"
	"go-sample/pkg/shared/view"

	"github.com/google/uuid"
	"github.com/mondegor/go-sysmess/mrerr"
	"github.com/mondegor/go-webcore/mrcore"
	"github.com/mondegor/go-webcore/mrserver"
	"github.com/mondegor/go-webcore/mrview"
)

const (
	categoryListURL             = "/v1/catalog/categories"
	categoryItemURL             = "/v1/catalog/categories/{id}"
	categoryItemChangeStatusURL = "/v1/catalog/categories/{id}/status"
)

type (
	// Category - comment struct.
	Category struct {
		parser     view_shared.RequestParser
		sender     mrserver.ResponseSender
		useCase    usecase.CategoryUseCase
		listSorter mrview.ListSorter
	}
)

// NewCategory - comment func.
func NewCategory(parser view_shared.RequestParser, sender mrserver.ResponseSender, useCase usecase.CategoryUseCase, listSorter mrview.ListSorter) *Category {
	return &Category{
		parser:     parser,
		sender:     sender,
		useCase:    useCase,
		listSorter: listSorter,
	}
}

// Handlers - comment method.
func (ht *Category) Handlers() []mrserver.HttpHandler {
	return []mrserver.HttpHandler{
		{Method: http.MethodGet, URL: categoryListURL, Func: ht.GetList},
		{Method: http.MethodPost, URL: categoryListURL, Func: ht.Create},

		{Method: http.MethodGet, URL: categoryItemURL, Func: ht.Get},
		{Method: http.MethodPut, URL: categoryItemURL, Func: ht.Store},
		{Method: http.MethodDelete, URL: categoryItemURL, Func: ht.Remove},

		{Method: http.MethodPatch, URL: categoryItemChangeStatusURL, Func: ht.ChangeStatus},
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
		Filter: entity.CategoryListFilter{
			SearchText: ht.parser.FilterString(r, module.ParamNameFilterSearchText),
			Statuses:   ht.parser.FilterStatusList(r, module.ParamNameFilterStatuses),
		},
		Sorter: ht.parser.SortParams(r, ht.listSorter),
		Pager:  ht.parser.PageParams(r),
	}
}

// Get - comment method.
func (ht *Category) Get(w http.ResponseWriter, r *http.Request) error {
	item, err := ht.useCase.GetItem(r.Context(), ht.getItemID(r))
	if err != nil {
		return ht.wrapError(err, r)
	}

	return ht.sender.Send(w, http.StatusOK, item)
}

// Create - comment method.
func (ht *Category) Create(w http.ResponseWriter, r *http.Request) error {
	request := CreateCategoryRequest{}

	if err := ht.parser.Validate(r, &request); err != nil {
		return err
	}

	item := entity.Category{
		Caption: request.Caption,
	}

	itemID, err := ht.useCase.Create(r.Context(), item)
	if err != nil {
		return ht.wrapError(err, r)
	}

	return ht.sender.Send(
		w,
		http.StatusCreated,
		view.SuccessCreatedItemResponse{
			ItemID: itemID.String(),
		},
	)
}

// Store - comment method.
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

// ChangeStatus - comment method.
func (ht *Category) ChangeStatus(w http.ResponseWriter, r *http.Request) error {
	request := view.ChangeItemStatusRequest{}

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

// Remove - comment method.
func (ht *Category) Remove(w http.ResponseWriter, r *http.Request) error {
	if err := ht.useCase.Remove(r.Context(), ht.getItemID(r)); err != nil {
		return ht.wrapError(err, r)
	}

	return ht.sender.SendNoContent(w)
}

func (ht *Category) getItemID(r *http.Request) uuid.UUID {
	return ht.parser.PathParamUUID(r, "id")
}

func (ht *Category) getRawItemID(r *http.Request) string {
	return ht.parser.PathParamString(r, "id")
}

func (ht *Category) wrapError(err error, r *http.Request) error {
	if mrcore.ErrUseCaseEntityNotFound.Is(err) {
		return api.ErrCategoryNotFound.Wrap(err, ht.getRawItemID(r))
	}

	if mrcore.ErrUseCaseEntityVersionInvalid.Is(err) {
		return mrerr.NewCustomError("tagVersion", err)
	}

	if mrcore.ErrUseCaseSwitchStatusRejected.Is(err) {
		return mrerr.NewCustomError("status", err)
	}

	return err
}
