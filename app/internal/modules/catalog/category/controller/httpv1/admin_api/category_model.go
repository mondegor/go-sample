package httpv1

import entity "go-sample/internal/modules/catalog/category/entity/admin_api"

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
