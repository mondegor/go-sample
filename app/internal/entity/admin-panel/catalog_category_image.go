package entity

import (
    "github.com/mondegor/go-storage/mrentity"
)

const (
    ModelNameCatalogCategoryImage = "admin.CatalogCategoryImage"
)

type (
    CatalogCategoryImageObject struct {
        CategoryId mrentity.KeyInt32
        File mrentity.File
    }
)
