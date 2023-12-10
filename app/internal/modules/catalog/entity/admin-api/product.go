package entity

import (
	entity_shared "go-sample/internal/modules/catalog/entity/shared"
	"time"

	"github.com/mondegor/go-webcore/mrenum"
	"github.com/mondegor/go-webcore/mrtype"
)

const (
	ModelNameProduct = "admin-api.CatalogProduct"
)

type (
	Product struct { // DB: ps_catalog.products
		ID         mrtype.KeyInt32 `json:"id"`                                   // product_id
		TagVersion int32           `json:"version"`                              // tag_version
		CreatedAt  time.Time       `json:"createdAt" sort:"createdAt"`           // datetime_created
		UpdatedAt  *time.Time      `json:"updatedAt,omitempty" sort:"updatedAt"` // datetime_updated

		CategoryID  mrtype.KeyInt32     `json:"categoryId" upd:"category_id"`   // ps_catalog.categories::category_id
		TrademarkID mrtype.KeyInt32     `json:"trademarkId" upd:"trademark_id"` // ps_catalog.trademarks::trademark_id
		Article     string              `json:"article" sort:"article" upd:"product_article"`
		Caption     string              `json:"caption" sort:"caption,default" upd:"product_caption"`
		Price       entity_shared.Money `json:"price" sort:"price" upd:"product_price"` // (coins)

		Status mrenum.ItemStatus `json:"status"` // product_status
	}

	ProductParams struct {
		Filter ProductListFilter
		Sorter mrtype.SortParams
		Pager  mrtype.PageParams
	}

	ProductListFilter struct {
		CategoryID mrtype.KeyInt32
		Trademarks []mrtype.KeyInt32
		SearchText string
		Price      mrtype.RangeInt64
		Statuses   []mrenum.ItemStatus
	}
)
