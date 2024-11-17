package httpv1

import (
	"github.com/mondegor/go-sample/internal/catalog/category/section/pub/entity"
)

type (
	// CategoryListResponse - comment struct.
	CategoryListResponse struct {
		Items []entity.Category `json:"items"`
		Total uint64            `json:"total"`
	}
)
