package view

import (
	entity "go-sample/internal/modules/catalog/entity/admin-api"
	entity_shared "go-sample/internal/modules/catalog/entity/shared"

	"github.com/mondegor/go-webcore/mrtype"
)

type (
	CreateProductRequest struct {
		CategoryID  mrtype.KeyInt32 `json:"categoryId" validate:"required,gte=1"`
		TrademarkID mrtype.KeyInt32 `json:"trademarkId" validate:"required,gte=1"`
		Article	 string `json:"article" validate:"required,min=4,max=32,article"`
		Caption	 string `json:"caption" validate:"required,max=128"`
		Price	   entity_shared.Money `json:"price" validate:"gte=0,lte=100000000001"`
	}

	StoreProductRequest struct {
		Version	 int32 `json:"version" validate:"required,gte=1"`
		CategoryID  mrtype.KeyInt32 `json:"categoryId" validate:"required,gte=1"`
		TrademarkID mrtype.KeyInt32 `json:"trademarkId" validate:"required,gte=1"`
		Article	 string `json:"article" validate:"required,min=4,max=32,article"`
		Caption	 string `json:"caption" validate:"required,max=128"`
		Price	   entity_shared.Money `json:"price" validate:"gte=0,lte=100000000001"`
	}

	ProductListResponse struct {
		Items []entity.Product `json:"items"`
		Total int64				   `json:"total"`
	}
)
