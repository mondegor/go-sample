package httpv1

import (
	"github.com/mondegor/go-sample/internal/catalog/trademark/section/adm/entity"
)

type (
	// CreateTrademarkRequest - comment struct.
	CreateTrademarkRequest struct {
		Caption string `json:"caption" validate:"required,max=128"`
	}

	// StoreTrademarkRequest - comment struct.
	StoreTrademarkRequest struct {
		TagVersion int32  `json:"tagVersion" validate:"required,gte=1"`
		Caption    string `json:"caption" validate:"required,max=128"`
	}

	// TrademarkListResponse - comment struct.
	TrademarkListResponse struct {
		Items []entity.Trademark `json:"items"`
		Total int64              `json:"total"`
	}
)
