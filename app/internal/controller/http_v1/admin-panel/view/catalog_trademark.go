package view

import (
    "go-sample/internal/entity/admin-panel"

    "github.com/mondegor/go-storage/mrentity"
)

type (
    CreateCatalogTrademarkRequest struct {
        Caption   string `json:"caption" validate:"required,max=64"`
    }

    StoreCatalogTrademarkRequest struct {
        Version   mrentity.Version `json:"version" validate:"required,gte=1"`
        Caption   string `json:"caption" validate:"required,max=64"`
    }

    CatalogTrademarkListResponse struct {
        Items []entity.CatalogTrademark `json:"items"`
        Total int64                     `json:"total"`
    }
)
