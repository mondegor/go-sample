package view

import (
	"go-sample/internal/modules/catalog/entity/admin-api"
)

type (
	CreateTrademarkRequest struct {
		Caption string `json:"caption" validate:"required,max=128"`
	}

	StoreTrademarkRequest struct {
		Version int32  `json:"version" validate:"required,gte=1"`
		Caption string `json:"caption" validate:"required,max=128"`
	}

	TrademarkListResponse struct {
		Items []entity.Trademark `json:"items"`
		Total int64              `json:"total"`
	}
)
