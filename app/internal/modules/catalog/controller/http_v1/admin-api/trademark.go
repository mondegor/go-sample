package http_v1

import (
	module "go-sample/internal/modules/catalog"
	"go-sample/internal/modules/catalog/controller/http_v1/admin-api/view"
	view_shared "go-sample/internal/modules/catalog/controller/http_v1/shared/view"
	entity "go-sample/internal/modules/catalog/entity/admin-api"
	usecase "go-sample/internal/modules/catalog/usecase/admin-api"
	usecase_shared "go-sample/internal/modules/catalog/usecase/shared"
	"net/http"
	"strconv"

	"github.com/mondegor/go-sysmess/mrerr"
	"github.com/mondegor/go-webcore/mrcore"
	"github.com/mondegor/go-webcore/mrctx"
	"github.com/mondegor/go-webcore/mrtype"
	"github.com/mondegor/go-webcore/mrview"
)

const (
	trademarkURL             = "/v1/catalog/trademarks"
	trademarkItemURL         = "/v1/catalog/trademarks/:id"
	trademarkChangeStatusURL = "/v1/catalog/trademarks/:id/status"
)

type (
	Trademark struct {
		section    mrcore.ClientSection
		service    usecase.TrademarkService
		listSorter mrview.ListSorter
	}
)

func NewTrademark(
	section mrcore.ClientSection,
	service usecase.TrademarkService,
	listSorter mrview.ListSorter,
) *Trademark {
	return &Trademark{
		section:    section,
		service:    service,
		listSorter: listSorter,
	}
}

func (ht *Trademark) AddHandlers(router mrcore.HttpRouter) {
	moduleAccessFunc := func(next mrcore.HttpHandlerFunc) mrcore.HttpHandlerFunc {
		return ht.section.MiddlewareWithPermission(module.PermissionCatalogTrademark, next)
	}

	router.HttpHandlerFunc(http.MethodGet, ht.section.Path(trademarkURL), moduleAccessFunc(ht.GetList()))
	router.HttpHandlerFunc(http.MethodPost, ht.section.Path(trademarkURL), moduleAccessFunc(ht.Create()))

	router.HttpHandlerFunc(http.MethodGet, ht.section.Path(trademarkItemURL), moduleAccessFunc(ht.Get()))
	router.HttpHandlerFunc(http.MethodPut, ht.section.Path(trademarkItemURL), moduleAccessFunc(ht.Store()))
	router.HttpHandlerFunc(http.MethodDelete, ht.section.Path(trademarkItemURL), moduleAccessFunc(ht.Remove()))

	router.HttpHandlerFunc(http.MethodPut, ht.section.Path(trademarkChangeStatusURL), moduleAccessFunc(ht.ChangeStatus()))
}

func (ht *Trademark) GetList() mrcore.HttpHandlerFunc {
	return func(c mrcore.ClientContext) error {
		items, totalItems, err := ht.service.GetList(c.Context(), ht.listParams(c))

		if err != nil {
			return err
		}

		return c.SendResponse(
			http.StatusOK,
			view.TrademarkListResponse{
				Items: items,
				Total: totalItems,
			},
		)
	}
}

func (ht *Trademark) listParams(c mrcore.ClientContext) entity.TrademarkParams {
	return entity.TrademarkParams{
		Filter: entity.TrademarkListFilter{
			SearchText: view_shared.ParseFilterString(c, module.ParamNameFilterSearchText),
			Statuses:   view_shared.ParseFilterStatusList(c, module.ParamNameFilterStatuses),
		},
		Sorter: view_shared.ParseSortParams(c, ht.listSorter),
		Pager:  view_shared.ParsePageParams(c),
	}
}

func (ht *Trademark) Get() mrcore.HttpHandlerFunc {
	return func(c mrcore.ClientContext) error {
		item, err := ht.service.GetItem(c.Context(), ht.getItemID(c))

		if err != nil {
			return ht.wrapError(err, ht.getRawItemID(c))
		}

		return c.SendResponse(http.StatusOK, item)
	}
}

func (ht *Trademark) Create() mrcore.HttpHandlerFunc {
	return func(c mrcore.ClientContext) error {
		request := view.CreateTrademarkRequest{}

		if err := c.Validate(&request); err != nil {
			return err
		}

		item := entity.Trademark{
			Caption: request.Caption,
		}

		if err := ht.service.Create(c.Context(), &item); err != nil {
			return ht.wrapError(err, ht.getRawItemID(c))
		}

		return c.SendResponse(
			http.StatusCreated,
			view.SuccessCreatedItemResponse{
				ItemID: strconv.Itoa(int(item.ID)),
				Message: mrctx.Locale(c.Context()).TranslateMessage(
					"msgTrademarkSuccessCreated",
					"entity has been success created",
				),
			},
		)
	}
}

func (ht *Trademark) Store() mrcore.HttpHandlerFunc {
	return func(c mrcore.ClientContext) error {
		request := view.StoreTrademarkRequest{}

		if err := c.Validate(&request); err != nil {
			return err
		}

		item := entity.Trademark{
			ID:         ht.getItemID(c),
			TagVersion: request.Version,
			Caption:    request.Caption,
		}

		if err := ht.service.Store(c.Context(), &item); err != nil {
			return ht.wrapError(err, ht.getRawItemID(c))
		}

		return c.SendResponseNoContent()
	}
}

func (ht *Trademark) ChangeStatus() mrcore.HttpHandlerFunc {
	return func(c mrcore.ClientContext) error {
		request := view.ChangeItemStatusRequest{}

		if err := c.Validate(&request); err != nil {
			return err
		}

		item := entity.Trademark{
			ID:         ht.getItemID(c),
			TagVersion: request.TagVersion,
			Status:     request.Status,
		}

		if err := ht.service.ChangeStatus(c.Context(), &item); err != nil {
			return ht.wrapError(err, ht.getRawItemID(c))
		}

		return c.SendResponseNoContent()
	}
}

func (ht *Trademark) Remove() mrcore.HttpHandlerFunc {
	return func(c mrcore.ClientContext) error {
		if err := ht.service.Remove(c.Context(), ht.getItemID(c)); err != nil {
			return ht.wrapError(err, ht.getRawItemID(c))
		}

		return c.SendResponseNoContent()
	}
}

func (ht *Trademark) getItemID(c mrcore.ClientContext) mrtype.KeyInt32 {
	return view_shared.ParseIDFromPath(c, "id")
}

func (ht *Trademark) getRawItemID(c mrcore.ClientContext) string {
	return c.ParamFromPath("id")
}

func (ht *Trademark) wrapError(err error, rawItemID string) error {
	if mrcore.FactoryErrServiceEntityNotFound.Is(err) {
		return usecase_shared.FactoryErrTrademarkNotFound.Wrap(err, rawItemID)
	}

	if mrcore.FactoryErrServiceEntityVersionInvalid.Is(err) {
		return mrerr.NewFieldError("version", err)
	}

	if mrcore.FactoryErrServiceSwitchStatusRejected.Is(err) {
		return mrerr.NewFieldError("status", err)
	}

	return err
}
