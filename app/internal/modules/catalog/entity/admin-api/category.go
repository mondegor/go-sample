package entity

import (
	"time"

	"github.com/mondegor/go-storage/mrentity"
	"github.com/mondegor/go-webcore/mrenum"
	"github.com/mondegor/go-webcore/mrtype"
)

const (
	ModelNameCategory = "admin-api.Catalog.Category"
)

type (
	Category struct { // DB: ps_catalog.categories
		ID         mrtype.KeyInt32 `json:"id"`                                   // category_id
		TagVersion int32           `json:"version"`                              // tag_version
		CreatedAt  time.Time       `json:"createdAt" sort:"createdAt"`           // datetime_created
		UpdatedAt  *time.Time      `json:"updatedAt,omitempty" sort:"updatedAt"` // datetime_updated

		Caption   string              `json:"caption" sort:"caption,default"` // category_caption
		ImageInfo *mrtype.ImageInfo   `json:"imageInfo,omitempty"`
		ImageMeta *mrentity.ImageMeta `json:"-"` // image_meta

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
