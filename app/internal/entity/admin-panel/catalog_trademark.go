package entity

import (
    "time"

    "github.com/mondegor/go-storage/mrentity"
    "github.com/mondegor/go-webcore/mrenum"
)

const (
    ModelNameCatalogTrademark = "admin.CatalogTrademark"
)

type (
    CatalogTrademark struct { // DB: catalog_laminate_types
        Id        mrentity.KeyInt32 `json:"id"` // type_id
        Version   mrentity.Version `json:"version"` // tag_version
        CreatedAt time.Time `json:"createdAt"` // datetime_created
        UpdateAt  time.Time `json:"updateAt"` // datetime_updated

        Caption   string `json:"caption" db:"type_caption" sort:"default"`

        Status    mrenum.ItemStatus `json:"status"` // type_status
    }

    CatalogTrademarkParams struct {
        Filter CatalogTrademarkListFilter
        Sorter mrentity.ListSorter
        Pager  mrentity.ListPager
    }

    CatalogTrademarkListFilter struct {
        SearchText string
        Statuses  []mrenum.ItemStatus
    }
)
