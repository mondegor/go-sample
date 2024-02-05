package entity

import (
	"time"

	"github.com/mondegor/go-webcore/mrenum"
	"github.com/mondegor/go-webcore/mrtype"
)

const (
	ModelNameTrademark = "admin-api.Catalog.Trademark"
)

type (
	Trademark struct { // DB: ps_catalog.trademarks
		ID         mrtype.KeyInt32 `json:"id"`                                   // type_id
		TagVersion int32           `json:"version"`                              // tag_version
		CreatedAt  time.Time       `json:"createdAt" sort:"createdAt"`           // created_at
		UpdatedAt  *time.Time      `json:"updatedAt,omitempty" sort:"updatedAt"` // updated_at

		Caption string `json:"caption" sort:"caption,default"` // category_caption

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
