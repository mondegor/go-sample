package entity

import (
	"github.com/google/uuid"
)

const (
	ModelNameCategory = "public-api.Catalog.Category"
)

type (
	Category struct { // DB: ps_catalog.categories
		ID       uuid.UUID `json:"id"`                 // category_id
		Caption  string    `json:"caption"`            // category_caption
		ImageURL string    `json:"imageUrl,omitempty"` // image_meta.path
	}

	CategoryParams struct {
		LanguageID uint16
		Filter     CategoryListFilter
	}

	CategoryListFilter struct {
		SearchText string
	}
)
