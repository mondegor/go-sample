package entity

import (
    entity_shared "go-sample/internal/entity/shared"
    "time"

    "github.com/mondegor/go-storage/mrentity"
    "github.com/mondegor/go-webcore/mrenum"
)

const (
    ModelNameCatalogProduct = "admin.CatalogProduct"
)

type (
    CatalogProduct struct { // DB: catalog_products
        Id         mrentity.KeyInt32 `json:"id"` // product_id
        Version    mrentity.Version `json:"version"` // tag_version
        CreatedAt  time.Time `json:"createdAt"` // datetime_created
        UpdateAt   time.Time `json:"updateAt"` // datetime_updated

        CategoryId mrentity.KeyInt32 `json:"categoryId" db:"category_id" update:"on"` // catalog_categories::category_id
        TrademarkId    mrentity.KeyInt32 `json:"trademarkId" db:"trademark_id" update:"on"` // catalog_trademarks::trademark_id
        Article   string `json:"article" db:"product_article" sort:"on" update:"on"`
        Caption   string `json:"caption" db:"product_caption" sort:"default" update:"on"`
        Price     entity_shared.Money `json:"price" db:"product_price" sort:"on" update:"on"` // (coins)

        Status    mrenum.ItemStatus `json:"status"` // product_status
    }

    CatalogProductParams struct {
        Filter CatalogProductListFilter
        Sorter mrentity.ListSorter
        Pager  mrentity.ListPager
    }

    CatalogProductListFilter struct {
        CategoryId mrentity.KeyInt32
        Trademarks []mrentity.KeyInt32
        SearchText string
        Price mrentity.RangeInt64
        Statuses  []mrenum.ItemStatus
    }
)
