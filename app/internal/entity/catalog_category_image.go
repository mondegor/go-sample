package entity

import (
    "github.com/mondegor/go-storage/mrentity"
    "github.com/mondegor/go-storage/mrstorage"
)

const (
    ModelNameCatalogCategoryImage = "CatalogCategoryImage"
)

type (
    CatalogCategoryImageObject struct {
        CategoryId mrentity.KeyInt32
        File mrstorage.File
    }
)
