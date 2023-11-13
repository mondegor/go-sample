package http_v1

import (
	"fmt"
	"go-sample/internal/global"
	"go-sample/internal/modules/catalog/controller/http_v1/admin-api/view"
	view_shared "go-sample/internal/modules/catalog/controller/http_v1/shared/view"
	entity "go-sample/internal/modules/catalog/entity/admin-api"
	usecase "go-sample/internal/modules/catalog/usecase/admin-api"
	"net/http"

	"github.com/mondegor/go-webcore/mrcore"
	"github.com/mondegor/go-webcore/mrctx"
	"github.com/mondegor/go-webcore/mrtype"
	"github.com/mondegor/go-webcore/mrview"
)

const (
	trademarkURL = "/v1/catalog/trademarks"
	trademarkItemURL = "/v1/catalog/trademarks/:id"
	trademarkChangeStatusURL = "/v1/catalog/trademarks/:id/status"
)

type (
	Trademark struct {
		section mrcore.ClientSection
		service usecase.TrademarkService
		listSorter mrview.ListSorter
	}
)

func NewTrademark(
	section mrcore.ClientSection,
	service usecase.TrademarkService,
	listSorter mrview.ListSorter,
) *Trademark {
	return &Trademark{
		section: section,
		service: service,
		listSorter: listSorter,
	}
}

func (ht *Trademark) AddHandlers(router mrcore.HttpRouter) {
	moduleAccessFunc := func (next mrcore.HttpHandlerFunc) mrcore.HttpHandlerFunc {
		return ht.section.MiddlewareWithPermission(global.PermissionCatalogTrademark, next)
	}

	router.HttpHandlerFunc(http.MethodGet, ht.section.Path(trademarkURL), moduleAccessFunc(ht.GetList()))
	router.HttpHandlerFunc(http.MethodPost, ht.section.Path(trademarkURL), moduleAccessFunc(ht.Create()))

	router.HttpHandlerFunc(http.MethodGet, ht.section.Path(trademarkItemURL), moduleAccessFunc(ht.Get()))
	router.HttpHandlerFunc(http.MethodPut, ht.section.Path(trademarkItemURL), moduleAccessFunc(ht.Store()))
	router.HttpHandlerFunc(http.MethodDelete, ht.section.Path(trademarkItemURL), moduleAccessFunc(ht.Remove()))

	router.HttpHandlerFunc(http.MethodPut, ht.section.Path(trademarkChangeStatusURL), moduleAccessFunc(ht.ChangeStatus()))
}

func (ht *Trademark) GetList() mrcore.HttpHandlerFunc {
	return func(c mrcore.ClientData) error {
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

func (ht *Trademark) listParams(c mrcore.ClientData) entity.TrademarkParams {
	return entity.TrademarkParams{
		Filter: entity.TrademarkListFilter{
			SearchText: view_shared.ParseFilterString(c, global.ParamNameFilterSearchText),
			Statuses: view_shared.ParseFilterStatusList(c, global.ParamNameFilterStatuses),
		},
		Sorter: view_shared.ParseListSorter(c, ht.listSorter),
		Pager: view_shared.ParseListPager(c),
	}
}

func (ht *Trademark) Get() mrcore.HttpHandlerFunc {
	return func(c mrcore.ClientData) error {
		item, err := ht.service.GetItem(c.Context(), ht.getItemID(c))

		if err != nil {
			return err
		}

		return c.SendResponse(http.StatusOK, item)
	}
}

func (ht *Trademark) Create() mrcore.HttpHandlerFunc {
	return func(c mrcore.ClientData) error {
		request := view.CreateTrademarkRequest{}

		if err := c.ParseAndValidate(&request); err != nil {
			return err
		}

		item := entity.Trademark{
			Caption: request.Caption,
		}

		if err := ht.service.Create(c.Context(), &item); err != nil {
			return err
		}

		return c.SendResponse(
			http.StatusCreated,
			view.CreateItemResponse{
				ItemID: fmt.Sprintf("%d", item.ID),
				Message: mrctx.Locale(c.Context()).TranslateMessage(
					"msgTrademarkSuccessCreated",
					"entity has been success created",
				),
			},
		)
	}
}

func (ht *Trademark) Store() mrcore.HttpHandlerFunc {
	return func(c mrcore.ClientData) error {
		request := view.StoreTrademarkRequest{}

		if err := c.ParseAndValidate(&request); err != nil {
			return err
		}

		item := entity.Trademark{
			ID:		 ht.getItemID(c),
			TagVersion: request.Version,
			Caption:	request.Caption,
		}

		if err := ht.service.Store(c.Context(), &item); err != nil {
			return err
		}

		return c.SendResponseNoContent()
	}
}

func (ht *Trademark) ChangeStatus() mrcore.HttpHandlerFunc {
	return func(c mrcore.ClientData) error {
		request := view.ChangeItemStatusRequest{}

		if err := c.ParseAndValidate(&request); err != nil {
			return err
		}

		item := entity.Trademark{
			ID:		 ht.getItemID(c),
			TagVersion: request.Version,
			Status:	 request.Status,
		}

		if err := ht.service.ChangeStatus(c.Context(), &item); err != nil {
			return err
		}

		return c.SendResponseNoContent()
	}
}

func (ht *Trademark) Remove() mrcore.HttpHandlerFunc {
	return func(c mrcore.ClientData) error {
		if err := ht.service.Remove(c.Context(), ht.getItemID(c)); err != nil {
			return err
		}

		return c.SendResponseNoContent()
	}
}

func (ht *Trademark) getItemID(c mrcore.ClientData) mrtype.KeyInt32 {
	return mrtype.KeyInt32(c.RequestPath().GetInt64("id"))
}
