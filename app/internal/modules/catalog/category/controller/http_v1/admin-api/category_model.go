package http_v1

import (
	"go-sample/internal/modules/catalog/category/entity/admin-api"
)

type (
	CreateCategoryRequest struct {
		Caption string `json:"caption" validate:"required,max=128"`
	}

	StoreCategoryRequest struct {
		Version int32  `json:"version" validate:"required,gte=1"`
		Caption string `json:"caption" validate:"required,max=128"`
	}

	CategoryListResponse struct {
		Items []entity.Category `json:"items"`
		Total int64             `json:"total"`
	}
)
