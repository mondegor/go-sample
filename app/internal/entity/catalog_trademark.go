package entity

import (
    "time"

    mrcom_status "github.com/mondegor/go-components/mrcom/status"
    "github.com/mondegor/go-storage/mrentity"
)

const (
    ModelNameCatalogTrademark = "CatalogTrademark"
)

type (
    CatalogTrademark struct { // DB: catalog_laminate_types
        Id        mrentity.KeyInt32 `json:"id"` // type_id
        Version   mrentity.Version `json:"version"` // tag_version
        CreatedAt time.Time `json:"createdAt"` // datetime_created
        UpdateAt  time.Time `json:"updateAt"` // datetime_updated
        Caption   string `json:"caption"` // type_caption
        Status    mrcom_status.ItemStatus `json:"status"` // type_status
    }

    CatalogTrademarkListFilter struct {
        Statuses  []mrcom_status.ItemStatus
    }
)