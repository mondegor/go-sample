package view

import (
    "github.com/mondegor/go-storage/mrentity"
    "github.com/mondegor/go-webcore/mrenum"
)

type (
    ChangeItemStatusRequest struct {
        Version mrentity.Version `json:"version" validate:"required,gte=1"`
        Status  mrenum.ItemStatus `json:"status" validate:"required"`
    }
)
