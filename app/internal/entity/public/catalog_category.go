package entity

import (
    "github.com/mondegor/go-storage/mrentity"
)

const (
    ModelNameCatalogCategory = "public.CatalogCategory"
)

type (
    CatalogCategory struct { // DB: catalog_categories
        Id        mrentity.KeyInt32 `json:"id"` // category_id
        Caption   string `json:"caption"` // category_caption
        ImagePath string `json:"imagePath"` // image_path
    }

    CatalogCategoryParams struct {
        Filter CatalogCategoryListFilter
        Sorter mrentity.ListSorter
        Pager  mrentity.ListPager
    }

    CatalogCategoryListFilter struct {
        SearchText string
    }
)
