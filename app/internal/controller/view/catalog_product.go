package view

import (
    "go-sample/internal/entity"

    "github.com/mondegor/go-storage/mrentity"
)

type (
    CreateCatalogProductRequest struct {
        TrademarkId mrentity.KeyInt32 `json:"trademarkId" validate:"required,gte=1"`
        Article     string `json:"article" validate:"required,min=3,max=32,article"`
        Caption     string `json:"caption" validate:"required,max=128"`
        Price       entity.Money `json:"price" validate:"required,gte=1,lte=100000000001"`
    }

    StoreCatalogProductRequest struct {
        Version     mrentity.Version `json:"version" validate:"required,gte=1"`
        TrademarkId mrentity.KeyInt32 `json:"trademarkId" validate:"required,gte=1"`
        Article     string `json:"article" validate:"required,min=3,max=32,article"`
        Caption     string `json:"caption" validate:"required,max=128"`
        Price       entity.Money `json:"price" validate:"required,gte=1,lte=100000000001"`
    }

    MoveCatalogProductRequest struct {
        AfterNodeId mrentity.KeyInt32 `json:"afterId"`
    }
)
