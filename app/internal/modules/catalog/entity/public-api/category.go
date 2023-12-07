package entity

import (
	"github.com/mondegor/go-webcore/mrtype"
)

const (
	ModelNameCatalogCategory = "public-api.CatalogCategory"
)

type (
	Category struct { // DB: catalog_categories
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
