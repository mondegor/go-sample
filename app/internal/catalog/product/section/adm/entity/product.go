package entity

import (
	"time"

	"github.com/google/uuid"
	"github.com/mondegor/go-webcore/mrenum"
	"github.com/mondegor/go-webcore/mrtype"

	"github.com/mondegor/go-sample/pkg/measure"
)

const (
	ModelNameProduct = "admin-api.Catalog.Product" // ModelNameProduct - название сущности
)

type (
	// Product - comment struct.
	Product struct { // DB: ps_catalog.products
		ID          mrtype.KeyInt32   `json:"id"` // product_id
		TagVersion  int32             `json:"tagVersion"`
		CategoryID  uuid.UUID         `json:"categoryId" upd:"category_id"` // categories::category_id
		Article     string            `json:"article" sort:"article" upd:"product_article"`
		Caption     string            `json:"caption" sort:"caption,default" upd:"product_caption"`
		TrademarkID mrtype.KeyInt32   `json:"trademarkId" upd:"trademark_id"`         // trademarks::trademark_id
		Price       *measure.Money    `json:"price" sort:"price" upd:"product_price"` // coins
		Status      mrenum.ItemStatus `json:"status"`
		CreatedAt   time.Time         `json:"createdAt" sort:"createdAt"`
		UpdatedAt   time.Time         `json:"updatedAt" sort:"updatedAt"`
	}

	// ProductParams - comment struct.
	ProductParams struct {
		Filter ProductListFilter
		Sorter mrtype.SortParams
		Pager  mrtype.PageParams
	}

	// ProductListFilter - comment struct.
	ProductListFilter struct {
		CategoryID   uuid.UUID
		SearchText   string
		TrademarkIDs []mrtype.KeyInt32
		Price        mrtype.RangeInt64
		Statuses     []mrenum.ItemStatus
	}
)
