package entity

import (
	"time"

	"github.com/mondegor/go-webcore/mrenum"
	"github.com/mondegor/go-webcore/mrtype"
)

const (
	ModelNameCatalogTrademark = "admin-api.CatalogTrademark"
)

type (
	Trademark struct { // DB: catalog_trademarks
		ID         mrtype.KeyInt32 `json:"id" sort:"type_id"`
		TagVersion int32           `json:"version"` // tag_version
		CreatedAt  time.Time       `json:"createdAt" sort:"datetime_created"`
		UpdateAt   time.Time       `json:"updateAt"` // datetime_updated

		Caption string `json:"caption" sort:"trademark_caption,default"`

		Status mrenum.ItemStatus `json:"status"` // type_status
	}

	TrademarkParams struct {
		Filter TrademarkListFilter
		Sorter mrtype.SortParams
		Pager  mrtype.PageParams
	}

	TrademarkListFilter struct {
		SearchText string
		Statuses   []mrenum.ItemStatus
	}
)
