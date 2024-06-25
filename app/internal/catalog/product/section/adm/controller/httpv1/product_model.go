package httpv1

import (
	"github.com/mondegor/go-sample/internal/catalog/product/section/adm/entity"
	pkgentity "github.com/mondegor/go-sample/pkg/entity"

	"github.com/google/uuid"
	"github.com/mondegor/go-webcore/mrtype"
)

type (
	// CreateProductRequest - comment struct.
	CreateProductRequest struct {
		CategoryID  uuid.UUID       `json:"categoryId" validate:"required,min=16,max=64"`
		Article     string          `json:"article" validate:"required,min=4,max=32,tag_article"`
		Caption     string          `json:"caption" validate:"required,max=128"`
		TrademarkID mrtype.KeyInt32 `json:"trademarkId" validate:"required,gte=1"`
		Price       pkgentity.Money `json:"price" validate:"gte=0"`
	}

	// StoreProductRequest - comment struct.
	StoreProductRequest struct {
		TagVersion  int32            `json:"tagVersion" validate:"required,gte=1"`
		Article     string           `json:"article" validate:"omitempty,min=4,max=32,tag_article"`
		Caption     string           `json:"caption" validate:"omitempty,max=128"`
		TrademarkID mrtype.KeyInt32  `json:"trademarkId" validate:"omitempty,gte=1"`
		Price       *pkgentity.Money `json:"price" validate:"omitempty,gte=0"`
	}

	// ProductListResponse - comment struct.
	ProductListResponse struct {
		Items []entity.Product `json:"items"`
		Total int64            `json:"total"`
	}
)
