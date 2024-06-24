package usecase_api

import (
	"context"

	"github.com/google/uuid"
	"github.com/mondegor/go-webcore/mrenum"
)

type (
	// CategoryStorage - comment interface.
	CategoryStorage interface {
		FetchStatus(ctx context.Context, rowID uuid.UUID) (mrenum.ItemStatus, error)
	}
)
