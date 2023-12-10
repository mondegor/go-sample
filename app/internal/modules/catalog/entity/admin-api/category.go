package entity

import (
	"time"

	"github.com/mondegor/go-webcore/mrenum"
	"github.com/mondegor/go-webcore/mrtype"
)

const (
	ModelNameCategory = "admin-api.CatalogCategory"
)

type (
	Category struct { // DB: ps_catalog.categories
		ID         mrtype.KeyInt32 `json:"id"`                                   // category_id
		TagVersion int32           `json:"version"`                              // tag_version
		CreatedAt  time.Time       `json:"createdAt" sort:"createdAt"`           // datetime_created
		UpdatedAt  *time.Time      `json:"updatedAt,omitempty" sort:"updatedAt"` // datetime_updated

		Caption   string           `json:"caption" sort:"caption,default"` // category_caption
		ImagePath string           `json:"-"`                              // image_path
		ImageInfo *mrtype.FileInfo `json:"imageInfo,omitempty"`

		Status mrenum.ItemStatus `json:"status"` // category_status
	}

	CategoryParams struct {
		Filter CategoryListFilter
		Sorter mrtype.SortParams
		Pager  mrtype.PageParams
	}

	CategoryListFilter struct {
		SearchText string
		Statuses   []mrenum.ItemStatus
	}
)
