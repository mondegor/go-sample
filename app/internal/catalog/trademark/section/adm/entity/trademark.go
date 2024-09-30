package entity

import (
	"time"

	"github.com/mondegor/go-webcore/mrenum"
	"github.com/mondegor/go-webcore/mrtype"
)

const (
	ModelNameTrademark = "admin-api.Catalog.Trademark" // ModelNameTrademark - название сущности
)

type (
	// Trademark - comment struct.
	Trademark struct { // DB: ps_catalog.trademarks
		ID         mrtype.KeyInt32   `json:"id"` // trademark_id
		TagVersion int32             `json:"tagVersion"`
		Caption    string            `json:"caption" sort:"caption,default"`
		Status     mrenum.ItemStatus `json:"status"`
		CreatedAt  time.Time         `json:"createdAt" sort:"createdAt"`
		UpdatedAt  time.Time         `json:"updatedAt" sort:"updatedAt"`
	}

	// TrademarkParams - comment struct.
	TrademarkParams struct {
		Filter TrademarkListFilter
		Sorter mrtype.SortParams
		Pager  mrtype.PageParams
	}

	// TrademarkListFilter - comment struct.
	TrademarkListFilter struct {
		SearchText string
		Statuses   []mrenum.ItemStatus
	}
)
