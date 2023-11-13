package view

import (
	"go-sample/internal/modules/catalog/entity/public-api"
)

type (
	CategoryListResponse struct {
		Items []entity.Category `json:"items"`
		Total int64					`json:"total"`
	}
)
