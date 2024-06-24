package httpv1

import entity "go-sample/internal/modules/catalog/category/entity/public_api"

type (
	// CategoryListResponse - comment struct.
	CategoryListResponse struct {
		Items []entity.Category `json:"items"`
		Total int64             `json:"total"`
	}
)
