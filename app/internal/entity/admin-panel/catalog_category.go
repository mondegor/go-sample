package entity

import (
    "time"

    "github.com/mondegor/go-storage/mrentity"
    "github.com/mondegor/go-webcore/mrenum"
)

const (
    ModelNameCatalogCategory = "admin.CatalogCategory"
)

type (
    CatalogCategory struct { // DB: catalog_categories
        Id        mrentity.KeyInt32 `json:"id"` // category_id
        Version   mrentity.Version `json:"version"` // tag_version
        CreatedAt time.Time `json:"createdAt"` // datetime_created
        UpdateAt  time.Time `json:"updateAt"` // datetime_updated

        Caption   string `json:"caption" db:"category_caption" sort:"default"`
        ImagePath string `json:"imagePath"` // image_path

        Status    mrenum.ItemStatus `json:"status"` // category_status
    }

    CatalogCategoryParams struct {
        Filter CatalogCategoryListFilter
        Sorter mrentity.ListSorter
        Pager  mrentity.ListPager
    }

    CatalogCategoryListFilter struct {
        SearchText string
        Statuses  []mrenum.ItemStatus
    }
)
