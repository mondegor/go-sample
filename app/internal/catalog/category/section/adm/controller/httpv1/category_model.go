package httpv1

import (
	"github.com/mondegor/go-sample/internal/catalog/category/section/adm/entity"
)

type (
	// CreateCategoryRequest - comment struct.
	CreateCategoryRequest struct {
		Caption string `json:"caption" validate:"required,max=128"`
	}

	// StoreCategoryRequest - comment struct.
	StoreCategoryRequest struct {
		TagVersion int32  `json:"tagVersion" validate:"required,gte=1"`
		Caption    string `json:"caption" validate:"required,max=128"`
	}

	// CategoryListResponse - comment struct.
	CategoryListResponse struct {
		Items []entity.Category `json:"items"`
		Total int64             `json:"total"`
	}
)
