package entity

import (
	"github.com/google/uuid"
)

const (
	ModelNameCategory = "public-api.Catalog.Category" // ModelNameCategory - название сущности
)

type (
	// Category - comment struct.
	Category struct { // DB: ps_catalog.categories
		ID       uuid.UUID `json:"id"` // category_id
		Caption  string    `json:"caption"`
		ImageURL string    `json:"imageUrl,omitempty"` // image_meta.path
	}

	// CategoryParams - comment struct.
	CategoryParams struct {
		LanguageID uint16
		Filter     CategoryListFilter
	}

	// CategoryListFilter - comment struct.
	CategoryListFilter struct {
		SearchText string
	}
)
