package entity

import (
    "time"

    "github.com/mondegor/go-components/mrcom"
    "github.com/mondegor/go-storage/mrentity"
)

const (
    ModelNameCatalogProduct = "CatalogProduct"
)

type (
    CatalogProduct struct { // DB: catalog_products
        Id         mrentity.KeyInt32 `json:"id"` // product_id
        Version    mrentity.Version `json:"version"` // tag_version
        CreatedAt  time.Time `json:"createdAt"` // datetime_created
        UpdateAt   time.Time `json:"updateAt"` // datetime_updated
        CategoryId mrentity.KeyInt32 `json:"categoryId"` // catalog_categories::category_id

        TrademarkId    mrentity.KeyInt32 `json:"trademarkId" db:"trademark_id"` // catalog_trademarks::trademark_id
        Article   string `json:"article" db:"product_article"`
        Caption   string `json:"caption" db:"product_caption"`
        Price     Money `json:"price" db:"product_price"` // (coins)

        Status    mrcom.ItemStatus `json:"status"` // product_status
    }

    CatalogProductListFilter struct {
        CategoryId mrentity.KeyInt32
        Trademarks []mrentity.KeyInt32
        Statuses   []mrcom.ItemStatus
    }
)
