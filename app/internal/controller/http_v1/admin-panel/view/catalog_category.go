package view

import (
    "go-sample/internal/entity/admin-panel"

    "github.com/mondegor/go-storage/mrentity"
)

type (
    CreateCatalogCategoryRequest struct {
        Caption  string `json:"caption" validate:"required,max=128"`
    }

    StoreCatalogCategoryRequest struct {
        Version  mrentity.Version `json:"version" validate:"required,gte=1"`
        Caption  string `json:"caption" validate:"required,max=128"`
    }

    CatalogCategoryListResponse struct {
        Items []entity.CatalogCategory `json:"items"`
        Total int64                    `json:"total"`
    }
)
