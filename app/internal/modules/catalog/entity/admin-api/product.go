package entity

import (
	"go-sample/internal/global"
	entity_shared "go-sample/internal/modules/catalog/entity/shared"
	"time"

	"github.com/mondegor/go-webcore/mrenum"
	"github.com/mondegor/go-webcore/mrtype"
)

const (
	ModelNameCatalogProduct = global.SectionAdminAPI + ".CatalogProduct"
)

type (
	Product struct { // DB: catalog_products
		ID         mrtype.KeyInt32 `json:"id" sort:"product_id"` // product_id
		TagVersion int32           `json:"version"`              // tag_version
		CreatedAt  time.Time       `json:"createdAt" sort:"datetime_created"`
		UpdateAt   time.Time       `json:"updateAt"` // datetime_updated

		CategoryID  mrtype.KeyInt32     `json:"categoryId" upd:"category_id"`   // catalog_categories::category_id
		TrademarkID mrtype.KeyInt32     `json:"trademarkId" upd:"trademark_id"` // catalog_trademarks::trademark_id
		Article     string              `json:"article" sort:"product_article" upd:"product_article"`
		Caption     string              `json:"caption" sort:"product_caption,default" upd:"product_caption"`
		Price       entity_shared.Money `json:"price" sort:"product_price" upd:"product_price"` // (coins)

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
