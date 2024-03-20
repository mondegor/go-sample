package entity

import (
	"time"

	"github.com/google/uuid"
	"github.com/mondegor/go-storage/mrentity"
	"github.com/mondegor/go-webcore/mrenum"
	"github.com/mondegor/go-webcore/mrtype"
)

const (
	ModelNameCategory      = "admin-api.Catalog.Category"
	ModelNameCategoryImage = "admin-api.Catalog.CategoryImage"
)

type (
	Category struct { // DB: ps_catalog.categories
		ID         uuid.UUID         `json:"id"` // category_id
		TagVersion int32             `json:"tagVersion"`
		Caption    string            `json:"caption" sort:"caption,default"`
		Status     mrenum.ItemStatus `json:"status"`
		CreatedAt  time.Time         `json:"createdAt" sort:"createdAt"`
		UpdatedAt  *time.Time        `json:"updatedAt,omitempty" sort:"updatedAt"`

		ImageInfo *mrtype.ImageInfo   `json:"imageInfo,omitempty"`
		ImageMeta *mrentity.ImageMeta `json:"-"` // image_meta
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
