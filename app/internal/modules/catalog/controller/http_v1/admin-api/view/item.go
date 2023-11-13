package view

import (
	"github.com/mondegor/go-webcore/mrenum"
	"github.com/mondegor/go-webcore/mrtype"
)

type (
	ChangeItemStatusRequest struct {
		Version int32 `json:"version" validate:"required,gte=1"`
		Status  mrenum.ItemStatus `json:"status" validate:"required"`
	}

	MoveItemRequest struct {
		AfterNodeID mrtype.KeyInt32 `json:"afterId"`
	}

	CreateItemResponse struct {
		ItemID string `json:"id"`
		Message string `json:"message,omitempty"`
	}
)
