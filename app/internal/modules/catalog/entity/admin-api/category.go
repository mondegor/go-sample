package entity

import (
	"go-sample/internal/global"
	"time"

	"github.com/mondegor/go-webcore/mrenum"
	"github.com/mondegor/go-webcore/mrtype"
)

const (
	ModelNameCatalogCategory = global.SectionAdminAPI + ".CatalogCategory"
)

type (
	Category struct { // DB: catalog_categories
		ID         mrtype.KeyInt32 `json:"id" sort:"category_id"`
		TagVersion int32           `json:"version"` // tag_version
		CreatedAt  time.Time       `json:"createdAt" sort:"datetime_created"`
		UpdateAt   time.Time       `json:"updateAt"` // datetime_updated

		Caption   string           `json:"caption" sort:"category_caption,default"`
		ImagePath string           `json:"-"` // image_path
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
