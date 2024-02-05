package http_v1

import (
	"go-sample/internal/modules/catalog/category/entity/public-api"
)

type (
	CategoryListResponse struct {
		Items []entity.Category `json:"items"`
		Total int64             `json:"total"`
	}
)
