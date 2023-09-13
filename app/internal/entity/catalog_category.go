package entity

import (
    "time"

    "github.com/mondegor/go-components/mrcom"
    "github.com/mondegor/go-storage/mrentity"
)

const ModelNameCatalogCategory = "CatalogCategory"

type (
    CatalogCategory struct { // DB: catalog_categories
        Id        mrentity.KeyInt32 `json:"id"` // category_id
        Version   mrentity.Version `json:"version"` // tag_version
        CreatedAt time.Time `json:"createdAt"` // datetime_created
        Caption   string `json:"caption"` // category_caption
        Status    mrcom.ItemStatus `json:"status"` // category_status
    }

    CatalogCategoryListFilter struct {
        Statuses  []mrcom.ItemStatus
    }
)
