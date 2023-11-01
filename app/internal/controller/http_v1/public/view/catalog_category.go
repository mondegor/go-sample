package view

import (
    "go-sample/internal/entity/public"
)

type (
    CatalogCategoryListResponse struct {
        Items []entity.CatalogCategory `json:"items"`
        Total int64                    `json:"total"`
    }
)
