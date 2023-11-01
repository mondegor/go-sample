package view

import (
    entity "go-sample/internal/entity/admin-panel"
    entity_shared "go-sample/internal/entity/shared"

    "github.com/mondegor/go-storage/mrentity"
)

type (
    CreateCatalogProductRequest struct {
        CategoryId  mrentity.KeyInt32 `json:"categoryId" validate:"required,gte=1"`
        TrademarkId mrentity.KeyInt32 `json:"trademarkId" validate:"required,gte=1"`
        Article     string `json:"article" validate:"required,min=3,max=32,article"`
        Caption     string `json:"caption" validate:"required,max=128"`
        Price       entity_shared.Money `json:"price" validate:"required,gte=1,lte=100000000001"`
    }

    StoreCatalogProductRequest struct {
        Version     mrentity.Version `json:"version" validate:"required,gte=1"`
        CategoryId  mrentity.KeyInt32 `json:"categoryId" validate:"required,gte=1"`
        TrademarkId mrentity.KeyInt32 `json:"trademarkId" validate:"required,gte=1"`
        Article     string `json:"article" validate:"required,min=3,max=32,article"`
        Caption     string `json:"caption" validate:"required,max=128"`
        Price       entity_shared.Money `json:"price" validate:"required,gte=1,lte=100000000001"`
    }

    MoveCatalogProductRequest struct {
        AfterNodeId mrentity.KeyInt32 `json:"afterId"`
    }

    CatalogProductListResponse struct {
        Items []entity.CatalogProduct `json:"items"`
        Total int64                   `json:"total"`
    }
)
