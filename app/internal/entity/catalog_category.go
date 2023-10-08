package entity

import (
    "time"

    mrcom_status "github.com/mondegor/go-components/mrcom/status"
    "github.com/mondegor/go-storage/mrentity"
)

const (
    ModelNameCatalogCategory = "CatalogCategory"
)

type (
    CatalogCategory struct { // DB: catalog_categories
        Id        mrentity.KeyInt32 `json:"id"` // category_id
        Version   mrentity.Version `json:"version"` // tag_version
        CreatedAt time.Time `json:"createdAt"` // datetime_created
        UpdateAt  time.Time `json:"updateAt"` // datetime_updated
        Caption   string `json:"caption"` // category_caption
        ImagePath string `json:"imagePath"` // image_path
        Status    mrcom_status.ItemStatus `json:"status"` // category_status
    }

    CatalogCategoryListFilter struct {
        Statuses  []mrcom_status.ItemStatus
    }
)
