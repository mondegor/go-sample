package entity

import (
	"github.com/mondegor/go-webcore/mrtype"
)

const (
	ModelNameCategory = "public-api.CatalogCategory"
)

type (
	Category struct { // DB: ps_catalog.categories
		ID        mrtype.KeyInt32 `json:"id"`       // category_id
		Caption   string          `json:"caption"`  // category_caption
		ImagePath string          `json:"imageURL"` // image_path
	}

	CategoryParams struct {
		Filter CategoryListFilter
	}

	CategoryListFilter struct {
		SearchText string
	}
)
