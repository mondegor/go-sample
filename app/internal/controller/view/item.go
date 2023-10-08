package view

import (
    mrcom_status "github.com/mondegor/go-components/mrcom/status"
    "github.com/mondegor/go-storage/mrentity"
)

type (
    ChangeItemStatusRequest struct {
        Version mrentity.Version `json:"version" validate:"required,gte=1"`
        Status  mrcom_status.ItemStatus `json:"status" validate:"required"`
    }
)
